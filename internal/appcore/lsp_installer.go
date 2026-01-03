package appcore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/javanhut/vem/internal/lsp"
)

type lspInstallSpec struct {
	name              string
	npmPackages       []string
	pipPackage        string
	brewPackages      []string
	aptPackages       []string
	dnfPackages       []string
	pacmanPackages    []string
	zypperPackages    []string
	wingetPackage     string
	chocoPackage      string
	scoopPackage      string
	goInstallModule   string
	rustupComponent   string
	preferSystemPkg   bool
	requiresSudoLinux bool
}

type installCommand struct {
	Name         string
	Args         []string
	RequiresSudo bool
	Fallbacks    []installCommand
}

func resolveLSPInstallCommand(cfg *lsp.ServerConfig) (*installCommand, error) {
	if cfg == nil {
		return nil, fmt.Errorf("missing config")
	}

	if len(cfg.InstallCommand) > 0 && hasCommand(cfg.InstallCommand[0]) {
		return &installCommand{
			Name: cfg.InstallCommand[0],
			Args: append([]string(nil), cfg.InstallCommand[1:]...),
		}, nil
	}

	spec := lspInstallSpecForConfig(cfg)
	if spec.name == "" {
		return nil, fmt.Errorf("no install spec for %s", cfg.Name)
	}

	if spec.rustupComponent != "" && hasCommand("rustup") {
		return &installCommand{
			Name: "rustup",
			Args: []string{"component", "add", spec.rustupComponent},
		}, nil
	}

	if spec.goInstallModule != "" && hasCommand("go") {
		return &installCommand{
			Name: "go",
			Args: []string{"install", spec.goInstallModule},
		}, nil
	}

	if len(spec.npmPackages) > 0 {
		if cmd, ok := npmInstallCommand(spec.npmPackages); ok {
			return cmd, nil
		}
	}

	if spec.pipPackage != "" {
		if cmd, ok := pipInstallCommand(spec.pipPackage); ok {
			return cmd, nil
		}
	}

	if cmd, ok := systemPackageInstallCommand(spec); ok {
		return cmd, nil
	}

	return nil, fmt.Errorf("no supported package manager available")
}

func resolveLSPUninstallCommand(cfg *lsp.ServerConfig) (*installCommand, error) {
	if cfg == nil {
		return nil, fmt.Errorf("missing config")
	}

	spec := lspInstallSpecForConfig(cfg)
	if spec.name == "" {
		return nil, fmt.Errorf("no uninstall spec for %s", cfg.Name)
	}

	if spec.rustupComponent != "" && hasCommand("rustup") {
		return &installCommand{
			Name: "rustup",
			Args: []string{"component", "remove", spec.rustupComponent},
		}, nil
	}

	if spec.goInstallModule != "" && hasCommand("go") {
		return &installCommand{
			Name: "go",
			Args: []string{"clean", "-i", spec.goInstallModule},
		}, nil
	}

	if len(spec.npmPackages) > 0 {
		if cmd, ok := npmUninstallCommand(spec.npmPackages); ok {
			return cmd, nil
		}
	}

	if spec.pipPackage != "" {
		if cmd, ok := pipUninstallCommand(spec.pipPackage); ok {
			return cmd, nil
		}
	}

	if cmd, ok := systemPackageUninstallCommand(spec); ok {
		return cmd, nil
	}

	return nil, fmt.Errorf("no supported package manager available")
}

func resolveToolInstallCommand(spec toolInstallSpec) (*installCommand, error) {
	if spec.rustupComponent != "" && hasCommand("rustup") {
		return &installCommand{Name: "rustup", Args: []string{"component", "add", spec.rustupComponent}}, nil
	}
	if spec.goInstall != "" && hasCommand("go") {
		return &installCommand{Name: "go", Args: []string{"install", spec.goInstall}}, nil
	}
	if spec.cargoCrate != "" && hasCommand("cargo") {
		return &installCommand{Name: "cargo", Args: []string{"install", spec.cargoCrate}}, nil
	}
	if len(spec.npmPackages) > 0 {
		if cmd, ok := npmInstallCommand(spec.npmPackages); ok {
			return cmd, nil
		}
	}
	if spec.pipPackage != "" {
		if cmd, ok := pipInstallCommand(spec.pipPackage); ok {
			return cmd, nil
		}
	}
	if spec.luarocksPackage != "" && hasCommand("luarocks") {
		return &installCommand{Name: "luarocks", Args: []string{"install", spec.luarocksPackage}}, nil
	}
	if cmd, ok := systemPackageInstallCommand(lspInstallSpec{
		brewPackages:   spec.brewPackages,
		aptPackages:    spec.aptPackages,
		dnfPackages:    spec.dnfPackages,
		pacmanPackages: spec.pacmanPackages,
		zypperPackages: spec.zypperPackages,
	}); ok {
		return cmd, nil
	}
	return nil, fmt.Errorf("no supported package manager available")
}

