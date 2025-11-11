package filesystem

import (
	"path/filepath"
	"strings"
)

// Nerd Font icon mappings
// Reference: https://www.nerdfonts.com/cheat-sheet
const (
	// Directories
	IconFolder     = "\uf07b" //
	IconFolderOpen = "\uf07c" //

	// Common file types
	IconFile       = "\uf15b" //
	IconMarkdown   = "\uf48a" //
	IconGo         = "\ue627" //
	IconJavaScript = "\ue74e" //
	IconTypeScript = "\ue628" //
	IconPython     = "\ue73c" //
	IconRust       = "\ue7a8" //
	IconC          = "\ue61e" //
	IconCPlusPlus  = "\ue61d" //
	IconJava       = "\ue738" //
	IconHTML       = "\ue736" //
	IconCSS        = "\ue749" //
	IconJSON       = "\ue60b" //
	IconXML        = "\ue619" //
	IconYAML       = "\uf481" //
	IconToml       = "\ue615" //
	IconSQL        = "\uf472" //
	IconDocker     = "\uf308" //
	IconGit        = "\ue702" //
	IconImage      = "\uf1c5" //
	IconVideo      = "\uf03d" //
	IconAudio      = "\uf001" //
	IconPDF        = "\uf1c1" //
	IconZip        = "\uf410" //
	IconText       = "\uf15c" //
	IconConfig     = "\ue615" //
	IconLog        = "\uf15c" //
	IconLock       = "\uf023" //
	IconKey        = "\uf084" //
	IconDatabase   = "\uf1c0" //
	IconShell      = "\uf489" //
	IconVim        = "\ue7c5" //
	IconBird       = "\uedea" //  Norse hammer icon for .crl Carrion files
)

// GetFileIcon returns the Nerd Font icon for a given file or directory
func GetFileIcon(name string, isDir bool) string {
	if isDir {
		return IconFolder
	}

	// Get extension (lowercase)
	ext := strings.ToLower(filepath.Ext(name))

	// Check for exact filename matches first (highest priority)
	switch strings.ToLower(name) {
	case "dockerfile", "dockerfile.dev", "dockerfile.prod":
		return IconDocker
	case ".gitignore", ".gitattributes", ".gitmodules":
		return IconGit
	case "makefile", "makefile.am":
		return "\uf489" //
	case "readme.md", "readme", "readme.txt":
		return IconMarkdown
	case "license", "license.md", "license.txt":
		return "\uf48a" //
	case "package.json", "package-lock.json":
		return "\ue71e" //
	case "go.mod", "go.sum":
		return IconGo
	case "cargo.toml", "cargo.lock":
		return IconRust
	case "requirements.txt", "setup.py":
		return IconPython
	case ".env", ".env.local", ".env.production":
		return IconConfig
	}

	// Extension-based matching
	switch ext {
	// Programming languages
	case ".go":
		return IconGo
	case ".js", ".mjs", ".cjs":
		return IconJavaScript
	case ".ts", ".tsx":
		return IconTypeScript
	case ".py", ".pyw", ".pyx":
		return IconPython
	case ".rs":
		return IconRust
	case ".c":
		return IconC
	case ".cpp", ".cc", ".cxx", ".hpp", ".h":
		return IconCPlusPlus
	case ".java":
		return IconJava

	// Custom file types
	case ".crl":
		return IconBird

	// Web
	case ".html", ".htm":
		return IconHTML
	case ".css", ".scss", ".sass", ".less":
		return IconCSS
	case ".vue":
		return "\ufd42" //
	case ".jsx":
		return "\ue7ba" //
	case ".php":
		return "\ue73d" //

	// Data/Config
	case ".json", ".jsonc":
		return IconJSON
	case ".xml":
		return IconXML
	case ".yaml", ".yml":
		return IconYAML
	case ".toml":
		return IconToml
	case ".ini", ".cfg", ".conf":
		return IconConfig
	case ".sql":
		return IconSQL

	// Markup/Documentation
	case ".md", ".markdown":
		return IconMarkdown
	case ".txt", ".text":
		return IconText
	case ".pdf":
		return IconPDF
	case ".doc", ".docx":
		return "\uf1c2" //

	// Images
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp", ".bmp":
		return IconImage

	// Video
	case ".mp4", ".avi", ".mov", ".mkv", ".webm":
		return IconVideo

	// Audio
	case ".mp3", ".wav", ".flac", ".ogg", ".m4a":
		return IconAudio

	// Archives
	case ".zip", ".tar", ".gz", ".bz2", ".7z", ".rar":
		return IconZip

	// Shell scripts
	case ".sh", ".bash", ".zsh", ".fish":
		return IconShell

	// Vim
	case ".vim", ".vimrc":
		return IconVim

	// Git
	case ".git":
		return IconGit

	// Logs
	case ".log":
		return IconLog

	// Lock files
	case ".lock":
		return IconLock

	// Keys/Certificates
	case ".pem", ".key", ".crt", ".cer", ".p12":
		return IconKey

	// Database
	case ".db", ".sqlite", ".sqlite3":
		return IconDatabase

	// Binary/Executable
	case ".exe", ".bin", ".dll", ".so", ".dylib":
		return "\uf489" //

	default:
		return IconFile
	}
}

// GetExpandIcon returns the expand/collapse icon for directories
func GetExpandIcon(expanded bool) string {
	if expanded {
		return "\uf078" //  (chevron-down)
	}
	return "\uf054" //  (chevron-right)
}
