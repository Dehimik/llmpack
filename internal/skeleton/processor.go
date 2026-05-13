package skeleton

import (
	"path/filepath"
)

// Process strategy
func Process(filename string, content []byte) ([]byte, error) {
	ext := filepath.Ext(filename)

	switch ext {
	case ".go":
		return reduceGo(content)
	case ".py":
		return reduceIndentation(content), nil
	case ".js", ".ts", ".tsx", ".jsx", ".java", ".cpp", ".c", ".h", ".hpp", ".cs", ".rs", ".php", ".swift", ".kt", ".dart", ".scala":
		return reduceBraces(content), nil
	default:
		return content, nil
	}
}
