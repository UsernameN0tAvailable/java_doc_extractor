package main

import (
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

func TestExtractMethodsAndInnerClasses(t *testing.T) {
	path := "./NastyInnerClasses.java"
	
	extractor := Extractor{classes: make([]Class, 0, 20000), interfaces: make([]Interface, 0, 10000), activeClasses: make([]*Class, 0, 200), activeClass: nil}

	extractor.Extract(path)

	//str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "test", "test")
}
