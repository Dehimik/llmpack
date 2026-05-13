package skeleton

import (
	"bytes"
	"strings"
)

// reduceBraces strips function/class bodies in brace-based languages (JS, TS, Java, etc.)
func reduceBraces(content []byte) []byte {
	blocks := scanBraces(content)
	if len(blocks) == 0 {
		return content
	}

	lines := strings.Split(string(content), "\n")
	var buf bytes.Buffer

	// We want to hide bodies of top-level blocks (Level 1)
	// But wait, if it's a class, we might want to keep its members but hide their bodies.
	// The previous implementation was hiding everything inside the first { } it found.
	
	hideRanges := make(map[int]int) // startLine -> endLine
	for _, b := range blocks {
		if b.Level == 1 {
			hideRanges[b.StartLine] = b.EndLine
		}
	}

	for i := 0; i < len(lines); i++ {
		lineNum := i + 1
		line := lines[i]

		if endLine, ok := hideRanges[lineNum]; ok {
			// Find where { is
			idx := strings.Index(line, "{")
			if idx != -1 {
				buf.WriteString(line[:idx+1])
				buf.WriteString("`... implementation hidden ...`}")
				
				// Find if there is anything after } on the endLine
				endLineText := lines[endLine-1]
				lastIdx := strings.LastIndex(endLineText, "}")
				if lastIdx != -1 && lastIdx+1 < len(endLineText) {
					buf.WriteString(endLineText[lastIdx+1:])
				}
				buf.WriteString("\n")
				i = endLine - 1 // Skip lines
				continue
			}
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
}

// reduceIndentation strips function/class bodies in indentation-based languages (Python)
func reduceIndentation(content []byte) []byte {
	blocks := scanIndentation(content)
	if len(blocks) == 0 {
		return content
	}

	lines := strings.Split(string(content), "\n")
	var buf bytes.Buffer

	hideRanges := make(map[int]int)
	for _, b := range blocks {
		if b.Level == 1 {
			hideRanges[b.StartLine] = b.EndLine
		}
	}

	for i := 0; i < len(lines); i++ {
		lineNum := i + 1
		line := lines[i]

		if endLine, ok := hideRanges[lineNum]; ok {
			buf.WriteString(line)
			buf.WriteString("\n")
			
			indent := getIndent(line)
			buf.WriteString(strings.Repeat(" ", indent+4) + "/* ... implementation hidden ... */\n")
			
			i = endLine - 1
			continue
		}
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	return bytes.TrimSuffix(buf.Bytes(), []byte("\n"))
}
