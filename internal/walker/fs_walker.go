package walker

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"

	"github.com/monochromegane/go-gitignore"
)

type Ignorer interface {
	Match(path string, isDir bool) bool
}

type noopIgnorer struct{}

func (n noopIgnorer) Match(path string, isDir bool) bool { return false }

type FSWalker struct {
	inputs         []string
	ignorePatterns []string
}

func New(inputs []string, ignorePatterns []string) (*FSWalker, error) {
	return &FSWalker{
		inputs:         inputs,
		ignorePatterns: ignorePatterns,
	}, nil
}

func (w *FSWalker) Walk() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for _, inputRoot := range w.inputs {
			var ignoreMatcher Ignorer = noopIgnorer{}

			info, err := os.Stat(inputRoot)
			if err != nil {
				if !yield(inputRoot, err) {
					return
				}
				continue
			}

			if !info.IsDir() {
				if !yield(inputRoot, nil) {
					return
				}
				continue
			}

			gitIgnorePath := filepath.Join(inputRoot, ".gitignore")
			if _, err := os.Stat(gitIgnorePath); err == nil {
				if m, err := gitignore.NewGitIgnore(gitIgnorePath); err == nil {
					ignoreMatcher = m
				}
			}

			// scan
			err = filepath.WalkDir(inputRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				relPath, _ := filepath.Rel(inputRoot, path)
				if relPath == "." {
					return nil
				}

				isDir := d.IsDir()
				name := d.Name()

				for _, pattern := range w.ignorePatterns {
					if name == pattern {
						if isDir {
							return filepath.SkipDir
						}
						return nil
					}
					if matched, _ := filepath.Match(pattern, name); matched {
						if isDir {
							return filepath.SkipDir
						}
						return nil
					}
				}

				// 1. Hardcoded Security Filters
				if isDir {
					name := d.Name()
					if name == ".git" || name == "node_modules" || name == ".idea" || name == ".vscode" || name == "vendor" || name == "dist" || name == "build" {
						return filepath.SkipDir
					}
				}

				// 2. .gitignore Check
				if ignoreMatcher.Match(relPath, isDir) {
					if isDir {
						return filepath.SkipDir
					}
					return nil
				}

				if isDir {
					return nil
				}

				if !yield(path, nil) {
					return filepath.SkipAll
				}

				return nil
			})

			if err != nil {
				// Log error logic if needed
			}
		}
	}
}
