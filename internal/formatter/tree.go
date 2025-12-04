package formatter

import "io"

type TreeFormatter struct{}

func NewTree() *TreeFormatter { return &TreeFormatter{} }

func (f *TreeFormatter) Name() string { return "tree" }

func (f *TreeFormatter) Start(w io.Writer) error {
	return nil
}

func (f *TreeFormatter) WriteTree(w io.Writer, tree string) error {
	_, err := io.WriteString(w, tree)
	return err
}

func (f *TreeFormatter) AddFile(w io.Writer, relPath string, content []byte) error {
	return nil
}

func (f *TreeFormatter) Close(w io.Writer) error {
	return nil
}
