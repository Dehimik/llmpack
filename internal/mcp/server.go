package mcp

import (
	"bytes"

	"github.com/dehimik/llmpack/internal/app"
	"github.com/dehimik/llmpack/internal/core"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type ListSymbolsArgs struct {
	Path string `json:"path" jsonschema:"required,description=Path to the directory or file to list symbols from"`
}

type GetCodeArgs struct {
	Path       string `json:"path" jsonschema:"required,description=Path to the file"`
	SymbolName string `json:"symbol_name" jsonschema:"required,description=Name of the symbol (function, class, method) to extract"`
}

type SearchArgs struct {
	Query string `json:"query" jsonschema:"required,description=Symbol name to search for across the project"`
	Focus bool   `json:"focus" jsonschema:"description=If true, return only the symbol implementation and skeletons for the rest of the file. Default is true.,default=true"`
}
func Start() error {
	s := mcp.NewServer(stdio.NewStdioServerTransport())

	// 1. List Symbols
	err := s.RegisterTool("list_symbols", "Get a compact index of all functions, classes, and methods in a path", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := core.Config{
			InputPaths:   []string{args.Path},
			SymbolsOnly:  true,
			CustomWriter: &buf,
			NoTree:       true,
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})
	if err != nil {
		return err
	}

	// 2. Get Code (Implementation)
	err = s.RegisterTool("get_code", "Extract full implementation of a specific symbol while skeletonizing the rest of the file", func(args GetCodeArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := core.Config{
			InputPaths:     []string{args.Path},
			Implementation: args.SymbolName,
			CustomWriter:   &buf,
			NoTree:         true,
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})
	if err != nil {
		return err
	}

	// 3. Search
	err = s.RegisterTool("search", "Find all files containing a specific symbol and return their content (optionally focused)", func(args SearchArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := core.Config{
			InputPaths:   []string{"."},
			FindSymbol:   args.Query,
			Focus:        args.Focus,
			CustomWriter: &buf,
			NoTree:       true,
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})
	if err != nil {
		return err
	}

	// 4. Security Scan
	s.RegisterTool("scan_security", "Check the project for hardcoded secrets and sensitive files", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var logBuf bytes.Buffer
		cfg := core.Config{
			InputPaths:      []string{args.Path},
			LogWriter:       &logBuf,
			DisableSecurity: false,
			Format:          "tree", // We just want to trigger the walk
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		_ = app.Run(cfg) // Run can fail on errors, but we care about log content

		if logBuf.Len() == 0 {
			return mcp.NewToolResponse(mcp.NewTextContent("No security issues detected.")), nil
		}

		return mcp.NewToolResponse(mcp.NewTextContent(logBuf.String())), nil
	})

	// 5. Estimate Cost
	type EstimateArgs struct {
		Path  string `json:"path" jsonschema:"required,description=Path to the project"`
		Model string `json:"model" jsonschema:"description=Model name (gpt-4o, claude-3-5-sonnet, gemini-1.5-pro)"`
	}
	s.RegisterTool("estimate_cost", "Calculate token count and estimated cost for the project", func(args EstimateArgs) (*mcp.ToolResponse, error) {
		var logBuf bytes.Buffer
		cfg := core.Config{
			InputPaths:  []string{args.Path},
			LogWriter:   &logBuf,
			CountTokens: true,
			ModelName:   args.Model,
			Format:      "tree", // Fast run
		}
		if cfg.ModelName == "" {
			cfg.ModelName = "gpt-4o"
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		_ = app.Run(cfg)

		return mcp.NewToolResponse(mcp.NewTextContent(logBuf.String())), nil
	})

	// 6. Get Tree
	s.RegisterTool("get_tree", "Get a visual directory tree of the project", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := core.Config{
			InputPaths:   []string{args.Path},
			Format:       "tree",
			CustomWriter: &buf,
		}

		cfg.IgnorePatterns = []string{".git", "node_modules", "vendor", "dist", "build"}

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})

	return s.Serve()
}

