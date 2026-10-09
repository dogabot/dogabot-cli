package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Execute runs the dogabot CLI.
func Execute() error {
	root := newRootCmd()
	err := root.Execute()
	if err != nil {
		printError(root, err)
	}
	return err
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "dogabot",
		Short:         "Official dogabot CLI for the public REST API",
		Long:          "Call every allowlisted dogabot REST operation (same surface as the official SDKs).\nMCP remains a separate agent protocol — see https://docs.dogabot.com/mcp/",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.PersistentFlags().String("api-key", "", "API key (prefer DOGABOT_API_KEY)")
	root.PersistentFlags().String("base-url", "", "API base URL (default https://api.dogabot.com)")
	root.PersistentFlags().Bool("json", false, "machine-readable JSON only (no decorative color)")
	root.PersistentFlags().Bool("yes", false, "confirm write operations (or set DOGABOT_YES=1)")
	root.PersistentFlags().String("color", colorAuto, "color output: auto|always|never (respects NO_COLOR; --json disables)")

	root.AddCommand(newVersionCmd())
	root.AddCommand(newOpsCmd())
	root.AddCommand(newCallCmd())
	registerGroupedCommands(root)
	registerFriendlyAliases(root)
	root.AddCommand(newCompletionCmd(root))
	installColoredHelp(root)

	return root
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "dogabot "+Version)
		},
	}
}

func newOpsCmd() *cobra.Command {
	var tag string
	cmd := &cobra.Command{
		Use:   "ops",
		Short: "List allowlisted operations (SDK / REST parity)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := parseColorMode(mustGetString(cmd, "color")); err != nil {
				return err
			}
			tag = strings.TrimSpace(strings.ToLower(tag))
			colorMode := mustGetString(cmd, "color")
			w := cmd.OutOrStdout()
			for _, op := range Operations {
				if tag != "" && op.TagSlug != tag && strings.ToLower(op.Tag) != tag {
					continue
				}
				printOpsLine(w, colorMode, op)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&tag, "tag", "", "filter by OpenAPI tag slug (e.g. account, terminal)")
	return cmd
}

func newCallCmd() *cobra.Command {
	return &cobra.Command{
		Use:                "call <operationId>",
		Short:              "Invoke an operation by exact SDK operationId (or kebab-case)",
		DisableFlagParsing: true,
		SilenceErrors:      true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("usage: dogabot call <operationId> [flags]")
			}
			op, ok := lookupOperation(args[0])
			if !ok {
				return fmt.Errorf("unknown operation %q (try: dogabot ops)", args[0])
			}
			inner := &cobra.Command{Use: "call", SilenceErrors: true, SilenceUsage: true}
			inner.Flags().AddFlagSet(cmd.Root().PersistentFlags())
			addOperationFlags(inner.Flags(), op)
			inner.SetArgs(args[1:])
			inner.SetContext(cmd.Context())
			inner.SetOut(cmd.OutOrStdout())
			inner.SetErr(cmd.ErrOrStderr())
			inner.SetIn(cmd.InOrStdin())
			if err := inner.ParseFlags(args[1:]); err != nil {
				return err
			}
			return runOperation(inner, op)
		},
	}
}

func registerGroupedCommands(root *cobra.Command) {
	byTag := map[string]*cobra.Command{}

	for _, op := range Operations {
		op := op
		parent, ok := byTag[op.TagSlug]
		if !ok {
			parent = &cobra.Command{
				Use:   op.TagSlug,
				Short: fmt.Sprintf("%s operations", op.Tag),
			}
			byTag[op.TagSlug] = parent
			root.AddCommand(parent)
		}
		c := &cobra.Command{
			Use:     op.Kebab,
			Aliases: []string{op.ID},
			Short:   fmt.Sprintf("%s %s", op.Method, op.Path),
			RunE: func(cmd *cobra.Command, _ []string) error {
				return runOperation(cmd, op)
			},
		}
		addOperationFlags(c.Flags(), op)
		parent.AddCommand(c)
	}
}

func registerFriendlyAliases(root *cobra.Command) {
	names := make([]string, 0, len(friendlyAliases))
	for name := range friendlyAliases {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		opID := friendlyAliases[name]
		found, ok := lookupOperation(opID)
		if !ok {
			continue
		}
		op := found
		c := &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Alias for %s", op.ID),
			RunE: func(cmd *cobra.Command, _ []string) error {
				return runOperation(cmd, op)
			},
		}
		addOperationFlags(c.Flags(), op)
		root.AddCommand(c)
	}
}

func newCompletionCmd(root *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion script",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return root.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell %q", args[0])
			}
		},
	}
}

// NewRootForTest exposes the root command for unit tests.
func NewRootForTest() *cobra.Command {
	return newRootCmd()
}
