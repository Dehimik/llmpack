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
	// In future .ts, .py, .java
	default:
		return content, nil
	}
}
