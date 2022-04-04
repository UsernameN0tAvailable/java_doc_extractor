package main

import (
	"strings"
)

type Scope interface {
	IsClass() bool
	IsInterface() bool
	AppendMethod(m Method)
	GetName() string
}


type Class struct {
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
}

func (c*Class) IsClass() bool {return true}
func (c*Class) IsInterface() bool {return false}

func NewClass(fullPath string, signature string, doc string, path string, imports *Imports, scope Scope) Class {

	fields := strings.Fields(strings.TrimSpace(RemoveTemplate(signature)))

	classIndex := -1 
	extendIndex := -1 
	implementsIndex := -1
	staticIndex := -1
	for i, p := range fields {
		if p == "class" {
			classIndex = i
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
	if scope != nil && scope.IsClass() {
		//fmt.Println("class", name, scope.GetName())	
		if staticIndex == -1 {
			pathSplt[len(pathSplt) - 1] = name 
			className = strings.Join(pathSplt, ".")
		}  else {
			className = scope.GetName() + "." + name
		}
		//fmt.Println(className)
	} else if scope != nil && scope.IsInterface() {
		//fmt.Println("interface", name, scope.GetName())
		pathSplt = append(pathSplt, name)
		className = scope.GetName() + "." + name
		//fmt.Println(className)
		//fmt.Println(scope.GetName(), name)
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

	return Class{
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

func (c*Class) IsInPackage(packSearched string) bool {
	pack, err := c.imports.GetPackage()
	if err != nil {
		return false
	}

	return pack == packSearched
}

func (c*Class) GetPackage() (string, error) {
	return c.imports.GetPackage()
}

func (c*Class) GetMethods() []Method {
	return c.Methods
}

func (c * Class) AppendMethod(m Method) {
	c.Methods = append(c.Methods, m)
}

func (c * Class) GetPath() string {
	return c.path
}

func (c * Class) GetFullPath() string {
	return c.fullPath
}

func (c* Class) GetDocLinesCount() int {
	if len(c.Doc) == 0 {
		return 0
	} 
	return len(strings.Split(c.Doc, "\n"))
}

func (c* Class) GetDoc() string {
	return c.Doc
}

func (c* Class) GetVis() string {
	return c.visibility
}

func (c* Class) GetName() string {
	return c.Name
}

func (c* Class) SetSuper(v string) {
	c.Super = v
}

func (c* Class) GetSuper() string {
	return c.Super
}

func (c* Class) GetInterfaces() []string {
	return c.Interfaces
}



type Interface struct {
	path string
	doc string
	visibility string
	name string
	super string
	methods []Method
}

func (c *Interface)IsClass() bool { return false}
func (c *Interface)IsInterface() bool {return true}

func NewInterface(signature string, doc string, path string, imports *Imports) Interface {

	fields := strings.Fields(RemoveTemplate(signature))

	classIndex := -1 
	extendIndex := -1 
	implementsIndex := -1
	for i, p := range fields {
		if p == "class" {
			classIndex = i
		} else if p == "extends" {
			extendIndex = i
		} else if p == "implements" {
			implementsIndex = i
		}
	}

	var vis string 

	if classIndex < 1 {
		vis = ""
	} else {
		vis = strings.Join(fields[:classIndex], " ")
	}

	className :=strings.Join(strings.Split(strings.Split(path, ".java")[0], "/"), ".")

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

	return Interface{ 
		path: path,
		doc: doc, 
		visibility: vis,
		name: className,
		super: super,
		methods: make([]Method, 0),
	} 
}

func (c * Interface) GetPath() string {
	return c.path
}

func (c* Interface) GetDocLinesCount() int {
	if len(c.doc) == 0 {
		return 0
	} 
	return len(strings.Split(c.doc, "\n"))
}

func (c* Interface) GetDoc() string {
	return c.doc
}

func (c* Interface) GetVis() string {
	return c.visibility
}

func (c* Interface) GetName() string {
	return c.name
}

func (c* Interface) GetSuper() string {
	return c.super
}

func (c * Interface) AppendMethod(m Method) {
	c.methods = append(c.methods, m)
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


