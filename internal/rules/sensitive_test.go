package rules

import "testing"

func TestCompileSensitive(t *testing.T) {
	patterns := []string{`(?i)token`, `password\s*=`}

	compiled, err := CompileSensitive(patterns)
	if err != nil {
		t.Fatalf("CompileSensitive returned error: %v", err)
	}

	if len(compiled) != len(patterns) {
		t.Fatalf("CompileSensitive returned %d patterns; want %d", len(compiled), len(patterns))
	}

	if compiled[0].ID != "S1" || compiled[1].ID != "S2" {
		t.Fatalf("unexpected IDs: got %q and %q", compiled[0].ID, compiled[1].ID)
	}

	if !compiled[0].Re.MatchString("TOKEN") {
		t.Fatalf("first regexp does not match expected string")
	}
	if !compiled[1].Re.MatchString("password = value") {
		t.Fatalf("second regexp does not match expected string")
	}
}

func TestCompileSensitive_InvalidPattern(t *testing.T) {
	_, err := CompileSensitive([]string{"("})
	if err == nil {
		t.Fatalf("expected error for invalid regexp, got nil")
	}
}
