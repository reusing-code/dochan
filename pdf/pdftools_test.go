package pdf

import (
	"fmt"
	"strings"
	"testing"
)

var expectedString = []struct {
	cont string
}{
	{cont: "Ihr Unternehmen"},
	{cont: "12345 Ihr Ort"},
	{cont: "CDEFGH"},
	{cont: "ZIELE"},
	{cont: "abcdefghijklm"},
	{cont: "FFFAAA"},
	{cont: "TestmeilensteinABC"},
}

func TestParsePDF(t *testing.T) {
	doc, err := parsePDF("testdata/Projektvorschlag.pdf")
	if err != nil {
		t.Fatalf("error parsing pdf: %q", err)
	}
	str := fmt.Sprint(doc)

	for _, test := range expectedString {
		if !strings.Contains(str, test.cont) {
			t.Errorf("String not in doc. Expected '%s'", test.cont)
		}
	}

}
