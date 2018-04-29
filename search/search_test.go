package search

import "testing"

/**
   	The English test data originates from the Wikipedia article
   	'Go (programming language)' ( https://en.wikipedia.org/wiki/Go_(programming_language) ),
   	which is released under the Creative Commons Attribution-Share-Alike License 3.0 ( https://creativecommons.org/licenses/by-sa/3.0/ ).

	The German test data originates from the Wikipedia article
   	'Go (Programmiersprache)' ( https://de.wikipedia.org/wiki/Go_(Programmiersprache) ),
   	which is released under the Creative Commons Attribution-Share-Alike License 3.0 ( https://creativecommons.org/licenses/by-sa/3.0/ ).
*/
var testDataEn = []string{
	"Go (often referred to as Golang) is a programming language created at Google[10] in 2009 by Robert Griesemer, Rob Pike, and Ken Thompson.",
	"Statically typed and scalable to large systems (like Java or C++)",
	"For a pair of types K, V, the type map[K]V is the type of hash tables mapping type-K keys to type-V values.",
}

var testDataGer = []string{
	"Go unterstützt objektorientierte Programmierung, diese ist jedoch nicht klassenbasiert.",
	"Zur Unterstützung der nebenläufigen Programmierung in Go wird das Konzept der Kanäle (channels) genutzt, ",
	"Die Entwürfe stammen von Robert Griesemer, Rob Pike und Ken Thompson.",
}

var searchTests = []struct {
	query  string
	result bool
}{
	{query: "language", result: true},
	{query: "Griesemer", result: true},
	{query: "griesemer", result: true},
	{query: "Gries", result: true},
	{query: "Griesa", result: false},
	{query: "riesemer", result: false},
	{query: "Kanäle", result: true},
	{query: "kanale", result: true},
	{query: "känäle", result: true},
	{query: "Kanele", result: false},
	{query: "Kanaele", result: false},
	{query: "Kan äle", result: false},
}

func TestSearch(t *testing.T) {
	s := MakeSearch()

	s.AddContent(testDataEn)
	s.AddContent(testDataGer)

	for _, tc := range searchTests {
		res := s.Search(tc.query)
		if res != tc.result {
			t.Errorf("Query %q resulted in wrong result. Want %t have %t", tc.query, tc.result, res)
		}
	}
}
