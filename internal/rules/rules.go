package rules

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/iconfire7/loglintergo/internal/config"
)

type RuleID string

const (
	RLowercaseStart RuleID = "LOG001"
	REnglishOnly    RuleID = "LOG002"
	RNoEmojiSpecial RuleID = "LOG003"
	RSensitive      RuleID = "LOG004"
)

type Violation struct {
	ID      RuleID
	Message string
}

// CheckAll проверяет все правила
func CheckAll(msg string, rulesConfig config.Rules, sensitive []*regexp.Regexp) []Violation {
	var out []Violation
	if rulesConfig.Lowercase {
		if v, ok := LowercaseStart(msg); ok {
			out = append(out, v)
		}
	}
	if rulesConfig.English {
		if v, ok := EnglishOnly(msg); ok {
			out = append(out, v)
		}
	}
	if rulesConfig.EmojiOrSpesial {
		if v, ok := NoEmojiOrSpecials(msg); ok {
			out = append(out, v)
		}
	}
	if rulesConfig.Sensitive {
		if v, ok := NoSensitivePatterns(msg, sensitive); ok {
			out = append(out, v)
		}
	}
	return out
}

// LowercaseStart проверяет на строчную букву в начале строки
func LowercaseStart(msg string) (Violation, bool) {
	s := strings.TrimLeft(msg, " \t\r\n")
	if s == "" {
		return Violation{}, false
	}

	r, _ := getFirstRune(s)
	if unicode.IsLetter(r) && unicode.IsUpper(r) {
		return Violation{ID: RLowercaseStart, Message: "log message must not start with an uppercase letter"}, true
	}
	return Violation{}, false
}

// EnglishOnly проверяет на английский язык
func EnglishOnly(msg string) (Violation, bool) {
	for _, r := range msg {
		if unicode.IsLetter(r) && !unicode.In(r, unicode.Latin) {
			return Violation{ID: REnglishOnly, Message: "log message must be in English (Latin letters only)"}, true
		}
	}
	return Violation{}, false
}

// NoEmojiOrSpecials проверяет на спецсимволы и эмодзи
func NoEmojiOrSpecials(msg string) (Violation, bool) {
	for _, r := range msg {
		if !IsAllowedLogChar(r) {
			return Violation{ID: RNoEmojiSpecial, Message: fmt.Sprintf("log message must not contain emoji or special characters (bad rune %U, %q)", r, string(r))}, true
		}
	}
	return Violation{}, false
}

// NoSensitivePatterns проверяет на чувствительные данные
func NoSensitivePatterns(msg string, patterns []*regexp.Regexp) (Violation, bool) {
	if len(patterns) == 0 {
		return Violation{}, false
	}
	for _, re := range patterns {
		if re.MatchString(msg) {
			return Violation{ID: RSensitive, Message: "log message matches sensitive pattern"}, true
		}
	}
	return Violation{}, false
}

// getFirstRune Маленький хелпер чтобы не тащить utf8 в каждый файл
func getFirstRune(s string) (rune, int) {
	for i, r := range s {
		return r, i
	}
	return 0, 0
}

func IsAllowedLogChar(r rune) bool {
	if r > 127 {
		return false
	}

	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}

	switch r {
	case ' ', '\t', '\n', '\r', ':', '.', '%', '[', ']', ',', '/', '_', '-', '=', '+':
		return true
	default:
		return false
	}
}
