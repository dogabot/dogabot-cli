package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func installColoredHelp(root *cobra.Command) {
	root.SetHelpFunc(coloredHelpFunc)
	root.SetUsageFunc(coloredUsageFunc)
}

func coloredHelpFunc(cmd *cobra.Command, _ []string) {
	printColoredHelp(cmd)
}

func coloredUsageFunc(cmd *cobra.Command) error {
	printColoredUsage(cmd)
	return nil
}

func printColoredHelp(cmd *cobra.Command) {
	w := cmd.OutOrStdout()
	s, on := helpStyles(cmd, w)

	if cmd.Long != "" {
		fmt.Fprintln(w, styleOrPlain(on, s.title, cmd.Long))
	} else if cmd.Short != "" {
		fmt.Fprintln(w, styleOrPlain(on, s.title, cmd.Short))
	}
	fmt.Fprintln(w)

	printColoredUsage(cmd)

	if cmds := helpCommands(cmd); len(cmds) > 0 {
		fmt.Fprintln(w, styleOrPlain(on, s.title, "Available Commands:"))
		nameWidth := 0
		for _, c := range cmds {
			if n := len(c.Name()); n > nameWidth {
				nameWidth = n
			}
		}
		for _, c := range cmds {
			name := fmt.Sprintf("%-*s", nameWidth, c.Name())
			short := strings.TrimSpace(c.Short)
			if on {
				fmt.Fprintf(w, "  %s  %s\n", s.method.Render(name), s.muted.Render(short))
			} else {
				fmt.Fprintf(w, "  %s  %s\n", name, short)
			}
		}
		fmt.Fprintln(w)
	}

	if cmd.HasAvailableLocalFlags() {
		fmt.Fprintln(w, styleOrPlain(on, s.title, "Flags:"))
		printFlagDefaults(w, s, on, cmd.LocalFlags())
		fmt.Fprintln(w)
	}

	if cmd.HasAvailableInheritedFlags() {
		fmt.Fprintln(w, styleOrPlain(on, s.title, "Global Flags:"))
		printFlagDefaults(w, s, on, cmd.InheritedFlags())
		fmt.Fprintln(w)
	}

	if cmd.HasAvailableSubCommands() {
		hint := fmt.Sprintf(`Use "%s [command] --help" for more information about a command.`, cmd.CommandPath())
		fmt.Fprintln(w, styleOrPlain(on, s.muted, hint))
	}
}

func printColoredUsage(cmd *cobra.Command) {
	w := cmd.OutOrStdout()
	s, on := helpStyles(cmd, w)
	fmt.Fprintln(w, styleOrPlain(on, s.title, "Usage:"))
	useLine := cmd.UseLine()
	if on {
		fmt.Fprintf(w, "  %s\n\n", s.method.Render(useLine))
	} else {
		fmt.Fprintf(w, "  %s\n\n", useLine)
	}
}

func helpStyles(cmd *cobra.Command, w io.Writer) (colorStyles, bool) {
	on := useColor(lookupColorMode(cmd), w, lookupJSON(cmd))
	return stylesFor(w, on), on
}

func styleOrPlain(on bool, style interface{ Render(...string) string }, text string) string {
	if !on {
		return text
	}
	return style.Render(text)
}

func helpCommands(cmd *cobra.Command) []*cobra.Command {
	out := make([]*cobra.Command, 0, len(cmd.Commands()))
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		out = append(out, c)
	}
	return out
}

func printFlagDefaults(w io.Writer, s colorStyles, on bool, fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		line := flagHelpLine(f)
		if on {
			// Color the flag names; keep the description muted after two spaces.
			name, desc, ok := strings.Cut(line, "  ")
			if !ok {
				fmt.Fprintln(w, s.muted.Render(line))
				return
			}
			// Trim leading indent from pflag's default format for re-indent.
			name = strings.TrimLeft(name, " ")
			fmt.Fprintf(w, "  %s  %s\n", s.method.Render(name), s.muted.Render(strings.TrimSpace(desc)))
			return
		}
		fmt.Fprintln(w, line)
	})
}

func flagHelpLine(f *pflag.Flag) string {
	// Match pflag's FlagUsLine shape closely enough for humans.
	line := "  "
	if f.Shorthand != "" && f.ShorthandDeprecated == "" {
		line += fmt.Sprintf("-%s, --%s", f.Shorthand, f.Name)
	} else {
		line += fmt.Sprintf("    --%s", f.Name)
	}
	varname, usage := pflag.UnquoteUsage(f)
	if varname != "" {
		line += " " + varname
	}
	if usage != "" {
		line += "  " + usage
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "[]" && f.DefValue != "0" {
		if f.Value.Type() == "string" {
			line += fmt.Sprintf(` (default %q)`, f.DefValue)
		} else {
			line += fmt.Sprintf(" (default %s)", f.DefValue)
		}
	}
	return line
}

func lookupColorMode(cmd *cobra.Command) string {
	if cmd == nil {
		return colorAuto
	}
	f := cmd.Flags().Lookup("color")
	if f == nil {
		f = cmd.InheritedFlags().Lookup("color")
	}
	if f == nil {
		return colorAuto
	}
	m, err := parseColorMode(f.Value.String())
	if err != nil {
		return colorAuto
	}
	return m
}

func lookupJSON(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	f := cmd.Flags().Lookup("json")
	if f == nil {
		f = cmd.InheritedFlags().Lookup("json")
	}
	if f == nil {
		return false
	}
	return f.Value.String() == "true"
}
