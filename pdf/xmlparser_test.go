package pdf

import (
	"bytes"
	"testing"

	"github.com/kokardy/saxlike"
)

var intTests = []struct {
	in  string
	out int32
}{
	{"11", 11},
	{"0", 0},
	{"", 0},
	{"abc", 0},
	{"555", 555},
}

func TestParseInt(t *testing.T) {
	for _, test := range intTests {
		if i := parseInt(test.in); i != test.out {
			t.Errorf("Wrong parsed value. Expected %d, was %d", test.out, i)
		}
	}

}

var xmlSample string = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE pdf2xml>

<pdf2xml producer="poppler" version="0.48.0">
<page number="1" position="absolute" top="0" left="0" height="1080" width="1920">
	<fontspec id="0" size="14" family="Times" color="#000000"/>
<text top="10" left="8" width="20" height="10000" font="0">Text 1</text>
<text top="30" left="4" width="200" height="11100" font="0">Sample 2</text>
<text top="50" left="2" width="2000" height="11111" font="0">XML3</text>
</page>
</pdf2xml>`

var xmlSampleElements = []textBlock{
	{posY: 10, posX: 8, sizeX: 20, sizeY: 10000, text: "Text 1"},
	{posY: 30, posX: 4, sizeX: 200, sizeY: 11100, text: "Sample 2"},
	{posY: 50, posX: 2, sizeX: 2000, sizeY: 11111, text: "XML3"},
}

func TestXmlSample(t *testing.T) {
	h := pdftohtmlHander{}
	r := bytes.NewReader([]byte(xmlSample))
	err := saxlike.Parse(r, &h, false)
	if err != nil {
		t.Error(err)
	}
	for i, block := range h.doc.pages[0].blocks {
		if block != xmlSampleElements[i] {
			t.Errorf("Textblock wrong. Expected '%v', was '%v'", xmlSampleElements[i], block)
		}
	}
}
