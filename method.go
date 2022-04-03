package main

type Method struct {
	signature string
	documentation string
}


func (m*Method) GetDoc() string {
	return m.documentation
}

func (m*Method) GetSignature() string {
	return m.signature
}

 
func NewMethod(s string, d string) Method {
	return Method{signature: s, documentation: d}
}
