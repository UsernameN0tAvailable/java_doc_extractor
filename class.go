package main

import (
	"strings"
)
type Class struct {
	path string
	doc string
	visibility string
	name string
	super string
	methods []Method
	interfaces []string
}

func NewClass(signature string, doc string, path string, imports *Imports, isInnerClass bool) Class {

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

	name := fields[classIndex + 1]

	var vis string 

	if classIndex < 1 {
		vis = ""
	} else {
		vis = strings.Join(fields[:classIndex], " ")
	}

	pathSplt := strings.Split(strings.Split(path, ".java")[0], "/")

	if isInnerClass {
		pathSplt[len(pathSplt) - 1] = name
	}


	className := strings.Join(pathSplt, ".")

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
		path: path,
		doc: doc, 
		visibility: vis,
		name: strings.TrimSpace(className),
		super: super,
		interfaces: implements,
		methods: make([]Method, 0, 20),
	} 
}

func (c*Class) GetMethods() []Method {
	return c.methods
}

func (c * Class) AppendMethod(m Method) {
	c.methods = append(c.methods, m)
}

func (c * Class) GetPath() string {
	return c.path
}

func (c* Class) GetDocLinesCount() int {
	if len(c.doc) == 0 {
		return 0
	} 
	return len(strings.Split(c.doc, "\n"))
}

func (c* Class) GetDoc() string {
	return c.doc
}

func (c* Class) GetVis() string {
	return c.visibility
}

func (c* Class) GetName() string {
	return c.name
}

func (c* Class) GetSuper() string {
	return c.super
}

func (c* Class) GetInterfaces() []string {
	return c.interfaces
}



type Interface struct {
	path string
	doc string
	visibility string
	name string
	super string
	methods []Method
}

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


