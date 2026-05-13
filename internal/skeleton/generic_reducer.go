package skeleton

import (
	"bytes"
	"strings"
)

// reduceBraces strips function/class bodies in brace-based languages (JS, TS, Java, etc.)
func reduceBraces(content []byte) []byte {
	var buf bytes.Buffer
	lines := strings.Split(string(content), "\n")
	
	inComment := false
	inString := false
	var stringChar byte
	braceLevel := 0
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// If we are deep in braces, we might want to skip or check for closure
		if braceLevel > 0 {
			// Process characters to find closing brace
			for i := 0; i < len(line); i++ {
				c := line[i]
				
				// String handling
				if !inComment && (c == '"' || c == '\'' || c == '`') {
					if !inString {
						inString = true
						stringChar = c
					} else if stringChar == c && (i == 0 || line[i-1] != '\\') {
						inString = false
					}
				}
				
				if inString {
					continue
				}
				
				// Comment handling
				if !inComment && i+1 < len(line) && line[i:i+2] == "/*" {
					inComment = true
					i++
					continue
				}
				if inComment && i+1 < len(line) && line[i:i+2] == "*/" {
					inComment = false
					i++
					continue
				}
				if inComment || (i+1 < len(line) && line[i:i+2] == "//") {
					break
				}
				
				if c == '{' {
					braceLevel++
				} else if c == '}' {
					braceLevel--
					if braceLevel == 0 {
						buf.WriteString("`... implementation hidden ...`}")
						if i+1 < len(line) {
							buf.WriteString(line[i+1:])
						}
						buf.WriteString("\n")
					}
				}
			}
			continue
		}

		// Look for start of a block
		if strings.Contains(trimmed, "{") && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "*") {
			idx := strings.Index(line, "{")
			buf.WriteString(line[:idx+1])
			
			// Initial balance check for the rest of the line
			braceLevel = 1
			for i := idx + 1; i < len(line); i++ {
				c := line[i]
				if c == '{' {
					braceLevel++
				} else if c == '}' {
					braceLevel--
				}
			}
			
			if braceLevel == 0 {
				buf.WriteString("`... implementation hidden ...`}")
				buf.WriteString("\n")
			}
		} else {
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
	
	return buf.Bytes()
}

// reduceIndentation strips function/class bodies in indentation-based languages (Python)
func reduceIndentation(content []byte) []byte {
	var buf bytes.Buffer
	lines := strings.Split(string(content), "\n")
	
	skipping := false
	skipIndent := 0
	
	for _, line := range lines {
		if line == "" {
			if !skipping {
				buf.WriteString("\n")
			}
			continue
		}
		
		indent := getIndent(line)
		trimmed := strings.TrimSpace(line)
		
		if skipping {
			if indent > skipIndent || trimmed == "" {
				continue
			}
			skipping = false
		}
		
		buf.WriteString(line)
		buf.WriteString("\n")
		
		// Start skipping after function or class definition
		if (strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "class ")) && strings.HasSuffix(trimmed, ":") {
			skipping = true
			skipIndent = indent
			buf.WriteString(strings.Repeat(" ", skipIndent+4) + "/* ... implementation hidden ... */\n")
		}
	}
	
	return buf.Bytes()
}

func getIndent(line string) int {
	indent := 0
	for _, r := range line {
		if r == ' ' {
			indent++
		} else if r == '\t' {
			indent += 4 // Assume 4 spaces for tab
		} else {
			break
		}
	}
	return indent
}
