package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the optional .codeindex.yml configuration.
type Config struct {
	SkipPatterns []string `yaml:"skip_patterns"`
}

func defaults() *Config {
	return &Config{
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
	return cfg, nil
}
