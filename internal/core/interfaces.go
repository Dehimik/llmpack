package core

import (
	"io"
	"iter"
)

// Formatter визначає, як ми пакуємо файли (XML, Markdown, Zip)
type Formatter interface {
	Name() string
	Start(w io.Writer) error
	WriteTree(w io.Writer, tree string) error
	AddFile(w io.Writer, relPath string, content []byte) error
	Close(w io.Writer) error
}

// TokenCounter абстрагує підрахунок токенів
type TokenCounter interface {
	Count(text string) int
}

// Filter вирішує, чи брати файл
type Filter interface {
	ShouldIgnore(path string, isDir bool) bool
}

// Walker повертає ітератор шляхів (Go 1.23 feature)
type Walker interface {
	Walk() iter.Seq2[string, error]
}

// Config тримає налаштування запуску
type Config struct {
	RootPath        string
	OutputPath      string
	Format          string // "xml", "markdown", "zip"
	IgnoreGit       bool
	CountTokens     bool
	CopyToClipboard bool
}
