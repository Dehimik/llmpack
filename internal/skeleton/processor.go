package skeleton

import (
	"path/filepath"
)

// Process strategy
func Process(filename string, content []byte) ([]byte, error) {
	return ProcessSpecific(filename, content, "")
}

// ProcessSpecific reduces code but preserves targetSymbol implementation
func ProcessSpecific(filename string, content []byte, targetSymbol string) ([]byte, error) {
	ext := filepath.Ext(filename)

	switch ext {
	case ".go":
		return reduceGo(content, targetSymbol)
	case ".py":
		return reduceIndentation(content, targetSymbol), nil
	case ".js", ".ts", ".tsx", ".jsx", ".java", ".cpp", ".c", ".h", ".hpp", ".cs", ".rs", ".php", ".swift", ".kt", ".dart", ".scala":
		return reduceBraces(content, targetSymbol), nil
	default:
		return content, nil
	}
}
