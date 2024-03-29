package pdf

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"

	"strconv"

	"github.com/kokardy/saxlike"
)

type filter struct {
	in io.Reader
}

func filterBoldTags(p []byte) []byte {
	ret := bytes.Replace(p, []byte{60, 98, 62}, []byte{}, -1)
	ret = bytes.Replace(ret, []byte{60, 47, 98, 62}, []byte{}, -1)
	return ret
}

func ParseFile(xmlFile string) (*Document, error) {
	f, err := os.Open(xmlFile)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	content = filterBoldTags(content)
	buf := bytes.NewBuffer(content)

	h := pdftohtmlHander{}
	err = saxlike.Parse(buf, &h, false)
	return &h.doc, err
}

type pdftohtmlHander struct {
	saxlike.VoidHandler
	doc         Document
	currentPage *page
	currentText *textBlock
}

func (h *pdftohtmlHander) StartElement(e xml.StartElement) {
	switch e.Name.Local {
	case "page":
		atts := parseAttributes(e.Attr)
		sizeX := parseInt(atts["width"])
		sizeY := parseInt(atts["height"])
		h.currentPage = &page{sizeX: sizeX, sizeY: sizeY}
	case "text":
		atts := parseAttributes(e.Attr)
		posX := parseInt(atts["left"])
		posY := parseInt(atts["top"])
		sizeX := parseInt(atts["width"])
		sizeY := parseInt(atts["height"])
		h.currentText = &textBlock{posX: posX, posY: posY, sizeX: sizeX, sizeY: sizeY}
	}
}

func (h *pdftohtmlHander) EndElement(e xml.EndElement) {
	switch e.Name.Local {
	case "page":
		h.doc.pages = append(h.doc.pages, *h.currentPage)
		h.currentPage = nil
	case "text":
		h.currentPage.blocks = append(h.currentPage.blocks, *h.currentText)
		h.currentText = nil
	}
}

func (h *pdftohtmlHander) CharData(c xml.CharData) {
	if h.currentText != nil {
		h.currentText.text += string(c[:])
	}
}

func parseAttributes(a []xml.Attr) map[string]string {
	m := make(map[string]string)
	for _, attr := range a {
		m[attr.Name.Local] = attr.Value
	}
	return m
}

func parseInt(s string) int32 {
	if s == "" {
		return 0
	}
	ret, _ := strconv.ParseInt(s, 10, 32)
	return int32(ret)
}
