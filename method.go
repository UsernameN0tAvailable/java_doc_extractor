package main

import (
	"strings"
)

var visibility = []string {"public", "private", "protected"}

type Method struct {
	Signature string `json:"signature"`
	Documentation string `json:"documentation"`
	signatureStart int
	Line int `json:"line"`
	Body string  `json:"code"`
	HasJavaDocComment bool `json:"hasJavaDocComment"`
	DocumentableItems int `json:"documentableItems"`
	DocumentedItems int `json:"documentedItems"`
	WordsInJavaDoc int
}

func (m *Method) AddBody (body string, currentIndex int) {
	m.Body = body[m.signatureStart:currentIndex]
}


func (m*Method) GetDoc() string {
	return m.Documentation
}

func (m*Method) GetSignature() string {
	return m.Signature
}

func (m*Method) IsStatic() bool {
	s := strings.Fields(m.Signature)

	for _, v := range s {
		if v == "static"{
			return true
		}

	}
	return false
}

func (m*Method) GetName() string {

	s := strings.Split(m.Signature, "(")

	if len(s) == 0 {
		return m.Signature
	}

	sp := strings.Fields(s[0])

	return sp[len(sp) - 1]
}
 
func NewMethod(s string, d string, signatureStart int, signatureLineStart int) Method {

	documentedItems, documentableItems := metrics(s, d)

	javaDocWords := 0

	if documentedItems > 0 {
		javaDocWords = wordsInJavaDoc(d)
	}

	return Method{
		Signature: s, 
		Documentation: d, 
		signatureStart: 
		signatureStart, 
		Line: signatureLineStart, 
		DocumentedItems: documentedItems,
		DocumentableItems: documentableItems,
		HasJavaDocComment: documentedItems > 0,
		WordsInJavaDoc: javaDocWords,
	}
}

func metrics(s string, doc string) (int, int) {

	documentedItems := 0
	signature := strings.Fields(s)

	returnType := findReturnValue(signature)

	documentableItems := countParams(s)

	if returnType != "void" {
		documentableItems++

		if len(strings.Split(doc, "@return")) > 1 {
			documentedItems++
		}

	}

	throws := sliceContains(signature, "throws")


	if throws {
		documentableItems++

		if len(strings.Split(doc, "@throws")) > 1 {
			documentedItems++
		}
	}

	documentedItems += len(strings.Split(doc, "@param")) -1


	return documentedItems, documentableItems
}

func wordsInJavaDoc(d string) int {

	c := strings.Join(strings.Split(d, "/*"), "")
	c = strings.Join(strings.Split(c, "*/"), "")
	c = strings.Join(strings.Split(c, "*"), "")
	c = strings.Join(strings.Split(c, "\n"), "")

	return len(strings.Fields(c))
}


func countParams(signature string) int {

	content := []byte(signature)

	parser := Parser{}

	start := 0
	end := 0

	for i, _ := range content {
		result := parser.Parse(content, i)

		if result == EnterParamsScope && parser.ParamScopeCount == 1 {
			start = i
		} else if result == LeaveParamsScope && parser.ParamScopeCount == 0 {
			end = i
			break
		}
	}

	if start + 1 == end {
		return 0
	}

	c := string(content[start:end])
	return len(strings.Split(c, ","))
}


func findReturnValue(signature []string) string {
	for i, e := range signature {
		if sliceContains(visibility, e) && len(signature) > i + 1 {
			if signature[i + 1] == "abstract" && len(signature) > i + 2 {
				return signature[i + 2]
			} else {
				return signature[i+1]
		}
		}
	}

	if signature[0] == "abstract" && len(signature) > 1 {
		return signature[1]
	}

	return signature[0]
}

func sliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}


