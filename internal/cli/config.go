package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds CLI runtime settings loaded from flags, env, and optional JSON.
type Config struct {
	APIKey  string
	BaseURL string
	JSON    bool
	Yes     bool
}

type fileConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func defaultConfigPath() string {
	if p := os.Getenv("DOGABOT_CONFIG"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// Prefer config.json; also accept legacy .toml path name pointing at JSON.
	return filepath.Join(home, ".config", "dogabot", "config.json")
}

func loadFileConfig(path string) (fileConfig, error) {
	var fc fileConfig
	if path == "" {
		return fc, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fc, nil
		}
		return fc, err
	}
	if err := json.Unmarshal(b, &fc); err != nil {
		return fc, fmt.Errorf("parse config %s: %w", path, err)
	}
	return fc, nil
}

func resolveConfig(apiKeyFlag, baseURLFlag string, jsonOut, yes bool) (Config, error) {
	fc, err := loadFileConfig(defaultConfigPath())
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		APIKey:  firstNonEmpty(apiKeyFlag, os.Getenv("DOGABOT_API_KEY"), fc.APIKey),
		BaseURL: firstNonEmpty(baseURLFlag, os.Getenv("DOGABOT_BASE_URL"), fc.BaseURL),
		JSON:    jsonOut,
		Yes:     yes || envTruthy("DOGABOT_YES"),
	}
	return cfg, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func envTruthy(name string) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	return v == "1" || v == "true" || v == "yes"
}
