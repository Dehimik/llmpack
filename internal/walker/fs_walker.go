package walker

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/dehimik/llmpack/internal/core"
	"github.com/monochromegane/go-gitignore"
)

type FSWalker struct {
	root   string
	ignore gitignore.IgnoreMatcher
}

func New(root string) (*FSWalker, error) {
	// Спробуємо завантажити .gitignore з кореня (спрощено для MVP)
	// У продакшені треба шукати .gitignore рекурсивно
	matcher := gitignore.NewMatcher([]gitignore.MatchPattern{})

	gitIgnorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitIgnorePath); err == nil {
		matcher, _ = gitignore.NewGitIgnore(gitIgnorePath)
	}

	return &FSWalker{
		root:   root,
		ignore: matcher,
	}, nil
}

// Walk - це імплементація Go 1.23 ітератора
func (w *FSWalker) Walk() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relPath, _ := filepath.Rel(w.root, path)
			if relPath == "." {
				return nil
			}

			// 1. Hardcoded Security Filters
			if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == ".idea") {
				return filepath.SkipDir
			}

			// 2. .gitignore check
			if w.ignore.Match(path, d.IsDir()) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				return nil
			}

			// Віддаємо шлях у цикл "for range"
			if !yield(path, nil) {
				return filepath.SkipAll
			}

			return nil
		})

		if err != nil {
			// Можна логувати помилку обходу, якщо треба
		}
	}
}
