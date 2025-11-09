package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

// FindAllFiles recursively finds all files starting from the given root directory.
// It returns a list of file paths relative to the root.
// Hidden files and directories (starting with .) are excluded.
func FindAllFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories we can't access
			return nil
		}

		// Skip hidden files and directories (starting with .)
		if info.Name() != "." && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip common build/cache directories
		if info.IsDir() {
			switch info.Name() {
			case "node_modules", "vendor", ".git", ".gocache", "dist", "build", "target":
				return filepath.SkipDir
			}
		}

		// Only include files, not directories
		if !info.IsDir() {
			// Make path relative to root
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}
			files = append(files, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
