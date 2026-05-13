package skeleton

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
)

type Symbol struct {
	Name      string
	Parent    string
	Type      string
	StartLine int
	EndLine   int
}

func ExtractSymbols(filename string, content []byte) ([]Symbol, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".go":
		return extractGoSymbols(content)
	case ".py":
		return extractIndentationSymbols(content)
	case ".js", ".ts", ".tsx", ".jsx", ".java", ".cpp", ".c", ".h", ".hpp", ".cs", ".rs", ".php", ".swift", ".kt", ".dart", ".scala":
		return extractBraceSymbols(content)
	default:
		return nil, nil
	}
}

func extractGoSymbols(content []byte) ([]Symbol, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var symbols []Symbol
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				symbols = append(symbols, Symbol{
					Name:      x.Name.Name,
					Type:      "struct",
					StartLine: fset.Position(x.Pos()).Line,
					EndLine:   fset.Position(x.End()).Line,
				})
			}
		case *ast.FuncDecl:
			sym := Symbol{
				Name:      x.Name.Name,
				Type:      "function",
				StartLine: fset.Position(x.Pos()).Line,
				EndLine:   fset.Position(x.End()).Line,
			}
			if x.Recv != nil && len(x.Recv.List) > 0 {
				sym.Type = "method"
				// Handle receiver name
				t := x.Recv.List[0].Type
				if star, ok := t.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						sym.Parent = ident.Name
					}
				} else if ident, ok := t.(*ast.Ident); ok {
					sym.Parent = ident.Name
				}
			}
			symbols = append(symbols, sym)
		}
		return true
	})
	return symbols, nil
}

var (
	pyClassRegex = regexp.MustCompile(`^class\s+([a-zA-Z0-9_]+)`)
	pyDefRegex   = regexp.MustCompile(`^def\s+([a-zA-Z0-9_]+)`)
	brClassRegex = regexp.MustCompile(`class\s+([a-zA-Z0-9_]+)`)
	brFuncRegex  = regexp.MustCompile(`function\s+([a-zA-Z0-9_]+)`)
	brMethRegex  = regexp.MustCompile(`([a-zA-Z0-9_]+)\s*\([^)]*\)\s*\{`)
)

func extractIndentationSymbols(content []byte) ([]Symbol, error) {
	blocks := scanIndentation(content)
	lines := strings.Split(string(content), "\n")
	var symbols []Symbol
	
	for _, b := range blocks {
		line := strings.TrimSpace(lines[b.StartLine-1])
		var name, symType string
		
		if matches := pyClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			name = matches[1]
			symType = "class"
		} else if matches := pyDefRegex.FindStringSubmatch(line); len(matches) > 1 {
			name = matches[1]
			symType = "function"
		}
		
		if name != "" {
			symbols = append(symbols, Symbol{
				Name:      name,
				Type:      symType,
				StartLine: b.StartLine,
				EndLine:   b.EndLine,
			})
		}
	}
	// Sort by start line so we can find parents easily
	for i := 0; i < len(symbols); i++ {
		for j := i + 1; j < len(symbols); j++ {
			if symbols[i].StartLine > symbols[j].StartLine {
				symbols[i], symbols[j] = symbols[j], symbols[i]
			}
		}
	}
	// Find parents
	for i := range symbols {
		for j := i - 1; j >= 0; j-- {
			if symbols[j].Type == "class" && symbols[i].StartLine > symbols[j].StartLine && symbols[i].EndLine <= symbols[j].EndLine {
				symbols[i].Parent = symbols[j].Name
				break
			}
		}
	}
	return symbols, nil
}

func extractBraceSymbols(content []byte) ([]Symbol, error) {
	blocks := scanBraces(content)
	lines := strings.Split(string(content), "\n")
	var symbols []Symbol
	
	for _, b := range blocks {
		line := strings.TrimSpace(lines[b.StartLine-1])
		var name, symType string
		
		if matches := brClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			name = matches[1]
			symType = "class"
		} else if matches := brFuncRegex.FindStringSubmatch(line); len(matches) > 1 {
			name = matches[1]
			symType = "function"
		} else if matches := brMethRegex.FindStringSubmatch(line); len(matches) > 1 {
			name = matches[1]
			symType = "method"
		}
		
		if name != "" {
			symbols = append(symbols, Symbol{
				Name:      name,
				Type:      symType,
				StartLine: b.StartLine,
				EndLine:   b.EndLine,
			})
		}
	}
	// Sort by start line so we can find parents easily
	for i := 0; i < len(symbols); i++ {
		for j := i + 1; j < len(symbols); j++ {
			if symbols[i].StartLine > symbols[j].StartLine {
				symbols[i], symbols[j] = symbols[j], symbols[i]
			}
		}
	}
	// Find parents
	for i := range symbols {
		for j := i - 1; j >= 0; j-- {
			if symbols[j].Type == "class" && symbols[i].StartLine > symbols[j].StartLine && symbols[i].EndLine <= symbols[j].EndLine {
				symbols[i].Parent = symbols[j].Name
				break
			}
		}
	}
	return symbols, nil
}
