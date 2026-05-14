package skeleton

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

func reduceGo(content []byte, targetSymbol string) ([]byte, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			// Check if this is the target symbol
			name := fn.Name.Name
			qualifiedName := name
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				t := fn.Recv.List[0].Type
				var parent string
				if star, ok := t.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						parent = ident.Name
					}
				} else if ident, ok := t.(*ast.Ident); ok {
					parent = ident.Name
				}
				if parent != "" {
					qualifiedName = parent + "." + name
				}
			}

			if targetSymbol != "" && (name == targetSymbol || qualifiedName == targetSymbol) {
				return true // Preserve
			}

			if fn.Body != nil {
				fn.Body.List = []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "`... implementation hidden ...`", // Or just comment
						},
					},
				}
			}
		}
		return true
	})

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
