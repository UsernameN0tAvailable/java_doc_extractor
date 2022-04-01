package main

import (
	"strings"
)


type Class struct {
	doc string
	visibility string
	name string
	super string
	interfaces []string
	methods []Method
}

func NewClass(signature string, doc string) Class {

	fields := strings.Fields(signature)

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

	className := fields[classIndex + 1]

	var super string

	if extendIndex < 1 {
		super = ""
	} else {
		super = fields[extendIndex + 1]
	}


	implements := make([]string, 0)

	if implementsIndex > 0 {
		tmp := strings.Join(fields[implementsIndex + 1:], " ")
		interfacesStr := strings.Split(tmp, "{")[0]

		for _,in := range strings.Split(interfacesStr, ",") {
			implements = append(implements, strings.TrimSpace(in))
		}
	}

	return Class{
		doc: doc, 
		visibility: vis,
		name: className,
		super: super,
		interfaces: implements,
		methods: make([]Method, 0),
	} 
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

