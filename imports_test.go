package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


// template tests

func TestImport(t *testing.T) {

	imports := "package org.elasticsearch.xpack.restart;\nimport org.elasticsearch.common.settings.Settings;\nimport org.elasticsearch.common.util.concurrent.ThreadContext; \nimport java.nio.charset.StandardCharsets;\nimport java.util.Base64;"


	imp := NewImports([]byte(imports))

	assert.Equal(t,"org.elasticsearch.common.util.concurrent.ThreadContext.SomeOtherContext", imp.GetPath("common.util.concurrent.ThreadContext.SomeOtherContext"))

	assert.Equal(t,"org.elasticsearch.xpack.restart.util.common.concurrent.ThreadContext.SomeOtherContext", imp.GetPath("util.common.concurrent.ThreadContext.SomeOtherContext"))
}
/*
func TestWithSuperAndMultipleIntefaces(t *testing.T) {

	class := NewClass("public class DiniMuetter extends DiVater implements Parents,Elders {", "")

	actualInt := class.GetInterfaces();
	expectedInt := [2]string{"Parents", "Elders"}

	for i,in := range expectedInt {
		assert.Equal(t, actualInt[i], in)
	}

	assert.Equal(t, "public", class.GetVis())
	assert.Equal(t, "DiniMuetter", class.GetName())
	assert.Equal(t, "DiVater", class.GetSuper())
}

func TestWithSuperAndMultipleIntefacesWithoutClosingBracket(t *testing.T) {

	class := NewClass("public class DiniMuetter extends DiVater implements Parents,Elders ", "")

	actualInt := class.GetInterfaces();
	expectedInt := [2]string{"Parents", "Elders"}

	for i,in := range expectedInt {
		assert.Equal(t, actualInt[i], in)
	}

	assert.Equal(t, "public", class.GetVis())
	assert.Equal(t, "DiniMuetter", class.GetName())
	assert.Equal(t, "DiVater", class.GetSuper())
}

func TestWithSuper(t *testing.T) {

	class := NewClass("public class DiniMuetter extends DiVater", "")

	assert.Equal(t, 0, len(class.GetInterfaces()))

	assert.Equal(t, "public", class.GetVis())
	assert.Equal(t, "DiniMuetter", class.GetName())
	assert.Equal(t, "DiVater", class.GetSuper())
}


func TestWithInterfacesOnly(t *testing.T) {

	class := NewClass("public class DiniMuetter implements Parents,Elders ", "")

	actualInt := class.GetInterfaces();
	expectedInt := [2]string{"Parents", "Elders"}

	for i,in := range expectedInt {
		assert.Equal(t, actualInt[i], in)
	}

	assert.Equal(t, "public", class.GetVis())
	assert.Equal(t, "DiniMuetter", class.GetName())	
} */

