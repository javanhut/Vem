package appcore

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/javanhut/vem/internal/lsp"
)

type formatterConfig struct {
	Name             string
	Command          string
	Args             []string
	UsesStdin        bool
	StdinFilePathArg bool
}

type linterConfig struct {
	Name    string
	Command string
	Args    []string
}

type toolInstallSpec struct {
	name            string
	goInstall       string
	rustupComponent string
	cargoCrate      string
	npmPackages     []string
	pipPackage      string
	luarocksPackage string
	brewPackages    []string
	aptPackages     []string
	dnfPackages     []string
	pacmanPackages  []string
	zypperPackages  []string
}

func toolInstallSpecsForLanguage(language string) []toolInstallSpec {
	switch language {
	case "go":
		return []toolInstallSpec{
			{
				name:           "golangci-lint",
				goInstall:      "github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
				brewPackages:   []string{"golangci-lint"},
				aptPackages:    []string{"golangci-lint"},
				dnfPackages:    []string{"golangci-lint"},
				pacmanPackages: []string{"golangci-lint"},
			},
		}
	case "python":
		return []toolInstallSpec{
			{name: "black", pipPackage: "black", brewPackages: []string{"black"}},
			{name: "ruff", pipPackage: "ruff", brewPackages: []string{"ruff"}},
		}
	case "rust":
		return []toolInstallSpec{
			{name: "rustfmt", rustupComponent: "rustfmt"},
			{name: "clippy", rustupComponent: "clippy"},
		}
	case "c", "cpp":
		return []toolInstallSpec{
			{
				name:           "clang-format",
				brewPackages:   []string{"llvm", "clang-format"},
				aptPackages:    []string{"clang-format", "clang-tools"},
				dnfPackages:    []string{"clang-tools-extra", "clang-format"},
				pacmanPackages: []string{"clang", "clang-tools-extra"},
				zypperPackages: []string{"clang-tools"},
			},
			{
				name:           "clang-tidy",
				brewPackages:   []string{"llvm", "clang-tidy"},
				aptPackages:    []string{"clang-tidy", "clang-tools"},
				dnfPackages:    []string{"clang-tools-extra", "clang-tidy"},
				pacmanPackages: []string{"clang", "clang-tools-extra"},
				zypperPackages: []string{"clang-tools"},
			},
		}
	case "lua":
		return []toolInstallSpec{
			{
				name:           "stylua",
				cargoCrate:     "stylua",
				brewPackages:   []string{"stylua"},
				aptPackages:    []string{"stylua"},
				dnfPackages:    []string{"stylua"},
				pacmanPackages: []string{"stylua"},
			},
			{
				name:            "luacheck",
				luarocksPackage: "luacheck",
				brewPackages:    []string{"luacheck"},
				aptPackages:     []string{"luacheck"},
				dnfPackages:     []string{"luacheck"},
				pacmanPackages:  []string{"luacheck"},
				zypperPackages:  []string{"luacheck"},
			},
		}
	case "javascript", "typescript", "js", "ts":
		return []toolInstallSpec{
			{name: "prettier", npmPackages: []string{"prettier"}},
			{name: "eslint", npmPackages: []string{"eslint"}},
		}
	case "html":
		return []toolInstallSpec{
			{name: "prettier", npmPackages: []string{"prettier"}},
			{name: "htmlhint", npmPackages: []string{"htmlhint"}},
		}
	default:
		return nil
	}
}

func getToolingForFile(filePath string) (formatterConfig, *linterConfig) {
	lang := languageKeyForFile(filePath)
	switch lang {
	case "go":
		return formatterConfig{
				Name:      "gofmt",
				Command:   "gofmt",
				Args:      []string{},
				UsesStdin: true,
			}, &linterConfig{
				Name:    "golangci-lint",
				Command: "golangci-lint",
				Args:    []string{"run", "{file}"},
			}
	case "python":
		return formatterConfig{
				Name:             "black",
				Command:          "black",
				Args:             []string{"--quiet", "--stdin-filename", "{file}", "-"},
				UsesStdin:        true,
				StdinFilePathArg: true,
			}, &linterConfig{
				Name:    "ruff",
				Command: "ruff",
				Args:    []string{"check", "{file}"},
			}
	case "rust":
		return formatterConfig{
				Name:      "rustfmt",
				Command:   "rustfmt",
				Args:      []string{"--emit=stdout"},
				UsesStdin: true,
			}, &linterConfig{
				Name:    "cargo-clippy",
				Command: "cargo",
				Args:    []string{"clippy", "--quiet"},
			}
	case "c", "cpp":
		return formatterConfig{
				Name:             "clang-format",
				Command:          "clang-format",
				Args:             []string{"-assume-filename", "{file}"},
				UsesStdin:        true,
				StdinFilePathArg: true,
			}, &linterConfig{
				Name:    "clang-tidy",
				Command: "clang-tidy",
				Args:    []string{"{file}"},
			}
	case "lua":
		return formatterConfig{
				Name:             "stylua",
				Command:          "stylua",
				Args:             []string{"--stdin-filepath", "{file}", "-"},
				UsesStdin:        true,
				StdinFilePathArg: true,
			}, &linterConfig{
				Name:    "luacheck",
				Command: "luacheck",
				Args:    []string{"{file}"},
			}
	case "js", "ts":
		return formatterConfig{
				Name:             "prettier",
				Command:          "prettier",
				Args:             []string{"--stdin-filepath", "{file}"},
				UsesStdin:        true,
				StdinFilePathArg: true,
			}, &linterConfig{
				Name:    "eslint",
				Command: "eslint",
				Args:    []string{"{file}"},
			}
	case "html":
		return formatterConfig{
				Name:             "prettier",
				Command:          "prettier",
				Args:             []string{"--stdin-filepath", "{file}"},
				UsesStdin:        true,
				StdinFilePathArg: true,
			}, &linterConfig{
				Name:    "htmlhint",
				Command: "htmlhint",
				Args:    []string{"{file}"},
			}
	default:
		return formatterConfig{}, nil
	}
}

