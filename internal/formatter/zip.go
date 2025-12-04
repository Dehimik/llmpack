package formatter

import (
	"archive/zip"
	"io"
)

type ZipFormatter struct {
	zw *zip.Writer
}

func NewZip() *ZipFormatter { return &ZipFormatter{} }

func (f *ZipFormatter) Name() string { return "zip" }

func (f *ZipFormatter) Start(w io.Writer) error {
	f.zw = zip.NewWriter(w)
	return nil
}

func (f *ZipFormatter) WriteTree(w io.Writer, tree string) error {
	wr, err := f.zw.Create("project_structure.txt")
	if err != nil {
		return err
	}
	_, err = wr.Write([]byte(tree))
	return err
}

func (f *ZipFormatter) AddFile(w io.Writer, relPath string, content []byte) error {
	fRel, err := f.zw.Create(relPath)
	if err != nil {
		return err
	}
	_, err = fRel.Write(content)
	return err
}

func (f *ZipFormatter) Close(w io.Writer) error {
	return f.zw.Close()
}
