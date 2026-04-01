package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the optional .codeindex.yml configuration.
type Config struct {
	SkipPatterns []string `yaml:"skip_patterns"`
	Gitignore    string   `yaml:"gitignore"` // "local" (default) or "global"
}

func defaults() *Config {
	return &Config{
		Gitignore: "local",
		SkipPatterns: []string{
			"**/*.test.ts",
			"**/*.spec.ts",
			"**/__mocks__/**",
			"*.generated.ts",
			"**/node_modules/**",
			"**/bin/**",
			"**/obj/**",
			"**/dist/**",
			"**/.git/**",
		},
	}
}

// LoadFile reads only the user-defined values from .codeindex.yml (no defaults merged).
// Returns an empty Config if the file does not exist.
func LoadFile(dir string) (*Config, error) {
	path := filepath.Join(dir, ".codeindex.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return &Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

// SaveFile writes only the user-defined config to .codeindex.yml (no defaults).
func SaveFile(dir string, cfg *Config) error {
	path := filepath.Join(dir, ".codeindex.yml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads .codeindex.yml from dir (or returns defaults if not found).
func Load(dir string) (*Config, error) {
	cfg := defaults()

	path := filepath.Join(dir, ".codeindex.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, nil
	}

	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return cfg, err
	}

	// Merge: append file patterns to defaults
	cfg.SkipPatterns = append(cfg.SkipPatterns, fileCfg.SkipPatterns...)
	if fileCfg.Gitignore != "" {
		cfg.Gitignore = fileCfg.Gitignore
	}
	return cfg, nil
}
