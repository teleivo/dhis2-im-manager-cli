package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConcat(t *testing.T) {
	if diff := cmp.Diff("ab", concat("a", "b")); diff != "" {
		t.Fatalf("concat(%s, %s) mismatch (-want +got): %s\n", "a", "b", diff)
	}
}
