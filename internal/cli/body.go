package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// parseBody reads JSON from a literal, @file, or "-" (stdin).
// Empty bodyFlag with hasBody returns nil (SDK may send empty body).
func parseBody(bodyFlag string, hasBody bool, stdin io.Reader) (any, error) {
	if !hasBody {
		if strings.TrimSpace(bodyFlag) != "" {
			return nil, fmt.Errorf("--body is not used by this operation")
		}
		return nil, nil
	}
	raw := strings.TrimSpace(bodyFlag)
	if raw == "" {
		return nil, nil
	}
	var r io.Reader
	switch {
	case raw == "-":
		r = stdin
	case strings.HasPrefix(raw, "@"):
		path := strings.TrimPrefix(raw, "@")
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	default:
		r = strings.NewReader(raw)
	}
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("decode body JSON: %w", err)
	}
	return v, nil
}
