package main

import (
	"testing"

	"github.com/dehimik/llmpack/internal/core"
)

func TestFlags(t *testing.T) {
	// Reset cfg
	cfg = core.Config{}

	// Mock appRun to avoid IO
	oldAppRun := appRun
	appRun = func(c core.Config) error { return nil }
	defer func() { appRun = oldAppRun }()

	setupFlags()
	// rootCmd is a global variable in main.go.
	rootCmd.SetArgs([]string{"--symbols", "--find", "mySymbol", "--focus", "test"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() failed: %v", err)
	}

	if !cfg.SymbolsOnly {
		t.Errorf("expected cfg.SymbolsOnly to be true")
	}
	if cfg.FindSymbol != "mySymbol" {
		t.Errorf("expected cfg.FindSymbol to be 'mySymbol', got '%s'", cfg.FindSymbol)
	}
	if !cfg.Focus {
		t.Errorf("expected cfg.Focus to be true")
	}

	// Test implementation flag
	cfg = core.Config{}
	rootCmd.SetArgs([]string{"--implementation", "myFunc", "test"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("rootCmd.Execute() failed for --implementation: %v", err)
	}
	if cfg.Implementation != "myFunc" {
		t.Errorf("expected cfg.Implementation to be 'myFunc', got '%s'", cfg.Implementation)
	}

	// Test validation: --symbols and --implementation cannot both be set
	cfg = core.Config{}
	rootCmd.SetArgs([]string{"--symbols", "--implementation", "myFunc", "test"})
	err = rootCmd.Execute()
	if err == nil {
		t.Errorf("expected error when both --symbols and --implementation are set")
	}
}
