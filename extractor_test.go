package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasty(t *testing.T) {


	extractor := NewExtractor()

	classes := extractor.Extract("java_test_data/CharMatcher.java")

	for _,c := range classes {
		fmt.Println("name:", c.GetName(), "super:", c.GetSuper())
	}

	assert.Equal(t, true, true) 
}
