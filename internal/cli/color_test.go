package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestParseColorMode(t *testing.T) {
	m, err := parseColorMode("")
	if err != nil || m != colorAuto {
		t.Fatalf("empty: %q %v", m, err)
	}
	m, err = parseColorMode("ALWAYS")
	if err != nil || m != colorAlways {
		t.Fatalf("always: %q %v", m, err)
	}
	if _, err := parseColorMode("rainbow"); err == nil {
		t.Fatal("expected error")
	}
}

func TestUseColorJSONDisables(t *testing.T) {
	t.Setenv("NO_COLOR", "") // empty → not disabling (no-color.org)
	if useColor(colorAlways, os.Stdout, true) {
		t.Fatal("--json should disable color")
	}
}

func TestUseColorNOCOLORWins(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if useColor(colorAlways, os.Stdout, false) {
		t.Fatal("NO_COLOR should win over always")
	}
}

func TestUseColorNever(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	if useColor(colorNever, os.Stdout, false) {
		t.Fatal("never")
	}
}

func TestUseColorAlwaysNonTTY(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	buf := &bytes.Buffer{}
	if !useColor(colorAlways, buf, false) {
		t.Fatal("always should color non-TTY writers")
	}
	if useColor(colorAuto, buf, false) {
		t.Fatal("auto should not color buffer")
	}
}

func TestOpsColorAlwaysContainsANSI(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"ops", "--tag", "account", "--color=always"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "getMe") {
		t.Fatalf("missing getMe: %q", out)
	}
	if !strings.Contains(out, "\x1b[") {
		n := 120
		if len(out) < n {
			n = len(out)
		}
		t.Fatalf("expected ANSI escapes with --color=always, got %q", out[:n])
	}
}

func TestOpsColorNeverPlain(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	root := NewRootForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"ops", "--tag", "account", "--color=never"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "\x1b[") {
		t.Fatal("unexpected ANSI with --color=never")
	}
}
