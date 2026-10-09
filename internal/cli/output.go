package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

type colorStyles struct {
	title  lipgloss.Style
	meta   lipgloss.Style
	err    lipgloss.Style
	method lipgloss.Style
	write  lipgloss.Style
	muted  lipgloss.Style
}

func stylesFor(w io.Writer, color bool) colorStyles {
	r := lipgloss.NewRenderer(w)
	if color {
		// Force ANSI even when w is not a TTY (--color=always, tests).
		r.SetColorProfile(termenv.ANSI)
	} else {
		r.SetColorProfile(termenv.Ascii)
	}
	return colorStyles{
		title:  r.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		meta:   r.NewStyle().Foreground(lipgloss.Color("245")),
		err:    r.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
		method: r.NewStyle().Foreground(lipgloss.Color("39")).Bold(true),
		write:  r.NewStyle().Foreground(lipgloss.Color("214")),
		muted:  r.NewStyle().Foreground(lipgloss.Color("245")),
	}
}

func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printResult(w io.Writer, jsonMode bool, colorMode string, op OperationMeta, v any) error {
	if jsonMode {
		return writeJSON(w, v)
	}
	if useColor(colorMode, w, false) {
		s := stylesFor(w, true)
		fmt.Fprintln(w, s.title.Render(op.Method+" "+op.Path))
		if len(op.Scopes) > 0 {
			fmt.Fprintln(w, s.meta.Render("scopes: "+joinComma(op.Scopes)))
		}
		fmt.Fprintln(w)
	}
	return writeJSON(w, v)
}

func printError(cmd *cobra.Command, err error) {
	out := io.Writer(os.Stderr)
	colorMode := colorAuto
	jsonMode := false
	if cmd != nil {
		out = cmd.ErrOrStderr()
		colorMode = mustGetString(cmd, "color")
		jsonMode = mustGetBool(cmd, "json")
	}
	prefix := "error:"
	if useColor(colorMode, out, jsonMode) {
		prefix = stylesFor(out, true).err.Render("error:")
	}
	fmt.Fprintln(out, prefix, err.Error())
}

func printOpsLine(w io.Writer, colorMode string, op OperationMeta) {
	id := fmt.Sprintf("%-40s", op.ID)
	method := fmt.Sprintf("%-8s", op.Method)
	path := fmt.Sprintf("%s/%s", op.TagSlug, op.Kebab)
	write := ""
	if op.Write {
		write = " write"
	}
	if !useColor(colorMode, w, false) {
		fmt.Fprintf(w, "%s %s %s%s\n", id, method, path, write)
		return
	}
	s := stylesFor(w, true)
	fmt.Fprintf(w, "%s %s %s%s\n",
		id,
		s.method.Render(method),
		s.muted.Render(path),
		s.write.Render(write),
	)
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	st, err := f.Stat()
	if err != nil {
		return false
	}
	return (st.Mode() & os.ModeCharDevice) != 0
}

func joinComma(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
