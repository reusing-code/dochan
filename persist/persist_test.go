package persist

import (
	"bytes"
	"testing"
)

var testData = []struct {
	key   string
	value []string
}{
	{"key", []string{"value"}},
	{"kay", []string{"velue"}},
	{"empty", []string{}},
	{"", []string{"empty"}},
	{"mutliple", []string{"v1", "v2", "v3"}},
}

func TestPersist(t *testing.T) {
	input := make(DataMap)
	for _, tc := range testData {
		input[tc.key] = tc.value
	}

	buf := &bytes.Buffer{}
	err := PersistData(input, buf)
	if err != nil {
		t.Fatal(err)
	}
	output, err := RetrieveData(buf)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testData {
		res, ok := output[tc.key]
		if !ok {
			t.Errorf("No result for %q", tc.key)
			continue
		}
		if len(res) != len(tc.value) {
			t.Errorf("Wrong len for %q: len(res) != len(tc.value, %d != %d", tc.key, len(res), len(tc.value))
			continue
		}
		for i, _ := range res {
			if res[i] != tc.value[i] {
				t.Errorf("Wrong value for %q[%d]: res[i] != tc.value[i], %q != %q", tc.key, i, res[i], tc.value[i])
				break
			}
		}
	}
}
