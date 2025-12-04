package app

import (
	"sort"
	"strings"
)

// Node representing a file or directory in the tree
type Node struct {
	Name     string
	Children map[string]*Node
	IsFile   bool
}

func newNode(name string, isFile bool) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		IsFile:   isFile,
	}
}

// buildTree converts a list of paths into a Node structure
func buildTree(paths []string) *Node {
	root := newNode(".", false)

	for _, path := range paths {
		cleanPath := strings.ReplaceAll(path, "\\", "/")
		parts := strings.Split(cleanPath, "/")

		current := root
		for i, part := range parts {
			if part == "" || part == "." {
				continue
			}

			isFile := i == len(parts)-1

			if _, exists := current.Children[part]; !exists {
				current.Children[part] = newNode(part, isFile)
			}
			current = current.Children[part]
		}
	}
	return root
}

// renderTree generates the string representation
func renderTree(root *Node) string {
	var sb strings.Builder
	printNodes(&sb, root, "", true)
	return sb.String()
}

func printNodes(sb *strings.Builder, node *Node, prefix string, isRoot bool) {
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, name := range keys {
		child := node.Children[name]
		isLast := i == len(keys)-1

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		if !isRoot {
			sb.WriteString(prefix + connector + name + "\n")
		} else {
			sb.WriteString(name + "\n")
		}

		newPrefix := prefix
		if !isRoot {
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
		}

		// recursion
		printNodes(sb, child, newPrefix, false)
	}
}
