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

func (n noopIgnorer) Match(path string, isDir bool) bool {
	return false
}

type FSWalker struct {
	inputs []string
	ignore Ignorer
}

func New(inputs []string) (*FSWalker, error) {
	var matcher Ignorer = noopIgnorer{}

	if cwd, err := os.Getwd(); err == nil {
		gitIgnorePath := filepath.Join(cwd, ".gitignore")
		if _, err := os.Stat(gitIgnorePath); err == nil {
			if m, err := gitignore.NewGitIgnore(gitIgnorePath); err == nil {
				matcher = m
			}
		}
	}

	return &FSWalker{
		inputs: inputs,
		ignore: matcher,
	}, nil
}

func (w *FSWalker) Walk() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for _, inputRoot := range w.inputs {

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

			err = filepath.WalkDir(inputRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				relPath, _ := filepath.Rel(inputRoot, path)
				if relPath == "." {
					return nil
				}

				isDir := d.IsDir()

				if isDir {
					name := d.Name()
					if name == ".git" || name == "node_modules" || name == ".idea" || name == ".vscode" || name == "vendor" {
						return filepath.SkipDir
					}
				}

				if w.ignore.Match(relPath, isDir) {
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
				// Log error logic
			}
		}
	}
}
