package loglinter

import (
	"github.com/iconfire7/loglintergo/internal/rules"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strings"
	"unicode"
)

func fixTargetForFirstArg(pass *analysis.Pass, call *ast.CallExpr) (pos, end token.Pos, ok bool) {
	if len(call.Args) == 0 {
		return token.NoPos, token.NoPos, false
	}
	expr := call.Args[0]

	if ce, ok := expr.(*ast.CallExpr); ok && isFmtSprintf(pass, ce) && len(ce.Args) > 0 {
		return ce.Args[0].Pos(), ce.Args[0].End(), true
	}

	return expr.Pos(), expr.End(), true
}

func suggestFixForViolation(id rules.RuleID, msg string) (string, bool) {
	switch id {
	case rules.RLowercaseStart:
		return fixLowercaseStart(msg)
	case rules.RNoEmojiSpecial:
		return fixNoEmojiOrSpecial(msg)
	case rules.RSensitive:
		return "", false
	case rules.REnglishOnly:
		return "", false
	default:
		return "", false
	}
}

func fixLowercaseStart(msg string) (string, bool) {
	rs := []rune(msg)

	i := 0
	for i < len(rs) {
		if rs[i] == ' ' || rs[i] == '\t' || rs[i] == '\n' || rs[i] == '\r' {
			i++
		} else {
			break
		}
	}
	if i >= len(rs) {
		return "", false
	}

	if unicode.IsLetter(rs[i]) && unicode.IsUpper(rs[i]) {
		rs[i] = unicode.ToLower(rs[i])
		return string(rs), true
	}
	return "", false
}

func fixNoEmojiOrSpecial(msg string) (string, bool) {
	var b strings.Builder
	b.Grow(len(msg))
	changed := false
	for _, r := range msg {
		if !rules.IsAllowedLogChar(r) {
			changed = true
			continue
		}
		b.WriteRune(r)
	}
	if !changed {
		return "", false
	}
	return b.String(), true
}