func resolveToolUninstallCommand(spec toolInstallSpec) (*installCommand, error) {
	if spec.rustupComponent != "" && hasCommand("rustup") {
		return &installCommand{Name: "rustup", Args: []string{"component", "remove", spec.rustupComponent}}, nil
	}
	if spec.goInstall != "" && hasCommand("go") {
		return &installCommand{Name: "go", Args: []string{"clean", "-i", spec.goInstall}}, nil
	}
	if spec.cargoCrate != "" && hasCommand("cargo") {
		return &installCommand{Name: "cargo", Args: []string{"uninstall", spec.cargoCrate}}, nil
	}
	if len(spec.npmPackages) > 0 {
		if cmd, ok := npmUninstallCommand(spec.npmPackages); ok {
			return cmd, nil
		}
	}
	if spec.pipPackage != "" {
		if cmd, ok := pipUninstallCommand(spec.pipPackage); ok {
			return cmd, nil
		}
	}
	if spec.luarocksPackage != "" && hasCommand("luarocks") {
		return &installCommand{Name: "luarocks", Args: []string{"remove", spec.luarocksPackage}}, nil
	}
	if cmd, ok := systemPackageUninstallCommand(lspInstallSpec{
		brewPackages:   spec.brewPackages,
		aptPackages:    spec.aptPackages,
		dnfPackages:    spec.dnfPackages,
		pacmanPackages: spec.pacmanPackages,
		zypperPackages: spec.zypperPackages,
	}); ok {
		return cmd, nil
	}
	return nil, fmt.Errorf("no supported package manager available")
}

func lspInstallSpecForConfig(cfg *lsp.ServerConfig) lspInstallSpec {
	name := strings.ToLower(cfg.Name)

	switch name {
	case "gopls":
		return lspInstallSpec{
			name:            "gopls",
			goInstallModule: "golang.org/x/tools/gopls@latest",
			brewPackages:    []string{"gopls"},
			aptPackages:     []string{"golang-gopls", "gopls"},
			dnfPackages:     []string{"golang-gopls", "gopls"},
			pacmanPackages:  []string{"gopls"},
			zypperPackages:  []string{"gopls"},
		}
	case "pyright":
		return lspInstallSpec{
			name:         "pyright",
			npmPackages:  []string{"pyright"},
			brewPackages: []string{"pyright"},
		}
	case "python-lsp-server":
		return lspInstallSpec{
			name:       "python-lsp-server",
			pipPackage: "python-lsp-server",
		}
	case "typescript-language-server":
		return lspInstallSpec{
			name:        "typescript-language-server",
			npmPackages: []string{"typescript-language-server", "typescript"},
		}
	case "rust-analyzer":
		return lspInstallSpec{
			name:            "rust-analyzer",
			rustupComponent: "rust-analyzer",
			brewPackages:    []string{"rust-analyzer"},
			aptPackages:     []string{"rust-analyzer"},
			dnfPackages:     []string{"rust-analyzer"},
			pacmanPackages:  []string{"rust-analyzer"},
			zypperPackages:  []string{"rust-analyzer"},
		}
	case "clangd":
		return lspInstallSpec{
			name:           "clangd",
			brewPackages:   []string{"llvm", "clangd"},
			aptPackages:    []string{"clangd", "clang-tools-extra"},
			dnfPackages:    []string{"clangd", "clang-tools-extra"},
			pacmanPackages: []string{"clang", "clang-tools-extra", "clangd"},
			zypperPackages: []string{"clangd", "clang-tools-extra"},
			wingetPackage:  "LLVM.LLVM",
			chocoPackage:   "llvm",
			scoopPackage:   "llvm",
		}
	case "lua-language-server":
		return lspInstallSpec{
			name:           "lua-language-server",
			brewPackages:   []string{"lua-language-server"},
			aptPackages:    []string{"lua-language-server", "lua-lsp"},
			dnfPackages:    []string{"lua-language-server", "lua-lsp"},
			pacmanPackages: []string{"lua-language-server"},
			zypperPackages: []string{"lua-language-server", "lua-lsp"},
		}
	case "vscode-html-language-server":
		return lspInstallSpec{
			name:        "vscode-html-language-server",
			npmPackages: []string{"vscode-langservers-extracted"},
		}
	default:
		return lspInstallSpec{}
	}
}

func hasCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func npmInstallCommand(pkgs []string) (*installCommand, bool) {
	if hasCommand("npm") {
		return &installCommand{Name: "npm", Args: append([]string{"install", "-g"}, pkgs...)}, true
	}
	if hasCommand("pnpm") {
		return &installCommand{Name: "pnpm", Args: append([]string{"add", "-g"}, pkgs...)}, true
	}
	if hasCommand("yarn") {
		return &installCommand{Name: "yarn", Args: append([]string{"global", "add"}, pkgs...)}, true
	}
	return nil, false
}

func npmUninstallCommand(pkgs []string) (*installCommand, bool) {
	if hasCommand("npm") {
		return &installCommand{Name: "npm", Args: append([]string{"uninstall", "-g"}, pkgs...)}, true
	}
	if hasCommand("pnpm") {
		return &installCommand{Name: "pnpm", Args: append([]string{"remove", "-g"}, pkgs...)}, true
	}
	if hasCommand("yarn") {
		return &installCommand{Name: "yarn", Args: append([]string{"global", "remove"}, pkgs...)}, true
	}
	return nil, false
}

func pipInstallCommand(pkg string) (*installCommand, bool) {
	if hasCommand("pipx") {
		return &installCommand{Name: "pipx", Args: []string{"install", pkg}}, true
	}
	if hasCommand("pip3") {
		return &installCommand{Name: "pip3", Args: []string{"install", pkg}}, true
	}
	if hasCommand("pip") {
		return &installCommand{Name: "pip", Args: []string{"install", pkg}}, true
	}
	return nil, false
}

func pipUninstallCommand(pkg string) (*installCommand, bool) {
	if hasCommand("pipx") {
		return &installCommand{Name: "pipx", Args: []string{"uninstall", pkg}}, true
	}
	if hasCommand("pip3") {
		return &installCommand{Name: "pip3", Args: []string{"uninstall", "-y", pkg}}, true
	}
	if hasCommand("pip") {
		return &installCommand{Name: "pip", Args: []string{"uninstall", "-y", pkg}}, true
	}
	return nil, false
}

