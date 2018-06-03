package search

type Search struct {
	root *node
}

type node struct {
	children map[rune]*node
	name     rune
	result   []interface{}
}

func creatNode(r rune) *node {
	n := &node{children: make(map[rune]*node), name: r}
	return n
}

func MakeSearch() *Search {
	return &Search{root: creatNode(0)}
}

func (s *Search) AddContent(content []string, result interface{}) {
	for _, str := range content {
		s.AddString(str, result)
	}
}

func (s *Search) AddString(str string, result interface{}) {
	tokens := Tokenize(str)
	for _, token := range tokens {
		s.addToken(token, result)
	}
}

func (s *Search) addToken(token string, result interface{}) {
	currentNode := s.root
	for _, r := range token {
		next, exists := currentNode.children[r]
		if !exists {
			next = creatNode(r)
			currentNode.children[r] = next
		}
		currentNode = next
	}
	currentNode.result = append(currentNode.result, result)
}

func (s *Search) Search(query string) []interface{} {
	result := make([]interface{}, 0)
	tokens := Tokenize(query)
	if len(tokens) > 1 {
		// multiple search words currently not supported
		return result
	}
	token := tokens[0]
	currentNode := s.root
	for _, r := range token {
		next, exists := currentNode.children[r]
		if !exists {
			return result
		}
		currentNode = next
	}
	if currentNode.result == nil {
		return result
	}
	return currentNode.result
}
