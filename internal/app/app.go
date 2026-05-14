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
	"github.com/dehimik/llmpack/internal/pricing"
	"github.com/dehimik/llmpack/internal/security"
	"github.com/dehimik/llmpack/internal/skeleton"
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

func isPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func Run(cfg core.Config) error {
	// Setup Log Writer
	logWriter := io.Writer(os.Stderr)
	if cfg.LogWriter != nil {
		logWriter = cfg.LogWriter
	}

	// Setup Formatter
	var fmtStrategy core.Formatter
	secScanner := security.New(cfg.DisableSecurity)
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
	wk, err := walker.New(cfg.InputPaths, cfg.IgnorePatterns)
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

	if cfg.CustomWriter != nil {
		writers = append(writers, cfg.CustomWriter)
	} else if cfg.OutputPath != "" && cfg.OutputPath != "-" {
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

	totalTokens := 0
	filesProcessed := 0

	if isPiped() {
		fmt.Println("Reading from STDIN...")

		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		if len(content) > 0 {
			// Security Check
			if err := secScanner.Scan("stdin_input", content); err != nil {
				fmt.Fprintf(logWriter, "SECURITY WARNING: Skipping STDIN -> %v\n", err)
			} else {
				if !isBinary(content) {
					if cfg.CountTokens {
						totalTokens += tk.Count(string(content))
					}
					if err := fmtStrategy.AddFile(multiWriter, "STDIN", content); err != nil {
						return err
					}
					fmt.Fprintf(logWriter, "Added content from STDIN (%d bytes)\n", len(content))
				}
			}
		}
	}

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
			fmt.Fprintf(logWriter, "Error accessing %s: %v\n", path, err)
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

	// Filter by symbol if --find is specified
	if cfg.FindSymbol != "" {
		fmt.Printf("Searching for symbol '%s'...\n", cfg.FindSymbol)
		var filteredFiles []string
		var filteredDisplayPaths []string

		for i, path := range files {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			if isBinary(content) {
				continue
			}

			symbols, _ := skeleton.ExtractSymbols(path, content)
			found := false
			for _, s := range symbols {
				qualified := s.Name
				if s.Parent != "" {
					qualified = s.Parent + "." + s.Name
				}
				if s.Name == cfg.FindSymbol || qualified == cfg.FindSymbol {
					found = true
					break
				}
			}

			if found {
				filteredFiles = append(filteredFiles, path)
				filteredDisplayPaths = append(filteredDisplayPaths, displayPaths[i])
			}
		}
		files = filteredFiles
		displayPaths = filteredDisplayPaths

		// Update header content if we filtered
		if cfg.Format != "tree" {
			headerContent = strings.Join(displayPaths, "\n")
		}
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
				fmt.Fprintf(logWriter, "Failed to copy to clipboard: %v\n", err)
			} else {
				fmt.Fprintf(logWriter, "Copied to clipboard!\n")
			}
		}
		return nil
	}

	// Process Content

	fmt.Printf("Packing %d files...\n", len(files))

	for i, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		// 1. Binary Check
		if isBinary(content) {
			continue
		}

		// 2. Security Check (До всього іншого)
		if err := secScanner.Scan(path, content); err != nil {
			fmt.Fprintf(logWriter, "SECURITY WARNING: Skipping %s -> %v\n", path, err)
			continue
		}

		display := displayPaths[i]

		// 3. Handle --symbols mode
		if cfg.SymbolsOnly {
			symbols, _ := skeleton.ExtractSymbols(path, content)
			if len(symbols) == 0 {
				continue // Skip files with no symbols in symbols-only mode
			}
			var symBuf bytes.Buffer
			symBuf.WriteString("Symbols in this file:\n")
			for _, s := range symbols {
				scope := ""
				if s.Parent != "" {
					scope = " scope=\"" + s.Parent + "\""
				}
				fmt.Fprintf(&symBuf, "  <symbol name=\"%s\" type=\"%s\"%s lines=\"%d-%d\" />\n", s.Name, s.Type, scope, s.StartLine, s.EndLine)
			}
			content = symBuf.Bytes()
		} else {
			// 4. Targeted Implementation Extraction (Focus)
			target := ""
			if cfg.Implementation != "" {
				target = cfg.Implementation
			} else if cfg.FindSymbol != "" && cfg.Focus {
				target = cfg.FindSymbol
			}

			if target != "" {
				reduced, err := skeleton.ProcessSpecific(path, content, target)
				if err != nil {
					fmt.Fprintf(logWriter, "Warning: failed to focus on %s in %s: %v\n", target, path, err)
				} else {
					content = reduced
				}
			} else if cfg.SkeletonMode {
				// 5. Standard Skeleton Mode
				reduced, err := skeleton.Process(path, content)
				if err != nil {
					fmt.Fprintf(logWriter, "Warning: failed to skeletonize %s: %v\n", path, err)
				} else {
					content = reduced
				}
			}
		}

		// 6. Token Counting (Тільки ОДИН раз, після всіх модифікацій)
		if cfg.CountTokens {
			totalTokens += tk.Count(string(content))
		}

		// 7. Write Output
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
			fmt.Fprintf(logWriter, "\nFailed to copy to clipboard: %v\n", err)
		} else {
			fmt.Fprintf(logWriter, "\nCopied to clipboard!\n")
		}
	}

	// stats
	fmt.Fprintf(logWriter, "\nDone! Processed: %d/%d files.\n", filesProcessed, len(files))
	if cfg.CountTokens {
		costStr := pricing.Estimate(totalTokens, cfg.ModelName)
		fmt.Fprintf(logWriter, "Total Tokens: ~%d (%s for %s)\n", totalTokens, costStr, cfg.ModelName)
	}

	if cfg.OutputPath != "" && cfg.OutputPath != "-" {
		fi, _ := os.Stat(cfg.OutputPath)
		fmt.Fprintf(logWriter, "Created: %s (%v bytes)\n", cfg.OutputPath, fi.Size())
	}

	return nil
}
