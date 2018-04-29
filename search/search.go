package search

type search struct {
	root *node
}

type node struct {
	children map[rune]*node
	name     rune
}

func creatNode(r rune) *node {
	n := &node{children: make(map[rune]*node), name: r}
	return n
}

func MakeSearch() *search {
	return &search{root: creatNode(0)}
}

func (s *search) AddContent(content []string) {
	for _, str := range content {
		s.AddString(str)
	}
}

func (s *search) AddString(str string) {
	tokens := Tokenize(str)
	for _, token := range tokens {
		s.addToken(token)
	}
}

func (s *search) addToken(token string) {
	currentNode := s.root
	for _, r := range token {
		next, exists := currentNode.children[r]
		if !exists {
			next = creatNode(r)
			currentNode.children[r] = next
		}
		currentNode = next
	}
}

func (s *search) Search(query string) bool {
	tokens := Tokenize(query)
	if len(tokens) > 1 {
		// multiple search words currently not supported
		return false
	}
	token := tokens[0]
	currentNode := s.root
	for _, r := range token {
		next, exists := currentNode.children[r]
		if !exists {
			return false
		}
		currentNode = next
	}
	return true
}
