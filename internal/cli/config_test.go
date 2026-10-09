package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigEnv(t *testing.T) {
	t.Setenv("DOGABOT_API_KEY", "from-env")
	t.Setenv("DOGABOT_BASE_URL", "")
	t.Setenv("DOGABOT_YES", "1")
	t.Setenv("DOGABOT_CONFIG", filepath.Join(t.TempDir(), "missing.json"))
	cfg, err := resolveConfig("", "", true, false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "from-env" || !cfg.JSON || !cfg.Yes {
		t.Fatalf("%#v", cfg)
	}
}

func TestResolveConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(`{"api_key":"from-file","base_url":"https://example.test"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DOGABOT_API_KEY", "")
	t.Setenv("DOGABOT_BASE_URL", "")
	t.Setenv("DOGABOT_YES", "")
	t.Setenv("DOGABOT_CONFIG", path)
	cfg, err := resolveConfig("", "", false, false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIKey != "from-file" || cfg.BaseURL != "https://example.test" {
		t.Fatalf("%#v", cfg)
	}
}
