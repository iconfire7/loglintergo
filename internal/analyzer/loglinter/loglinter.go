package loglinter

import (
	"github.com/iconfire7/loglintergo/internal/config"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"

	"github.com/iconfire7/loglintergo/internal/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var logMethods = map[string]bool{
	"Debug": true,
	"Info":  true,
	"Warn":  true,
	"Error": true,
}

func New(cfg config.Config, sensitive []*regexp.Regexp) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "loglintergo",
		Doc:      "checks log messages for style/safety rules",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run(cfg, sensitive),
	}
}

// run — основная функция анализа пакета.
func run(cfg config.Config, sensitive []*regexp.Regexp) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
			call := n.(*ast.CallExpr)

			kind := detectLoggerCall(pass, call)
			if kind == "" {
				return
			}

			msg, pos, ok := extractFirstStringArg(pass, call)
			if !ok {
				return
			}

			violations := rules.CheckAll(msg, cfg.Rules, sensitive)
			if len(violations) == 0 {
				return
			}

			hasSensitiveDynamic := false
			if len(call.Args) > 0 {
				for _, v := range violations {
					if v.ID == rules.RSensitive && HasDynamicTail(pass, call.Args[0]) {
						hasSensitiveDynamic = true
						break
					}
				}
			}

			fixableViolationID, fixableText, hasFixableViolation := pickSingleSuggestedFix(violations, msg)

			for _, v := range violations {
				diag := analysis.Diagnostic{
					Pos:     pos,
					Message: string(v.ID) + " " + v.Message + " (" + kind + ")",
				}

				if hasSensitiveDynamic {
					if v.ID != rules.RSensitive {
						pass.Report(diag)
						continue
					}

					if prefix, ok := safePrefixForSensitive(pass, call.Args[0]); ok {
						if fixPos, fixEnd, ok2 := fixTargetForFirstArgWhole(call); ok2 {
							diag.SuggestedFixes = []analysis.SuggestedFix{
								{
									Message: "remove sensitive dynamic data from log message",
									TextEdits: []analysis.TextEdit{
										{
											Pos:     fixPos,
											End:     fixEnd,
											NewText: []byte(strconv.Quote(prefix)),
										},
									},
								},
							}
						}
					}

					pass.Report(diag)
					continue
				}

				if hasFixableViolation && v.ID == fixableViolationID {
					if fixPos, fixEnd, ok2 := fixTargetForFirstArg(pass, call); ok2 {
						diag.SuggestedFixes = []analysis.SuggestedFix{
							{
								Message: "apply fix for " + string(v.ID),
								TextEdits: []analysis.TextEdit{
									{
										Pos:     fixPos,
										End:     fixEnd,
										NewText: []byte(strconv.Quote(fixableText)),
									},
								},
							},
						}
					}
				}

				pass.Report(diag)
			}
		})

		return nil, nil
	}
}

func pickSingleSuggestedFix(vs []rules.Violation, msg string) (rules.RuleID, string, bool) {
	order := []rules.RuleID{rules.RSensitive, rules.RNoEmojiSpecial, rules.RLowercaseStart}

	present := map[rules.RuleID]struct{}{}
	for _, v := range vs {
		present[v.ID] = struct{}{}
	}

	for _, id := range order {
		if _, ok := present[id]; !ok {
			continue
		}
		fixed, ok := suggestFixForViolation(id, msg)
		if !ok {
			continue
		}
		if fixed == msg {
			continue
		}
		return id, fixed, true
	}
	return "", "", false
}

// detectLoggerCall определяет, является ли вызов CallExpr логированием через slog или zap.
func detectLoggerCall(pass *analysis.Pass, call *ast.CallExpr) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	if pkgIdent, ok := sel.X.(*ast.Ident); ok {
		if pkgName, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName); ok {
			if pkgName.Imported() != nil && pkgName.Imported().Path() == "log/slog" {
				if logMethods[sel.Sel.Name] {
					return "slog"
				}
			}
		}
	}

	recvT := pass.TypesInfo.TypeOf(sel.X)
	if recvT == nil {
		return ""
	}
	named := derefNamed(recvT)
	if named == nil || named.Obj() == nil || named.Obj().Pkg() == nil {
		return ""
	}
	pkgPath := named.Obj().Pkg().Path()
	typeName := named.Obj().Name()

	if logMethods[sel.Sel.Name] {
		// slog: *slog.Logger
		if pkgPath == "log/slog" && typeName == "Logger" {
			return "slog"
		}
		// zap: *zap.Logger
		if pkgPath == "go.uber.org/zap" && typeName == "Logger" {
			return "zap"
		}
		// zap: *zap.SugaredLogger
		if pkgPath == "go.uber.org/zap" && typeName == "SugaredLogger" {
			return "zap-sugar"
		}
	}

	return ""
}

