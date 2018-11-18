package searchTree

type SearchTree struct {
	root *node
}

type node struct {
	children map[rune]*node
	name     rune
	result   *resultSet
}

type resultSet struct {
	data map[string]bool
}

func newResultSet() *resultSet {
	return &resultSet{make(map[string]bool)}
}

func (r *resultSet) add(res string) {
	r.data[res] = true
}

func (r *resultSet) addAll(other *resultSet) {
	for k, _ := range other.data {
		r.data[k] = true
	}
}

func (r *resultSet) contains(item string) bool {
	_, res := r.data[item]
	return res
}

func (r *resultSet) GetRes() map[string]bool {
	return r.data
}

func (r *resultSet) GetResSlice() []string {
	slice := make([]string, len(r.data))

	i := 0
	for k := range r.data {
		slice[i] = k
		i++
	}
	return slice
}

func creatNode(r rune) *node {
	n := &node{children: make(map[rune]*node), name: r, result: newResultSet()}
	return n
}

func MakeSearchTree() *SearchTree {
	return &SearchTree{root: creatNode(0)}
}

func (s *SearchTree) AddContent(content []string, result string) {
	for _, str := range content {
		s.AddString(str, result)
	}
}

func (s *SearchTree) AddString(str string, result string) {
	tokens := Tokenize(str)
	for _, token := range tokens {
		s.addToken(token, result)
	}
}

func (s *SearchTree) addToken(token string, result string) {
	currentNode := s.root
	for _, r := range token {
		next, exists := currentNode.children[r]
		if !exists {
			next = creatNode(r)
			currentNode.children[r] = next
		}
		currentNode = next
	}
	currentNode.result.add(result)
}

func (s *SearchTree) Search(query string, prefix bool) *resultSet {
	result := newResultSet()
	tokens := Tokenize(query)
	if len(tokens) > 1 {
		// multiple search words currently not supported
		return result
	}
	if len(tokens) == 0 {
		if prefix {
			return collectResults(s.root)
		} else {
			return result
		}
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
	if prefix {
		return collectResults(currentNode)
	} else {
		return currentNode.result
	}
}

func collectResults(n *node) *resultSet {
	result := n.result
	for _, child := range n.children {
		result.addAll(collectResults(child))
	}
	return result
}
