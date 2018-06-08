package search

import (
	"testing"

	"github.com/reusing-code/dochan/search/wiki_testdata"
)

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
	result []interface{}
}{
	{query: "language", result: []interface{}{"en"}},
	{query: "Griesemer", result: []interface{}{"en", "ger"}},
	{query: "griesemer", result: []interface{}{"en", "ger"}},
	{query: "Gries", result: []interface{}{}},
	{query: "Griesa", result: []interface{}{}},
	{query: "riesemer", result: []interface{}{}},
	{query: "Kanäle", result: []interface{}{"ger"}},
	{query: "kanale", result: []interface{}{"ger"}},
	{query: "känäle", result: []interface{}{"ger"}},
	{query: "Kanele", result: []interface{}{}},
	{query: "Kanöle", result: []interface{}{}},
	{query: "Kanaele", result: []interface{}{"ger"}},
	{query: "Kan äle", result: []interface{}{}},
	{query: "2009", result: []interface{}{"en"}},
	{query: "2010", result: []interface{}{}},
	{query: "C++", result: []interface{}{"en"}},
	{query: "", result: []interface{}{}},
	{query: "C+++++---/(&", result: []interface{}{"en"}}, // hm...
}

func TestSearch(t *testing.T) {
	s := MakeSearch()

	s.AddContent(testDataEn, "en")
	s.AddContent(testDataGer, "ger")

	for _, tc := range searchTests {
		res := s.Search(tc.query, false)
		for _, val := range tc.result {
			if !res.contains(val) {
				t.Errorf("Query %q resulted in wrong result. Want %v have %v", tc.query, val, res)
			}
		}
		if len(tc.result) != len(res.data) {
			t.Errorf("Query %q resulted in wrong number of results. Want %v have %v", tc.query, len(tc.result), len(res.data))
		}
	}
}

var prefixSearchTests = []struct {
	query  string
	result []interface{}
}{
	{query: "Griesemer", result: []interface{}{"en", "ger"}},
	{query: "griesemer", result: []interface{}{"en", "ger"}},
	{query: "Gries", result: []interface{}{"en", "ger"}},
	{query: "Griese", result: []interface{}{"en", "ger"}},
	{query: "g", result: []interface{}{"en", "ger"}},
	{query: "", result: []interface{}{"en", "ger"}},
	{query: "Griesea", result: []interface{}{}},
}

func TestPrefixSearch(t *testing.T) {
	s := MakeSearch()

	s.AddContent(testDataEn, "en")
	s.AddContent(testDataGer, "ger")

	for _, tc := range prefixSearchTests {
		res := s.Search(tc.query, true)
		for _, val := range tc.result {
			if !res.contains(val) {
				t.Errorf("Query %q resulted in wrong result. Want %v have %v", tc.query, val, res)
			}
		}
	}
}

var benchSearchWordsDE = []string{
	"friedrich",
	"dampfschiff",
	"keinWorttttttt",
	"und",
	"sandkast",
	"bei",
}

func BenchmarkWikiDataDE(b *testing.B) {
	s := MakeSearch()
	wiki_testdata.ParseDataDE(func(data string) {
		s.AddString(data, "test")
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, word := range benchSearchWordsDE {
			s.Search(word, false)
		}
	}
}
