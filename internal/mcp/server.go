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
func Start(baseCfg core.Config) error {
	s := mcp.NewServer(stdio.NewStdioServerTransport(), mcp.WithName("llmpack"), mcp.WithVersion("1.0.0"))

	// Helper to create a tool config based on base config
	newToolCfg := func() core.Config {
		c := baseCfg
		c.NoStdin = true
		c.NoTree = true
		return c
	}

	// 1. List Symbols
	err := s.RegisterTool("list_symbols", "Get a compact index of all functions, classes, and methods in a path", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := newToolCfg()
		cfg.InputPaths = []string{args.Path} // Explicitly target requested path
		cfg.SymbolsOnly = true
		cfg.CustomWriter = &buf

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
		cfg := newToolCfg()
		cfg.InputPaths = []string{args.Path} // Explicitly target target file
		cfg.Implementation = args.SymbolName
		cfg.CustomWriter = &buf

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
		cfg := newToolCfg()
		cfg.InputPaths = []string{"."} // Scan project root
		cfg.FindSymbol = args.Query
		cfg.Focus = args.Focus
		cfg.CustomWriter = &buf

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})
	if err != nil {
		return err
	}

	// 4. Security Scan
	err = s.RegisterTool("scan_security", "Check the project for hardcoded secrets and sensitive files", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var logBuf bytes.Buffer
		var discard bytes.Buffer
		cfg := newToolCfg()
		cfg.InputPaths = []string{args.Path}
		cfg.LogWriter = &logBuf
		cfg.CustomWriter = &discard
		cfg.DisableSecurity = false
		cfg.Format = "tree" // Use tree format for fast walk

		_ = app.Run(cfg)

		if logBuf.Len() == 0 {
			return mcp.NewToolResponse(mcp.NewTextContent("No security issues detected.")), nil
		}

		return mcp.NewToolResponse(mcp.NewTextContent(logBuf.String())), nil
	})
	if err != nil {
		return err
	}

	// 5. Estimate Cost
	type EstimateArgs struct {
		Path  string `json:"path" jsonschema:"required,description=Path to the project"`
		Model string `json:"model" jsonschema:"description=Model name (gpt-4o, claude-3-5-sonnet, gemini-1.5-pro)"`
	}
	err = s.RegisterTool("estimate_cost", "Calculate token count and estimated cost for the project", func(args EstimateArgs) (*mcp.ToolResponse, error) {
		var logBuf bytes.Buffer
		var discard bytes.Buffer
		cfg := newToolCfg()
		cfg.InputPaths = []string{args.Path}
		cfg.LogWriter = &logBuf
		cfg.CustomWriter = &discard
		cfg.CountTokens = true
		cfg.ModelName = args.Model
		cfg.Format = "tree" // Fast run

		if cfg.ModelName == "" {
			cfg.ModelName = "gpt-4o"
		}

		_ = app.Run(cfg)

		return mcp.NewToolResponse(mcp.NewTextContent(logBuf.String())), nil
	})
	if err != nil {
		return err
	}

	// 6. Get Tree
	err = s.RegisterTool("get_tree", "Get a visual directory tree of the project", func(args ListSymbolsArgs) (*mcp.ToolResponse, error) {
		var buf bytes.Buffer
		cfg := newToolCfg()
		cfg.InputPaths = []string{args.Path}
		cfg.Format = "tree"
		cfg.CustomWriter = &buf
		cfg.NoTree = false // We WANT the tree here

		if err := app.Run(cfg); err != nil {
			return nil, err
		}

		return mcp.NewToolResponse(mcp.NewTextContent(buf.String())), nil
	})
	if err != nil {
		return err
	}

	if err := s.Serve(); err != nil {
		return err
	}

	// Serve is asynchronous in metoro-io/mcp-golang, so we wait
	select {}
}

