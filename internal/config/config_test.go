package config

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()

	if !cfg.Rules.Lowercase || !cfg.Rules.English || !cfg.Rules.EmojiOrSpesial || !cfg.Rules.Sensitive {
		t.Fatalf("all default rules must be enabled: %+v", cfg.Rules)
	}

	if len(cfg.SensitivePatterns) == 0 {
		t.Fatalf("default sensitive patterns must not be empty")
	}
}
