package core

import (
	"io"
	"iter"
)

type Formatter interface {
	Name() string
	Start(w io.Writer) error
	WriteTree(w io.Writer, tree string) error
	AddFile(w io.Writer, relPath string, content []byte) error
	Close(w io.Writer) error
}

type TokenCounter interface {
	Count(text string) int
}

type Filter interface {
	ShouldIgnore(path string, isDir bool) bool
}

type Walker interface {
	Walk() iter.Seq2[string, error]
}

type Config struct {
	InputPaths      []string
	OutputPath      string
	Format          string // "xml", "markdown", "zip", etc.
	IgnoreGit       bool
	CountTokens     bool
	CopyToClipboard bool
	NoTree          bool
	SkeletonMode    bool
	IgnorePatterns  []string
	DisableSecurity bool
	ModelName       string
}