func systemPackageInstallCommand(spec lspInstallSpec) (*installCommand, bool) {
	switch runtime.GOOS {
	case "darwin":
		if len(spec.brewPackages) > 0 && hasCommand("brew") {
			cmd := &installCommand{Name: "brew", Args: []string{"install", spec.brewPackages[0]}}
			for i := 1; i < len(spec.brewPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name: "brew",
					Args: []string{"install", spec.brewPackages[i]},
				})
			}
			return cmd, true
		}
	case "windows":
		if spec.wingetPackage != "" && hasCommand("winget") {
			return &installCommand{Name: "winget", Args: []string{"install", "--id", spec.wingetPackage}}, true
		}
		if spec.chocoPackage != "" && hasCommand("choco") {
			return &installCommand{Name: "choco", Args: []string{"install", "-y", spec.chocoPackage}}, true
		}
		if spec.scoopPackage != "" && hasCommand("scoop") {
			return &installCommand{Name: "scoop", Args: []string{"install", spec.scoopPackage}}, true
		}
	default:
		if len(spec.aptPackages) > 0 && hasCommand("apt-get") {
			cmd := &installCommand{Name: "apt-get", Args: []string{"install", "-y", spec.aptPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.aptPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "apt-get",
					Args:         []string{"install", "-y", spec.aptPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.dnfPackages) > 0 && hasCommand("dnf") {
			cmd := &installCommand{Name: "dnf", Args: []string{"install", "-y", spec.dnfPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.dnfPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "dnf",
					Args:         []string{"install", "-y", spec.dnfPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.pacmanPackages) > 0 && hasCommand("pacman") {
			cmd := &installCommand{Name: "pacman", Args: []string{"-S", "--noconfirm", spec.pacmanPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.pacmanPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "pacman",
					Args:         []string{"-S", "--noconfirm", spec.pacmanPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.zypperPackages) > 0 && hasCommand("zypper") {
			cmd := &installCommand{Name: "zypper", Args: []string{"--non-interactive", "install", spec.zypperPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.zypperPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "zypper",
					Args:         []string{"--non-interactive", "install", spec.zypperPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
	}
	return nil, false
}

func systemPackageUninstallCommand(spec lspInstallSpec) (*installCommand, bool) {
	switch runtime.GOOS {
	case "darwin":
		if len(spec.brewPackages) > 0 && hasCommand("brew") {
			cmd := &installCommand{Name: "brew", Args: []string{"uninstall", spec.brewPackages[0]}}
			for i := 1; i < len(spec.brewPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name: "brew",
					Args: []string{"uninstall", spec.brewPackages[i]},
				})
			}
			return cmd, true
		}
	case "windows":
		if spec.wingetPackage != "" && hasCommand("winget") {
			return &installCommand{Name: "winget", Args: []string{"uninstall", "--id", spec.wingetPackage}}, true
		}
		if spec.chocoPackage != "" && hasCommand("choco") {
			return &installCommand{Name: "choco", Args: []string{"uninstall", "-y", spec.chocoPackage}}, true
		}
		if spec.scoopPackage != "" && hasCommand("scoop") {
			return &installCommand{Name: "scoop", Args: []string{"uninstall", spec.scoopPackage}}, true
		}
	default:
		if len(spec.aptPackages) > 0 && hasCommand("apt-get") {
			cmd := &installCommand{Name: "apt-get", Args: []string{"remove", "-y", spec.aptPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.aptPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "apt-get",
					Args:         []string{"remove", "-y", spec.aptPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.dnfPackages) > 0 && hasCommand("dnf") {
			cmd := &installCommand{Name: "dnf", Args: []string{"remove", "-y", spec.dnfPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.dnfPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "dnf",
					Args:         []string{"remove", "-y", spec.dnfPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.pacmanPackages) > 0 && hasCommand("pacman") {
			cmd := &installCommand{Name: "pacman", Args: []string{"-Rns", "--noconfirm", spec.pacmanPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.pacmanPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "pacman",
					Args:         []string{"-Rns", "--noconfirm", spec.pacmanPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
		if len(spec.zypperPackages) > 0 && hasCommand("zypper") {
			cmd := &installCommand{Name: "zypper", Args: []string{"--non-interactive", "remove", spec.zypperPackages[0]}, RequiresSudo: true}
			for i := 1; i < len(spec.zypperPackages); i++ {
				cmd.Fallbacks = append(cmd.Fallbacks, installCommand{
					Name:         "zypper",
					Args:         []string{"--non-interactive", "remove", spec.zypperPackages[i]},
					RequiresSudo: true,
				})
			}
			return cmd, true
		}
	}
	return nil, false
}

func maybeSudo(cmd *installCommand) *installCommand {
	if cmd == nil || !cmd.RequiresSudo {
		return cmd
	}
	rewrapped := &installCommand{}
	if runtime.GOOS != "windows" && hasCommand("pkexec") {
		rewrapped.Name = "pkexec"
		rewrapped.Args = append([]string{cmd.Name}, cmd.Args...)
	}
	if rewrapped.Name == "" && hasCommand("sudo") {
		rewrapped.Name = "sudo"
		rewrapped.Args = append([]string{cmd.Name}, cmd.Args...)
	}
	if rewrapped.Name == "" {
		return cmd
	}
	if len(cmd.Fallbacks) > 0 {
		for _, fb := range cmd.Fallbacks {
			wrapped := installCommand{Name: rewrapped.Name, Args: append([]string{fb.Name}, fb.Args...)}
			rewrapped.Fallbacks = append(rewrapped.Fallbacks, wrapped)
		}
	}
	return rewrapped
}

func ensureGoBinInPath() {
	if !hasCommand("go") {
		return
	}
	goBin := strings.TrimSpace(runCommandOutput("go", "env", "GOBIN"))
	if goBin == "" {
		goPath := strings.TrimSpace(runCommandOutput("go", "env", "GOPATH"))
		if goPath != "" {
			goBin = filepath.Join(goPath, "bin")
		}
	}
	if goBin == "" {
		return
	}
	appendToPath(goBin)
}

func ensureCargoBinInPath() {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}
	appendToPath(filepath.Join(home, ".cargo", "bin"))
}

func appendToPath(dir string) {
	if dir == "" {
		return
	}
	pathEnv := os.Getenv("PATH")
	for _, part := range strings.Split(pathEnv, string(os.PathListSeparator)) {
		if part == dir {
			return
		}
	}
	if pathEnv == "" {
		os.Setenv("PATH", dir)
		return
	}
	os.Setenv("PATH", pathEnv+string(os.PathListSeparator)+dir)
}

func runCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}
