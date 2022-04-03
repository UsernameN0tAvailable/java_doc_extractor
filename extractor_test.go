package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// template tests
/*
func TestExtractMethods(t *testing.T) {
	path := "./NastyMethods.java"

	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]*Class, 0, 200), activeClass: nil}

	extractor.Extract(path)

	//str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "test", "test")
}*/
/*
func TestExtractMethodsAndInnerClasses(t *testing.T) {
	path := "./NastyInnerClasses.java"
	
	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]Scope, 0, 200), activeClass: nil}

	extractor.Extract(path)

	//str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "test", "test")
}*/


func TestNasty(t *testing.T) {
	path := "./Action.java"
	
	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]Scope, 0, 200), activeClass: nil}

	classes, interfaces := extractor.Extract(path)


	for _,class := range classes {
		fmt.Println(class.GetName())
	}

	for _,in := range interfaces {
		fmt.Println(in.GetName(), "in")
	}

	//str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, true, true) 
}
