package rules

import "testing"

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
	cases := []struct {
		in   string
		want bool
	}{
		{"User logged in", false},
		{"User token expired", true},
		{"Password is invalid", true},
	}
	for _, tc := range cases {
		_, got := NoSensitiveKeywords(tc.in)
		if got != tc.want {
			t.Fatalf("NoSensitiveKeywords(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}
