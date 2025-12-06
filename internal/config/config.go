package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var DefaultIgnores = []string{
	".git", ".idea", ".vscode", ".obsidian",
	"node_modules", "vendor", "dist", "build",
	"*.lock", "*.log", "*.exe", "*.bin",
	".DS_Store",
}

type Settings struct {
	Format       string `yaml:"format"`
	IgnoreGit    bool   `yaml:"ignore_git"`
	SkeletonMode bool   `yaml:"skeleton"`
	NoTree       bool   `yaml:"no_tree"`
	Tokens       bool   `yaml:"tokens"`
	ModelName    string `yaml:"model_name"`
}

type FileConfig struct {
	Global   Settings            `yaml:"global"`
	Profiles map[string]Settings `yaml:"profiles"`
	Ignore   []string            `yaml:"ignore"` // global ignore list
}

func Load(customPath string) (*FileConfig, error) {
	cfg := &FileConfig{
		Global: Settings{
			Format:       "xml",
			IgnoreGit:    true,
			SkeletonMode: false,
			Tokens:       true,
			ModelName:    "gpt-4o", // Default
		},
		Ignore:   DefaultIgnores,
		Profiles: make(map[string]Settings),
	}

	var paths []string

	// Explicit path via flag
	if customPath != "" {
		paths = append(paths, customPath)
	}

	// Current Directory
	paths = append(paths, ".llmpack.yaml")

	// Home Directory
	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".llmpack.yaml"))
	}

	var configPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			configPath = p
			break
		}
	}

	if customPath != "" && configPath != customPath {
		return nil, fmt.Errorf("config file not found at %s", customPath)
	}

	if configPath == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if len(cfg.Ignore) == 0 {
		cfg.Ignore = DefaultIgnores
	}

	return cfg, nil
}
