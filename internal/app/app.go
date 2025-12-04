package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dehimik/llmpack/internal/core"
	"github.com/dehimik/llmpack/internal/formatter"
	"github.com/dehimik/llmpack/internal/tokenizer"
	"github.com/dehimik/llmpack/internal/walker"
)

// for checking if file binary(true->dont add to file)
func isBinary(content []byte) bool {
	const maxBytesToCheck = 8000
	length := len(content)
	if length > maxBytesToCheck {
		length = maxBytesToCheck
	}

	for _, b := range content[:length] {
		if b == 0 {
			return true
		}
	}
	return false
}

func Run(cfg core.Config) error {
	// Setup Formatter
	var fmtStrategy core.Formatter
	switch cfg.Format {
	case "zip":
		fmtStrategy = formatter.NewZip()
	case "markdown", "md":
		fmtStrategy = formatter.NewMarkdown()
	case "tree":
		fmtStrategy = formatter.NewTree()
	default:
		fmtStrategy = formatter.NewXML()
	}

	// Setup Walker
	wk, err := walker.New(cfg.InputPaths)
	if err != nil {
		return fmt.Errorf("failed to init walker: %w", err)
	}

	// Setup Tokenizer
	var tk *tokenizer.TikToken
	if cfg.CountTokens {
		tk = tokenizer.New()
	}

	// Output Destination Logic
	var writers []io.Writer

	if cfg.OutputPath != "" && cfg.OutputPath != "-" {
		f, err := os.Create(cfg.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		writers = append(writers, f)
	} else if cfg.OutputPath == "-" {
		writers = append(writers, os.Stdout)
	}

	var clipboardBuf *bytes.Buffer
	if cfg.CopyToClipboard {
		clipboardBuf = new(bytes.Buffer)
		writers = append(writers, clipboardBuf)
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)

	// get pretty path
	getDisplayPath := func(originalPath string) string {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, originalPath); err == nil {
				if !strings.HasPrefix(rel, "..") {
					return rel
				}
				return rel
			}
		}
		return originalPath
	}

	// Generate Tree & Collect Paths
	var files []string
	var displayPaths []string

	fmt.Println("Scanning files...")

	for path, err := range wk.Walk() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			continue
		}

		files = append(files, path)

		display := getDisplayPath(path)
		displayPaths = append(displayPaths, display)
	}

	// prepare header content
	var headerContent string

	if cfg.Format == "tree" {
		// Only if user want visual tree
		rootNode := buildTree(displayPaths)
		headerContent = renderTree(rootNode)
	} else {
		// For AI no ASCII, only clear paths
		headerContent = strings.Join(displayPaths, "\n")
	}

	// write header / start
	if err := fmtStrategy.Start(multiWriter); err != nil {
		return err
	}

	shouldWriteHeader := cfg.Format == "tree" || !cfg.NoTree

	if shouldWriteHeader {
		if err := fmtStrategy.WriteTree(multiWriter, headerContent); err != nil {
			return err
		}
	}

	// Optimization: Exit if tree-only mode
	if cfg.Format == "tree" {
		fmt.Println("Tree generated.")
		if cfg.CopyToClipboard && clipboardBuf != nil {
			if err := clipboard.WriteAll(clipboardBuf.String()); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Copied to clipboard!\n")
			}
		}
		return nil
	}

	// Process Content
	totalTokens := 0
	filesProcessed := 0

	fmt.Printf("Packing %d files...\n", len(files))

	for i, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		if isBinary(content) {
			continue
		}

		display := displayPaths[i]

		if cfg.CountTokens {
			totalTokens += tk.Count(string(content))
		}

		if err := fmtStrategy.AddFile(multiWriter, display, content); err != nil {
			return err
		}
		filesProcessed++
	}

	if err := fmtStrategy.Close(multiWriter); err != nil {
		return err
	}

	// final
	if cfg.CopyToClipboard && clipboardBuf != nil {
		if err := clipboard.WriteAll(clipboardBuf.String()); err != nil {
			fmt.Fprintf(os.Stderr, "\nFailed to copy to clipboard: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "\nCopied to clipboard!\n")
		}
	}

	// stats
	fmt.Fprintf(os.Stderr, "\nDone! Processed: %d/%d files.\n", filesProcessed, len(files))
	if cfg.CountTokens {
		fmt.Fprintf(os.Stderr, "Total Tokens: ~%d\n", totalTokens)
	}

	if cfg.OutputPath != "" && cfg.OutputPath != "-" {
		fi, _ := os.Stat(cfg.OutputPath)
		fmt.Fprintf(os.Stderr, "Created: %s (%v bytes)\n", cfg.OutputPath, fi.Size())
	}

	return nil
}
