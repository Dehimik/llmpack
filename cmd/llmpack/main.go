package main

import (
	"fmt"
	"os"

	"github.com/dehimik/llmpack/internal/app"
	"github.com/dehimik/llmpack/internal/core"
	"github.com/spf13/cobra"
)

var cfg core.Config

var rootCmd = &cobra.Command{
	Use:   "llmpack [path]",
	Short: "Pack your code into LLM-friendly context",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg.InputPaths = args

		if cfg.OutputPath == "" && !cfg.CopyToClipboard {
			if cfg.Format == "markdown" {
				cfg.OutputPath = "context.md"
			} else if cfg.Format == "zip" {
				cfg.OutputPath = "context.zip"
			} else {
				cfg.OutputPath = "context.xml"
			}
		}

		if err := app.Run(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.Flags().StringVarP(&cfg.OutputPath, "output", "o", "", "Output file path (default: context.xml/.md)")
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "xml", "Output format (xml, markdown, zip, tree)")

	rootCmd.Flags().BoolVar(&cfg.IgnoreGit, "ignore-git", true, "Use .gitignore")
	rootCmd.Flags().BoolVar(&cfg.CountTokens, "tokens", true, "Count tokens")
	rootCmd.Flags().BoolVar(&cfg.NoTree, "no-tree", false, "Disable file tree in output header")
	rootCmd.Flags().BoolVarP(&cfg.CopyToClipboard, "clipboard", "c", false, "Copy output to clipboard")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
