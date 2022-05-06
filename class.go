package main

import (
	//	"fmt"
//	"fmt"
	"strings"
)

type Scope struct {
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
	IsTest bool `json:"isTest"`
	Tests []string `json:"testClasses"`
	SubClasses []string `json:"subClasses"`
	ImplementedBy []string `json:"implementedBy"`
	tmpUses []string
	Uses []string `json:"uses"`
	UsedBy []string `json:"usedBy"`
	IsPrivate bool `json:"isPrivate"`
	body string
	InnerClasses []string
}

func (s*Scope) IsClass() bool {return s.ScopeType == "class"}
func (s*Scope) IsInterface() bool {return s.ScopeType == "interface"}
func (s*Scope) IsEnum() bool {return s.ScopeType == "enum"}
func (s*Scope) IsRecord() bool {return s.ScopeType == "record"}
func (s*Scope) IsATest() bool {return s.IsTest}
func (s*Scope) IsAnnotation() bool {return s.ScopeType == "annotation"}

func NewScope(fullPath string, signature string, doc string, imports *Imports, scope *Scope) Scope {

	fields := strings.Fields(strings.TrimSpace(RemoveTemplate(signature)))

	classIndex := -1 
	extendIndex := -1 
	implementsIndex := -1
	staticIndex := -1
	scopeType := "" 
	visible := false

	for i, p := range fields {
		if p == "class" {
			classIndex = i
			scopeType = "class"
		} else if p == "interface" {
			classIndex = i
			scopeType = "interface"
		} else if p == "@interface" {
			classIndex = i
			scopeType = "annotation"
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
		} else if p == "protected" || p == "public" {
			visible = true
		}
	}

	name := fields[classIndex + 1]

	var vis string 

	if classIndex < 1 {
		vis = ""
	} else {
		vis = strings.Join(fields[:classIndex], " ")
	}

	pack := imports.GetPackage()

	var className string

	if staticIndex == -1 {
		className = pack + "." + name
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

	isTest := strings.Contains(className, ".test.") || strings.Contains(className, "Test") || strings.Contains(className, "Benchmark")

	annotations := findAnnotations(signature, imports)

	return Scope{
		IsTest: isTest,
		ScopeType: scopeType,
		fullPath: fullPath,
		staticIndex: staticIndex,
		Doc: doc, 
		visibility: vis,
		Name: strings.TrimSpace(strings.Split(className, "(")[0]),
		Super: super,
		Interfaces: implements,
		Methods: make([]Method, 0, 20),
		imports: *imports,
		Tests: make([]string, 0, 20),
		SubClasses: make([]string, 0, 20),
		ImplementedBy: make([]string, 0, 20),
		tmpUses: annotations,
		Uses: make([]string, 0, 20),
		UsedBy: make([]string, 0, 20),
		IsPrivate: !visible,
		body: "",
		InnerClasses: make([]string, 0, 20),
	} 
}

func findAnnotations(content string, imports *Imports) []string {

	annotations := make([]string, 0, 10)
	signatureLines := strings.Split(content, "\n")
	for _, a := range signatureLines {
		v := strings.TrimSpace(a)
		if  len(v) > 0 && v[0] == byte('@') {
			v = strings.Split(v[1:], "(")[0]	

			found := imports.GetPath(v)

			if len(strings.Split(found, ".")) == 1 {
				found = imports.GetPackage() + "." + found
			}

			annotations = append(annotations,  found)
		}
	}

	return annotations
}


func (s *Scope) HasSuper() bool {return len(s.Super) > 0}
func (s *Scope) HasChildren() bool {return len(s.SubClasses) > 0}

func (s *Scope) AddInnerClass(class string) {
	s.InnerClasses = append(s.InnerClasses, class)
	s.imports.imports = append(s.imports.imports, class)
}

func (s *Scope) AddBody(body string, imports *Imports) {
	s.body = s.removeInnerClassesDeclarations(body)
	s.tmpUses = append(s.tmpUses, findAnnotations(body, imports)...)
}

/* So that they dont appear in use */
func (s *Scope) removeInnerClassesDeclarations(body string) string {

	for _, icName := range s.InnerClasses {

		icSplit := strings.Split(icName, ".")
		ic := icSplit[len(icSplit) - 1]

		body = strings.Join(strings.Split(body, "inteface " + ic), "")
		body = strings.Join(strings.Split(body, "class " + ic), "")
		body = strings.Join(strings.Split(body, "enum " + ic), "")
		body = strings.Join(strings.Split(body, "record " + ic), "")
	}

	return body
}

func (s *Scope) IsVisible() bool {
	return !s.IsPrivate
}

func (s* Scope) AppendUses(u string) {
	if !isStored(s.Uses, u) {
		s.Uses = append(s.Uses, u)
	}
}

func (s*Scope) AppendUsedBy(use string) {
	if !isStored(s.UsedBy, use) {
		s.UsedBy = append(s.UsedBy, use)
	}
}

func (s*Scope) AppendTestCase(testCase string) {
	if !isStored(s.Tests, testCase) {
		s.Tests = append(s.Tests, testCase)
	}
}


func (c*Scope) AppendImplementedBy(inter string) {
	if !isStored(c.ImplementedBy, inter) {
		c.ImplementedBy = append(c.ImplementedBy, inter)
	}

}

func (c*Scope) AppendSubClass(subClass string) {
	if !isStored(c.SubClasses, subClass) {
		c.SubClasses = append(c.SubClasses, subClass)
	}
}


func isStored(stack []string, hay string) bool {
	contains := false

	for _, v := range stack {
		if v == hay {
			contains = true
			break
		}

	}
	return contains
}

func (c*Scope) Imports(scope *Scope) bool {

	if c.imports.IsImported(scope) {
		return true
	}

	for _, uses := range c.tmpUses {
		for _,m := range scope.GetStaticMethods() {
			if m == uses {
				return true
			}
		}
	}


	return false
}

func (c *Scope) UsesClass(class *Scope) bool {

	if !c.Imports(class) {return false}

	className := class.GetName()

	for _, u := range c.tmpUses {

		if u == className {
			return true
		}else {
			for _, m := range class.GetStaticMethods() {
				if u == m {
					return true
				}
			}
		}
	}

	return c.imports.IsClassUsed(class, c.body)
}

func (c*Scope) IsInPackage(packSearched string) bool {
	return c.imports.IsInPackage(packSearched)
}

func (c*Scope) GetPackage() string {
	return c.imports.GetPackage()
}

func (c*Scope) GetMethods() []Method {
	return c.Methods
}

func (c*Scope) GetStaticMethods() []string {
	staticMethods := make([]string, 0, 10)

	for _, m := range c.Methods {
		if m.IsStatic() {
			staticMethods = append(staticMethods, c.GetName() + "." + m.GetName())
		}
	}
	return staticMethods

}

func (c * Scope) AppendMethod(m Method) {
	c.Methods = append(c.Methods, m)
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

	content := []byte(name)

	start := 0
	end := len(name) 

	count := 0

	result := ""

	parser := NewParser()

	for i, _ := range name {

		switch parser.Parse(content, i) {
		case EnterTemplate:
			count++
			if count == 1 {
				end = i 
			}
		case LeaveTemplate:
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


