package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBodyLiteral(t *testing.T) {
	v, err := parseBody(`{"a":1}`, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("want map, got %T", v)
	}
	n, ok := m["a"].(json.Number)
	if !ok || n.String() != "1" {
		t.Fatalf("want number 1, got %#v", m["a"])
	}
}

func TestParseBodyStdin(t *testing.T) {
	v, err := parseBody("-", true, strings.NewReader(`{"ok":true}`))
	if err != nil {
		t.Fatal(err)
	}
	m := v.(map[string]any)
	if m["ok"] != true {
		t.Fatalf("got %#v", v)
	}
}

func TestParseBodyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "body.json")
	if err := os.WriteFile(path, []byte(`{"x":"y"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	v, err := parseBody("@"+path, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if v.(map[string]any)["x"] != "y" {
		t.Fatalf("got %#v", v)
	}
}

func TestParseBodyRejectedWhenNoBody(t *testing.T) {
	_, err := parseBody(`{}`, false, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
