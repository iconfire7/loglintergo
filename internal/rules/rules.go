package rules

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

func CheckAll(msg string) []Violation {
	var out []Violation
	if v, ok := LowercaseStart(msg); ok {
		out = append(out, v)
	}
	if v, ok := EnglishOnly(msg); ok {
		out = append(out, v)
	}
	if v, ok := NoEmojiOrSpecials(msg); ok {
		out = append(out, v)
	}
	if v, ok := NoSensitiveKeywords(msg); ok {
		out = append(out, v)
	}
	return out
}

func LowercaseStart(msg string) (Violation, bool) {
	// TODO
	return Violation{}, false
}

func EnglishOnly(msg string) (Violation, bool) {
	// TODO
	return Violation{}, false
}

func NoEmojiOrSpecials(msg string) (Violation, bool) {
	// TODO
	return Violation{}, false
}

func NoSensitiveKeywords(msg string) (Violation, bool) {
	// TODO
	return Violation{}, false
}
