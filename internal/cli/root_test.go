package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"version"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "dogabot ") {
		t.Fatalf("got %q", buf.String())
	}
}

func TestOpsListsGetMe(t *testing.T) {
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"ops", "--tag", "account"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "getMe") {
		t.Fatalf("got %q", buf.String())
	}
}

func TestWriteRequiresYes(t *testing.T) {
	root := NewRootForTest()
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"automations", "post-followers",
		"--api-key", "dbk_test",
		"--body", `{}`,
	})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("want --yes error, got %v", err)
	}
}

func TestGetMeAgainstFakeServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/me" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer dbk_test" {
			t.Fatalf("auth %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "u1"})
	}))
	defer srv.Close()

	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{
		"account", "get-me",
		"--api-key", "dbk_test",
		"--base-url", srv.URL,
		"--json",
	})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"id"`) {
		t.Fatalf("got %q", buf.String())
	}
}

func TestMissingPathParam(t *testing.T) {
	root := NewRootForTest()
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	root.SetArgs([]string{
		"automations", "get-followers-by-id",
		"--api-key", "dbk_test",
	})
	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "--id") {
		t.Fatalf("want --id error, got %v", err)
	}
}

func TestWhoamiAlias(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"whoami", "--api-key", "k", "--base-url", srv.URL, "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "ok") {
		t.Fatalf("got %q", buf.String())
	}
}
