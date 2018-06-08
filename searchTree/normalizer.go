package searchTree

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

var umlautReplacer = strings.NewReplacer("ae", "a", "oe", "o", "ue", "u")

func Tokenize(str string) []string {
	var i norm.Iter
	var result []string
	var curString string
	i.InitString(norm.NFKD, str)
	for !i.Done() {
		curRune, _ := utf8.DecodeRune(i.Next())
		if unicode.IsLetter(curRune) || unicode.IsDigit(curRune) {
			curString += string(unicode.ToLower(curRune))
		} else if unicode.IsSpace(curRune) || unicode.IsPunct(curRune) {
			if len(curString) > 0 {
				curString = umlautReplacer.Replace(curString)
				result = append(result, curString)
				curString = ""
			}
		}
	}
	if len(curString) > 0 {
		curString = umlautReplacer.Replace(curString)
		result = append(result, curString)
	}
	return result
}
