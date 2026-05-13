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

	return s.Serve()
}