func languageKeyForFile(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".py", ".pyi":
		return "python"
	case ".rs":
		return "rust"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc", ".hh", ".cxx", ".hxx", ".c++", ".h++":
		return "cpp"
	case ".lua":
		return "lua"
	case ".js", ".jsx", ".mjs", ".cjs", ".ts", ".tsx":
		return "js"
	case ".html", ".htm":
		return "html"
	default:
		return ""
	}
}

func formatBufferWithTool(bufContent, filePath string, cfg formatterConfig) (string, error) {
	if cfg.Command == "" {
		return "", errors.New("no formatter configured")
	}
	cmdPath := resolveToolPath(cfg.Command)
	if cmdPath == "" {
		return "", errors.New("formatter not installed")
	}
	args := replaceFileArgs(cfg.Args, filePath)
	cmd := exec.Command(cmdPath, args...)

	if cfg.UsesStdin {
		cmd.Stdin = strings.NewReader(bufContent)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}

	return out.String(), nil
}

func lintFileWithTool(filePath string, cfg linterConfig) error {
	if cfg.Command == "" {
		return nil
	}
	cmdPath := resolveToolPath(cfg.Command)
	if cmdPath == "" {
		return errors.New("linter not installed")
	}
	args := replaceFileArgs(cfg.Args, filePath)
	cmd := exec.Command(cmdPath, args...)

	if cfg := lsp.GetConfigForFile(filePath); cfg != nil {
		if root := lsp.FindProjectRoot(filePath, cfg.RootPatterns); root != "" {
			cmd.Dir = root
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return errors.New(msg)
	}
	return nil
}

func replaceFileArgs(args []string, filePath string) []string {
	out := make([]string, len(args))
	for i, arg := range args {
		out[i] = strings.ReplaceAll(arg, "{file}", filePath)
	}
	return out
}

func resolveToolPath(name string) string {
	if path, err := exec.LookPath(name); err == nil {
		return path
	}
	if name == "golangci-lint" {
		if path := findGoBinary(name); path != "" {
			return path
		}
	}
	if name == "stylua" {
		if path := findCargoBinary(name); path != "" {
			return path
		}
	}
	if npmBin := strings.TrimSpace(runCommandOutput("npm", "bin", "-g")); npmBin != "" {
		if path := filepath.Join(npmBin, name); fileExists(path) {
			return path
		}
	}
	if runtime.GOOS == "windows" {
		if path := filepath.Join(npmBinFromEnv(), name+".cmd"); fileExists(path) {
			return path
		}
	}
	return ""
}

func npmBinFromEnv() string {
	pathEnv := runCommandOutput("npm", "bin", "-g")
	return strings.TrimSpace(pathEnv)
}

func findGoBinary(name string) string {
	if runCommandOutput("go", "env", "GOPATH") == "" {
		return ""
	}
	binName := withExecutableSuffix(name)
	if gobin := strings.TrimSpace(runCommandOutput("go", "env", "GOBIN")); gobin != "" {
		if path := filepath.Join(gobin, binName); fileExists(path) {
			return path
		}
	}
	if gopath := strings.TrimSpace(runCommandOutput("go", "env", "GOPATH")); gopath != "" {
		if path := filepath.Join(gopath, "bin", binName); fileExists(path) {
			return path
		}
	}
	return ""
}

func findCargoBinary(name string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return ""
	}
	path := filepath.Join(homeDir, ".cargo", "bin", withExecutableSuffix(name))
	if fileExists(path) {
		return path
	}
	return ""
}

func withExecutableSuffix(name string) string {
	if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
		return name + ".exe"
	}
	return name
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
