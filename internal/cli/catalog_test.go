package cli

import (
	"testing"
)

func TestOperationCatalogNonEmpty(t *testing.T) {
	if len(Operations) == 0 {
		t.Fatal("Operations empty")
	}
	if len(OperationIDs) != len(Operations) {
		t.Fatalf("OperationIDs=%d Operations=%d", len(OperationIDs), len(Operations))
	}
	seen := map[string]bool{}
	for _, op := range Operations {
		if op.ID == "" || op.Kebab == "" || op.TagSlug == "" {
			t.Fatalf("incomplete op %#v", op)
		}
		if seen[op.ID] {
			t.Fatalf("duplicate %s", op.ID)
		}
		seen[op.ID] = true
	}
	if !seen["getMe"] {
		t.Fatal("missing getMe")
	}
}

func TestLookupOperation(t *testing.T) {
	op, ok := lookupOperation("getMe")
	if !ok || op.Kebab != "get-me" {
		t.Fatalf("getMe: %#v ok=%v", op, ok)
	}
	op2, ok := lookupOperation("get-me")
	if !ok || op2.ID != "getMe" {
		t.Fatalf("kebab: %#v ok=%v", op2, ok)
	}
}
