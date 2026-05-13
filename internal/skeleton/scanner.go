package skeleton

import (
	"strings"
)

// Block represents a range of lines (1-indexed)
type Block struct {
	StartLine int
	EndLine   int
	Level     int
}

// scanBraces finds blocks delimited by { }
func scanBraces(content []byte) []Block {
	var blocks []Block
	lines := strings.Split(string(content), "\n")
	
	inComment := false
	inString := false
	var stringChar byte
	
	type openBlock struct {
		startLine int
		level     int
	}
	var stack []openBlock
	
	for lineIdx, line := range lines {
		lineNum := lineIdx + 1
		
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
				stack = append(stack, openBlock{startLine: lineNum, level: len(stack) + 1})
			} else if c == '}' {
				if len(stack) > 0 {
					top := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					blocks = append(blocks, Block{
						StartLine: top.startLine,
						EndLine:   lineNum,
						Level:     top.level,
					})
				}
			}
		}
	}
	
	return blocks
}

// scanIndentation finds blocks in indentation-based languages
func scanIndentation(content []byte) []Block {
	var blocks []Block
	lines := strings.Split(string(content), "\n")
	
	type openBlock struct {
		startLine int
		indent    int
	}
	var stack []openBlock
	
	for lineIdx, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		
		lineNum := lineIdx + 1
		indent := getIndent(line)
		trimmed := strings.TrimSpace(line)
		
		// Close blocks if current indent is less than or equal to stack top
		for len(stack) > 0 && indent <= stack[len(stack)-1].indent {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			
			// Find actual end line (previous non-empty line)
			endLine := lineNum - 1
			for endLine > top.startLine && strings.TrimSpace(lines[endLine-1]) == "" {
				endLine--
			}
			
			blocks = append(blocks, Block{
				StartLine: top.startLine,
				EndLine:   endLine,
				Level:     len(stack) + 1,
			})
		}
		
		// Start skipping after function or class definition
		if (strings.HasPrefix(trimmed, "def ") || strings.HasPrefix(trimmed, "class ")) && strings.HasSuffix(trimmed, ":") {
			stack = append(stack, openBlock{
				startLine: lineNum,
				indent:    indent,
			})
		}
	}
	
	// Close remaining blocks
	lastLine := len(lines)
	for lastLine > 0 && strings.TrimSpace(lines[lastLine-1]) == "" {
		lastLine--
	}
	
	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		blocks = append(blocks, Block{
			StartLine: top.startLine,
			EndLine:   lastLine,
			Level:     len(stack) + 1,
		})
	}
	
	return blocks
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
