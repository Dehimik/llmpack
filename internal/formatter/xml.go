package formatter

import (
	"fmt"
	"io"
	"strings"
)

type XMLFormatter struct{}

func NewXML() *XMLFormatter { return &XMLFormatter{} }

func (f *XMLFormatter) Name() string { return "xml" }

func (f *XMLFormatter) Start(w io.Writer) error {
	_, err := io.WriteString(w, "<project_context>\n")
	return err
}

func (f *XMLFormatter) WriteTree(w io.Writer, tree string) error {
	_, err := fmt.Fprintf(w, "<file_tree>\n%s\n</file_tree>\n\n", tree)
	return err
}

func (f *XMLFormatter) AddFile(w io.Writer, relPath string, content []byte) error {
	cleanContent := strings.ReplaceAll(string(content), "]]>", "]]]]><![CDATA[>")

	header := fmt.Sprintf("<file path=\"%s\">\n<![CDATA[\n", relPath)
	footer := "\n]]>\n</file>\n\n"

	if _, err := io.WriteString(w, header); err != nil {
		return err
	}
	if _, err := w.Write([]byte(cleanContent)); err != nil {
		return err
	}
	if _, err := io.WriteString(w, footer); err != nil {
		return err
	}

	return nil
}

func (f *XMLFormatter) Close(w io.Writer) error {
	_, err := io.WriteString(w, "</project_context>")
	return err
}
