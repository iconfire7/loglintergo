package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/iconfire7/loglintergo/internal/rules"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglintergo",
	Doc:  "checks log messages for style/safety rules",
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
	Run: run,
}

var logMethods = map[string]bool{
	"Debug": true,
	"Info":  true,
	"Warn":  true,
	"Error": true,
}

// run — основная функция анализа пакета.
func run(pass *analysis.Pass) (any, error) {
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

		violations := rules.CheckAll(msg)

		for _, v := range violations {
			pass.Reportf(
				pos,
				"%s %s (%s)",
				v.ID,
				v.Message,
				kind,
			)
		}
	})

	return nil, nil
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
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", token.NoPos, false
	}

	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", token.NoPos, false
	}
	return s, lit.Pos(), true
}
