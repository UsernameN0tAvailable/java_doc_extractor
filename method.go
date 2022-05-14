package main

import (
	"strings"
)

type Method struct {
	Signature string `json:"signature"`
	Documentation string `json:"documentation"`
	signatureStart int
	Line int `json:"line"`
	Body string  `json:"code"`
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
	return Method{Signature: s, Documentation: d, signatureStart: signatureStart, Line: signatureLineStart}
}
