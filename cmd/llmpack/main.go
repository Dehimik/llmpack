package main

import (
	"fmt"
	"os"

	"github.com/dehimik/llmpack/internal/app"
	"github.com/dehimik/llmpack/internal/config"
	"github.com/dehimik/llmpack/internal/core"
	"github.com/spf13/cobra"
)

var (
	cfg         core.Config
	profileName string
	configPath  string
)

func hasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

var rootCmd = &cobra.Command{
	Use:   "llmpack [path]",
	Short: "Pack your code into LLM-friendly context",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && !hasStdinData() {
			return fmt.Errorf("requires at least 1 arg OR data from stdin")
		}
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fileCfg, err := config.Load(configPath)
		if err != nil {
			if !os.IsNotExist(err) && err.Error() != "config file not found" {
				fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
			}
		}

		settings := fileCfg.Global
		if profileName != "" {
			if p, ok := fileCfg.Profiles[profileName]; ok {
				settings = p
			} else {
				fmt.Fprintf(os.Stderr, "Warning: Profile '%s' not found in config, using global settings.\n", profileName)
			}
		}

		cfg.IgnorePatterns = fileCfg.Ignore

		if !cmd.Flags().Changed("format") && settings.Format != "" {
			cfg.Format = settings.Format
		}

		if !cmd.Flags().Changed("ignore-git") {
			cfg.IgnoreGit = settings.IgnoreGit
		}

		if !cmd.Flags().Changed("tokens") {
			cfg.CountTokens = settings.Tokens
		}

		if !cmd.Flags().Changed("model") {
			cfg.ModelName = settings.ModelName
		}

		if !cmd.Flags().Changed("skeleton") {
			cfg.SkeletonMode = settings.SkeletonMode
		}

		if !cmd.Flags().Changed("no-tree") {
			cfg.NoTree = settings.NoTree
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg.InputPaths = args

		if cfg.OutputPath == "" && !cfg.CopyToClipboard {
			if cfg.Format == "markdown" || cfg.Format == "md" {
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
	rootCmd.Flags().StringVarP(&cfg.OutputPath, "output", "o", "", "Output file path")
	rootCmd.Flags().StringVarP(&cfg.Format, "format", "f", "xml", "Output format (xml, markdown, zip, tree)")

	rootCmd.Flags().BoolVar(&cfg.IgnoreGit, "ignore-git", true, "Use .gitignore")
	rootCmd.Flags().BoolVar(&cfg.CountTokens, "tokens", true, "Count tokens")
	rootCmd.Flags().StringVarP(&cfg.ModelName, "model", "m", "gpt-4o", "Model for cost estimation (gpt-4o, claude-3-5-sonnet, etc.)")

	rootCmd.Flags().BoolVar(&cfg.NoTree, "no-tree", false, "Disable file tree in output header")
	rootCmd.Flags().BoolVarP(&cfg.CopyToClipboard, "clipboard", "c", false, "Copy output to clipboard")
	rootCmd.Flags().BoolVar(&cfg.DisableSecurity, "no-security", false, "Disable security checks (secrets detection)")

	rootCmd.Flags().BoolVarP(&cfg.SkeletonMode, "skeleton", "s", false, "Strip function bodies (skeleton mode)")
	rootCmd.Flags().StringVarP(&profileName, "profile", "p", "", "Configuration profile to use (defined in .llmpack.yaml)")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config file (default .llmpack.yaml)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
