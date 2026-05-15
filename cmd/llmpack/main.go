package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dehimik/llmpack/internal/app"
	"github.com/dehimik/llmpack/internal/config"
	"github.com/dehimik/llmpack/internal/core"
	"github.com/spf13/cobra"
)

var (
	cfg         core.Config
	profileName string
	configPath  string
	appRun      = app.Run
)

func hasStdinData() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

var rootCmd = &cobra.Command{
	Use:   "llmpack",
	Short: "Pack your code into LLM-friendly context",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

		if cfg.SymbolsOnly && cfg.Implementation != "" {
			return fmt.Errorf("--symbols and --implementation cannot both be set")
		}
		return nil
	},
}

var packCmd = &cobra.Command{
	Use:   "pack [path]",
	Short: "Pack files into context (default)",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		packRun(cmd, args)
	},
}

func packRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasStdinData() {
		_ = cmd.Help()
		return
	}

	cfg.InputPaths = args

	if cfg.OutputPath == "" && !cfg.CopyToClipboard && len(args) > 0 {
		if cfg.Format == "markdown" || cfg.Format == "md" {
			cfg.OutputPath = "context.md"
		} else if cfg.Format == "zip" {
			cfg.OutputPath = "context.zip"
		} else if cfg.Format == "llms-txt" || cfg.Format == "llms" {
			cfg.OutputPath = "llms.txt"
		} else if cfg.Format == "tree" {
			cfg.OutputPath = "-"
		} else {
			cfg.OutputPath = "context.xml"
		}
	}

	if err := appRun(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setupFlags() {
	// Flags are now persistent so they apply to both 'pack' and 'root' (and thus subcommands)
	rootCmd.PersistentFlags().StringVarP(&cfg.OutputPath, "output", "o", "", "Output file path")
	rootCmd.PersistentFlags().StringVarP(&cfg.Format, "format", "f", "xml", "Output format (xml, markdown, zip, tree, llms-txt)")

	// Autocompletion for formats
	_ = rootCmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"xml", "markdown", "md", "llms-txt", "llms", "zip", "tree"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.PersistentFlags().BoolVar(&cfg.IgnoreGit, "ignore-git", true, "Use .gitignore")
	rootCmd.PersistentFlags().BoolVar(&cfg.CountTokens, "tokens", true, "Count tokens")
	rootCmd.PersistentFlags().StringVarP(&cfg.ModelName, "model", "m", "gpt-4o", "Model for cost estimation (gpt-4o, claude-3-5-sonnet, etc.)")

	// Autocompletion for models
	_ = rootCmd.RegisterFlagCompletionFunc("model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"gpt-4o", "gpt-4-turbo", "gpt-3.5-turbo", "claude-3-5-sonnet", "claude-3-opus", "gemini-1.5-pro", "gemini-1.5-flash"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootCmd.PersistentFlags().BoolVar(&cfg.NoTree, "no-tree", false, "Disable file tree in output header")
	rootCmd.PersistentFlags().BoolVarP(&cfg.CopyToClipboard, "clipboard", "c", false, "Copy output to clipboard")
	rootCmd.PersistentFlags().BoolVar(&cfg.DisableSecurity, "no-security", false, "Disable security checks (secrets detection)")

	rootCmd.PersistentFlags().BoolVarP(&cfg.SkeletonMode, "skeleton", "s", false, "Strip function bodies (skeleton mode)")
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "Configuration profile to use (defined in .llmpack.yaml)")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file (default .llmpack.yaml)")

	rootCmd.PersistentFlags().BoolVar(&cfg.SymbolsOnly, "symbols", false, "List all symbols in AI-friendly format")
	rootCmd.PersistentFlags().StringVar(&cfg.Implementation, "implementation", "", "Extract full implementation of a symbol")
	rootCmd.PersistentFlags().StringVar(&cfg.FindSymbol, "find", "", "Find files containing a specific symbol")
	rootCmd.PersistentFlags().BoolVar(&cfg.Focus, "focus", false, "In find mode, return only the symbol implementation + skeleton")
}

var treeCmd = &cobra.Command{
	Use:   "tree [path]",
	Short: "Show visual directory tree",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Format = "tree"
		cfg.OutputPath = "-"
		cfg.NoTree = false
		if len(args) == 0 {
			args = []string{"."}
		}
		packRun(cmd, args)
	},
}

func main() {
	setupFlags()
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(treeCmd)

	// If no args or first arg is not a command, default to 'pack'
	if len(os.Args) > 1 {
		found := false
		cmdName := os.Args[1]

		// Check registered commands
		for _, c := range rootCmd.Commands() {
			if c.Name() == cmdName || c.HasAlias(cmdName) {
				found = true
				break
			}
		}

		// Also check built-in Cobra commands and aliases
		if !found {
			builtIns := []string{"completion", "help", "__complete"}
			for _, b := range builtIns {
				if cmdName == b {
					found = true
					break
				}
			}
		}

		// If it's a flag, it's also for the root/pack
		if !found && !strings.HasPrefix(cmdName, "-") {
			os.Args = append([]string{os.Args[0], "pack"}, os.Args[1:]...)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
