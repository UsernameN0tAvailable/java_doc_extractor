package main

import (
	"strings"
)

type Scope struct {
	path string 
	Doc string  `json:"documentation"`
	visibility string
	Name string  `json:"name"`
	Super string `json:"extends"`
	Methods []Method
	Interfaces []string `json:"interfaces"`
	staticIndex int
	fullPath string
	imports Imports
	ScopeType string `json:"type"`
}

func (s*Scope) IsClass() bool {return s.ScopeType == "class"}

func (s*Scope) IsInterface() bool {return s.ScopeType == "interface"}

func (s*Scope) IsEnum() bool {return s.ScopeType == "enum"}
func (s*Scope) IsRecord() bool {return s.ScopeType == "record"}

func NewScope(fullPath string, signature string, doc string, path string, imports *Imports, scope *Scope) Scope {

	fields := strings.Fields(strings.TrimSpace(RemoveTemplate(signature)))

	classIndex := -1 
	extendIndex := -1 
	implementsIndex := -1
	staticIndex := -1
	scopeType := "" 
	for i, p := range fields {
		if p == "class" {
			classIndex = i
			scopeType = "class"
		} else if p == "interface" {
			classIndex = i
			scopeType = "interface"
		} else if p == "enum" {
			classIndex = i
			scopeType = "enum"
		} else if p == "record" {
			classIndex = i
			scopeType = "record"
		} else if p == "extends" {
			extendIndex = i
		} else if p == "implements" {
			implementsIndex = i
		} else if p == "static" {
			staticIndex = i
		}
	}

	name := fields[classIndex + 1]

	var vis string 

	if classIndex < 1 {
		vis = ""
	} else {
		vis = strings.Join(fields[:classIndex], " ")
	}

	pathSplt := strings.Split(strings.Split(path, ".java")[0], "/")


	className := strings.Join(pathSplt, ".")

	if staticIndex == -1 {
		pathSplt[len(pathSplt) - 1] = name 
		className = strings.Join(pathSplt, ".")
	}  else {
		className = scope.GetName() + "." + name
	}

	var super string

	if extendIndex < 1 {
		super = ""
	} else {
		toFind :=RemoveTemplate(fields[extendIndex + 1])
		super = imports.GetPath(toFind)
	}

	implements := make([]string, 0)

	if implementsIndex > 0 {
		tmp := strings.Join(fields[implementsIndex + 1:], " ")
		interfacesStr := strings.Split(tmp, "{")[0]

		for _,in := range strings.Split(interfacesStr, ",") {
			toFind := RemoveTemplate(strings.TrimSpace(in))
			implements = append(implements, imports.GetPath(toFind))
		}	
	}

	return Scope{
		ScopeType: scopeType,
		fullPath: fullPath,
		staticIndex: staticIndex,
		path: path,
		Doc: doc, 
		visibility: vis,
		Name: strings.TrimSpace(className),
		Super: super,
		Interfaces: implements,
		Methods: make([]Method, 0, 20),
		imports: *imports,
	} 
}

func (c*Scope) IsInPackage(packSearched string) bool {
	pack, err := c.imports.GetPackage()
	if err != nil {
		return false
	}

	return pack == packSearched
}

func (c*Scope) GetPackage() (string, error) {
	return c.imports.GetPackage()
}

func (c*Scope) GetMethods() []Method {
	return c.Methods
}

func (c * Scope) AppendMethod(m Method) {
	c.Methods = append(c.Methods, m)
}

func (c * Scope) GetPath() string {
	return c.path
}

func (c * Scope) GetFullPath() string {
	return c.fullPath
}

func (c* Scope) GetDocLinesCount() int {
	if len(c.Doc) == 0 {
		return 0
	} 
	return len(strings.Split(c.Doc, "\n"))
}

func (c* Scope) GetDoc() string {
	return c.Doc
}

func (c* Scope) GetVis() string {
	return c.visibility
}

func (c* Scope) GetName() string {
	return c.Name
}

func (c* Scope) SetSuper(v string) {
	c.Super = v
}

func (c* Scope) GetSuper() string {
	return c.Super
}

func (c* Scope) GetInterfaces() []string {
	return c.Interfaces
}

func (c* Scope) SetInterface(v string, index int) {
	c.Interfaces[index] = v
}


//public to tst helper
func RemoveTemplate(name string) string {

	start := 0
	end := len(name) 

	count := 0

	result := ""

	for i, s := range name {

		if string(s) == "<" {

			count++
			if count == 1 {
				end = i 
			}
		} else if string(s) == ">" {
			count--

			if count == 0 {
				result += name[start:end]
				start = i + 1
				end = len(name) 
			}
		} 
	}

	result += name[start:end]

	return result
}


