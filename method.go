package main

type Method struct {
	Signature string `json:"signature"`
	Documentation string `json:"documentation"`
}


func (m*Method) GetDoc() string {
	return m.Documentation
}

func (m*Method) GetSignature() string {
	return m.Signature
}

 
func NewMethod(s string, d string) Method {
	return Method{Signature: s, Documentation: d}
}
