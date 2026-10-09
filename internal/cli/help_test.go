package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestMarketsHelpColorAlwaysContainsANSI(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"markets", "--help", "--color=always"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Available Commands:") {
		t.Fatalf("missing commands section: %q", out)
	}
	if !strings.Contains(out, "get-markets") {
		t.Fatalf("missing get-markets: %q", out)
	}
	if !strings.Contains(out, "\x1b[") {
		n := 160
		if len(out) < n {
			n = len(out)
		}
		t.Fatalf("expected ANSI with --color=always, got %q", out[:n])
	}
}

func TestMarketsHelpColorNeverPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"markets", "--help", "--color=never"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Available Commands:") {
		t.Fatalf("missing commands section: %q", out)
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatal("unexpected ANSI with --color=never")
	}
}

func TestRootHelpColorAlwaysContainsANSI(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"--help", "--color=always"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "markets") {
		t.Fatalf("missing markets: %q", out)
	}
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI on root help, got %q", out[:min(160, len(out))])
	}
}
