package searchTree

import (
	"reflect"
	"testing"
)

var normTests = []struct {
	input  string
	output []string
}{
	{"abc", []string{"abc"}},
	{"ABC", []string{"abc"}},
	{"AbCd", []string{"abcd"}},
	{"Äbcöß", []string{"abcoß"}},
	{"A.b!c", []string{"a", "b", "c"}},
	{"A A A", []string{"a", "a", "a"}},
	{"a2B3 01337", []string{"a2b3", "01337"}},
	{"hae hoe hue", []string{"ha", "ho", "hu"}},
}

func TestNormalization(t *testing.T) {
	for _, test := range normTests {
		normStr := Tokenize(test.input)
		if !reflect.DeepEqual(normStr, test.output) {
			t.Errorf("Normalization of '%s' failed: expected '%s', was '%s'", test.input, test.output, normStr)
		}
	}
}
