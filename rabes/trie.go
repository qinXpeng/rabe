package rabes

import "fmt"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (nd *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", nd.pattern, nd.part, nd.isWild)
}

func (nd *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		nd.pattern = pattern
		return
	}

	part := parts[height]
	child := nd.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		nd.children = append(nd.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (nd *node) matchChild(part string) *node {
	for _, child := range nd.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