// derefNamed убирает указатель и возвращает именованный тип
func derefNamed(t types.Type) *types.Named {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	if n, ok := t.(*types.Named); ok {
		return n
	}
	return nil
}

// extractFirstStringArg извлекает строковое сообщение логгера.
func extractFirstStringArg(pass *analysis.Pass, call *ast.CallExpr) (msg string, pos token.Pos, ok bool) {
	if len(call.Args) == 0 {
		return "", token.NoPos, false
	}

	expr := call.Args[0]
	s, ok := extractStaticText(pass, expr)
	if !ok {
		return "", token.NoPos, false
	}
	return s, expr.Pos(), true
}

// extractStaticText извлекает статический текст из выражения.
func extractStaticText(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	switch e := expr.(type) {

	// "literal"
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", false
		}
		return s, true

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}

		if pass != nil && pass.TypesInfo != nil {
			if t := pass.TypesInfo.TypeOf(e); t != nil && t.String() != "string" {
				return "", false
			}
		}

		if left, ok := extractStaticText(pass, e.X); ok {
			if right, ok2 := extractStaticText(pass, e.Y); ok2 {
				return left + right, true
			}
			return left, true
		}
		if right, ok := extractStaticText(pass, e.Y); ok {
			return right, true
		}
		return "", false

	case *ast.CallExpr:
		if isFmtSprintf(pass, e) {
			if len(e.Args) == 0 {
				return "", false
			}
			return extractStaticText(pass, e.Args[0])
		}
		return "", false

	case *ast.ParenExpr:
		return extractStaticText(pass, e.X)

	default:
		return "", false
	}
}

// isFmtSprintf проверяет, что выражение — это именно fmt.Sprintf.
func isFmtSprintf(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	if sel.Sel == nil || sel.Sel.Name != "Sprintf" {
		return false
	}

	id, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if pass == nil || pass.TypesInfo == nil {
		return false
	}

	pkgName, ok := pass.TypesInfo.Uses[id].(*types.PkgName)
	if !ok || pkgName.Imported() == nil {
		return false
	}

	return pkgName.Imported().Path() == "fmt"
}

// HasDynamicTail возвращает true, если первый аргумент потенциально добавляет динамические данные.
func HasDynamicTail(pass *analysis.Pass, expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		_, lok := extractStaticText(pass, e.X)
		_, rok := extractStaticText(pass, e.Y)
		return lok != rok

	case *ast.CallExpr:
		return isFmtSprintf(pass, e) && len(e.Args) > 1

	default:
		return false
	}
}

// safePrefixForSensitive строит безопасный префикс, который НЕ печатает значение секрета.
func safePrefixForSensitive(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", false
		}
		s, ok := extractStaticText(pass, e)
		if !ok {
			return "", false
		}
		return s, true

	case *ast.CallExpr:
		if isFmtSprintf(pass, e) && len(e.Args) > 0 {
			formatStr, ok := extractStaticText(pass, e.Args[0])
			if !ok {
				return "", false
			}
			prefix := stripFmtDirectives(formatStr)
			return prefix, true
		}
		return "", false

	default:
		return "", false
	}
}

// stripFmtDirectives отрезает format string по первому реальному спецификатору '%'
func stripFmtDirectives(format string) string {
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		// "%%" -> literal '%'
		if i+1 < len(format) && format[i+1] == '%' {
			i++
			continue
		}
		return format[:i]
	}
	return format
}

// fixTargetForFirstArgWhole возвращает диапазон исходника, который нужно заменить
func fixTargetForFirstArgWhole(call *ast.CallExpr) (pos, end token.Pos, ok bool) {
	if len(call.Args) == 0 {
		return token.NoPos, token.NoPos, false
	}
	return call.Args[0].Pos(), call.Args[0].End(), true
}
