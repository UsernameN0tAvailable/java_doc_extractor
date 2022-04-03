package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// template tests

func TestExtractMethods(t *testing.T) {
	fmt.Println("yooo")
	str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "public class QueryToFilterAdapter", RemoveTemplate(str))
}

func TestNastyClasses(t *testing.T) {
	str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "public class QueryToFilterAdapter", RemoveTemplate(str))
}
