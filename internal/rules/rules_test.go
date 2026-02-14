package rules

import (
	"regexp"
	"testing"

	"github.com/iconfire7/loglintergo/internal/config"
)

func TestLowercaseStart(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"user created", false},
		{"User created", true},
		{" User created", true},
		{"123 ok", false},
		{"", false},
		{"add testLower", false},
	}

	for _, tc := range cases {
		_, got := LowercaseStart(tc.in)
		if got != tc.want {
			t.Errorf("LowercaseStart(%q) = %v; want %v", tc.in, got, tc.want)
		}
	}
}

func TestEnglishOnly(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"User created", false},
		{"Создан пользователь", true},
		{"User создан", true},
	}

	for _, tc := range cases {
		_, got := EnglishOnly(tc.in)
		if got != tc.want {
			t.Errorf("EnglishOnly(%q) = %v; want %v", tc.in, got, tc.want)
		}
	}
}

func TestNoEmojiOrSpecials(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"OK", false},
		{"token=", false},
		{"Hello🙂", true},
		{"Привет", true},
	}
	for _, tc := range cases {
		_, got := NoEmojiOrSpecials(tc.in)
		if got != tc.want {
			t.Fatalf("NoEmojiOrSpecials(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestNoSensitiveKeywords(t *testing.T) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\btoken\b`),
		regexp.MustCompile(`(?i)\bpassword\b`),
	}

	cases := []struct {
		in   string
		want bool
	}{
		{"User logged in", false},
		{"User token expired", true},
		{"Password is invalid", true},
	}
	for _, tc := range cases {
		_, got := NoSensitivePatterns(tc.in, patterns)
		if got != tc.want {
			t.Fatalf("NoSensitivePatterns(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}

func TestCheckAll(t *testing.T) {
	rulesConfig := config.Rules{
		Lowercase:      true,
		English:        true,
		EmojiOrSpesial: true,
		Sensitive:      true,
	}

	sensitive := []*regexp.Regexp{regexp.MustCompile(`(?i)token`)}
	violations := CheckAll("User token🙂", rulesConfig, sensitive)
	if len(violations) != 3 {
		t.Fatalf("CheckAll returned %d violations; want 3", len(violations))
	}

	wantIDs := map[RuleID]bool{
		RLowercaseStart: true,
		RSensitive:      true,
		RNoEmojiSpecial: true,
	}
	for _, v := range violations {
		if !wantIDs[v.ID] {
			t.Fatalf("unexpected violation ID: %s", v.ID)
		}
	}
}

func TestGetFirstRune(t *testing.T) {
	r, idx := getFirstRune("Привет")
	if r != 'П' || idx != 0 {
		t.Fatalf("getFirstRune returned (%q, %d); want ('П', 0)", r, idx)
	}

	r, idx = getFirstRune("")
	if r != 0 || idx != 0 {
		t.Fatalf("getFirstRune for empty string returned (%q, %d); want (0, 0)", r, idx)
	}
}

func TestIsAllowedLogChar(t *testing.T) {
	cases := []struct {
		name string
		r    rune
		want bool
	}{
		{name: "latin letter", r: 'a', want: true},
		{name: "digit", r: '7', want: true},
		{name: "allowed punctuation", r: ':', want: true},
		{name: "emoji", r: '🙂', want: false},
		{name: "cyrillic", r: 'П', want: false},
		{name: "not allowed punctuation", r: '!', want: false},
	}

	for _, tc := range cases {
		got := IsAllowedLogChar(tc.r)
		if got != tc.want {
			t.Fatalf("IsAllowedLogChar(%q) = %v; want %v", tc.r, got, tc.want)
		}
	}
}
