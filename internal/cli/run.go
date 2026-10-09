package cli

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	dogabot "github.com/dogabot/dogabot-sdk-go"
)

func flagName(param string) string {
	return strings.ReplaceAll(param, "_", "-")
}

func newIdempotencyKey() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func runOperation(cmd *cobra.Command, op OperationMeta) error {
	colorMode, err := parseColorMode(mustGetString(cmd, "color"))
	if err != nil {
		return err
	}
	cfg, err := resolveConfig(
		mustGetString(cmd, "api-key"),
		mustGetString(cmd, "base-url"),
		mustGetBool(cmd, "json"),
		mustGetBool(cmd, "yes"),
	)
	if err != nil {
		return err
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("missing API key: set DOGABOT_API_KEY or --api-key (or ~/.config/dogabot/config.json)")
	}
	if op.Write && !cfg.Yes {
		return fmt.Errorf("write operation %s requires --yes or DOGABOT_YES=1", op.ID)
	}

	pathParams := map[string]string{}
	for _, p := range op.PathParams {
		v, _ := cmd.Flags().GetString(flagName(p))
		v = strings.TrimSpace(v)
		if v == "" {
			return fmt.Errorf("missing required path flag --%s", flagName(p))
		}
		pathParams[p] = v
	}

	query := url.Values{}
	for _, q := range op.QueryParams {
		v, _ := cmd.Flags().GetString(flagName(q.Name))
		v = strings.TrimSpace(v)
		if v == "" {
			if q.Required {
				return fmt.Errorf("missing required query flag --%s", flagName(q.Name))
			}
			continue
		}
		query.Set(q.Name, v)
	}

	bodyFlag, _ := cmd.Flags().GetString("body")
	body, err := parseBody(bodyFlag, op.HasBody, cmd.InOrStdin())
	if err != nil {
		return err
	}

	var opts []dogabot.RequestOption
	if op.Write {
		key, _ := cmd.Flags().GetString("idempotency-key")
		key = strings.TrimSpace(key)
		if key == "" {
			key = newIdempotencyKey()
		}
		opts = append(opts, dogabot.WithIdempotencyKey(key))
	}

	clientOpts := []dogabot.ClientOption{}
	if cfg.BaseURL != "" {
		clientOpts = append(clientOpts, dogabot.WithBaseURL(cfg.BaseURL))
	}
	client := dogabot.New(cfg.APIKey, clientOpts...)

	ctx, cancel := context.WithTimeout(cmd.Context(), 120*time.Second)
	defer cancel()

	result, err := Call(client, ctx, op.ID, pathParams, query, body, opts...)
	if err != nil {
		return err
	}
	return printResult(cmd.OutOrStdout(), cfg.JSON, colorMode, op, result)
}

func mustGetString(cmd *cobra.Command, name string) string {
	if f := cmd.Flags().Lookup(name); f != nil {
		return f.Value.String()
	}
	if f := cmd.InheritedFlags().Lookup(name); f != nil {
		return f.Value.String()
	}
	for c := cmd; c != nil; c = c.Parent() {
		if f := c.PersistentFlags().Lookup(name); f != nil {
			return f.Value.String()
		}
	}
	return ""
}

func mustGetBool(cmd *cobra.Command, name string) bool {
	s := mustGetString(cmd, name)
	return s == "true"
}

func addOperationFlags(fs *pflag.FlagSet, op OperationMeta) {
	for _, p := range op.PathParams {
		fs.String(flagName(p), "", fmt.Sprintf("path parameter %s (required)", p))
	}
	for _, q := range op.QueryParams {
		req := ""
		if q.Required {
			req = " (required)"
		}
		fs.String(flagName(q.Name), "", fmt.Sprintf("query parameter %s%s", q.Name, req))
	}
	if op.HasBody {
		fs.String("body", "", "JSON body, @file.json, or - for stdin")
	}
	if op.Write {
		fs.String("idempotency-key", "", "Idempotency-Key (auto-generated if omitted)")
	}
}

func lookupOperation(id string) (OperationMeta, bool) {
	for _, op := range Operations {
		if op.ID == id || op.Kebab == id {
			return op, true
		}
	}
	return OperationMeta{}, false
}
