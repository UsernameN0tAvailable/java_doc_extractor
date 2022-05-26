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
	ReturnType string
	Throws bool
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

	signatureFields := strings.Fields(s)


	return Method{
		Signature: s, 
		Documentation: d, 
		signatureStart: 
		signatureStart, 
		Line: signatureLineStart, 
		Throws: sliceContains(signatureFields, "throws"), 
		ReturnType: findReturnValue(signatureFields) }
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


