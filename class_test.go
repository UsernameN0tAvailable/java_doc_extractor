package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


// template tests

func TestWithSuperAndMultipleIntefaces(t *testing.T) {


	assert.Equal(t, "DiVater", "DiVater")
}

func TestRemoveTemplate(t *testing.T) {

	str := "BaseGeometryTestCase<GeometryCollection<Geometry>>"

	assert.Equal(t, "BaseGeometryTestCase", RemoveTemplate(str))
}

func TestRemoveMultipleTemplates(t *testing.T) {

	str := "BaseGeometryTestCase<GeometryCollection<Geometry>> extends SomeSuperClass<SomeTemplate extends SomeSuperTemplate> implements SomeInterface<SomeInterfaceTemplate>"
	assert.Equal(t, "BaseGeometryTestCase extends SomeSuperClass implements SomeInterface", RemoveTemplate(str))
}
func TestDoNotAlter(t *testing.T) {

	str := "BaseGeometryTestCase extends SomeSuperClass implements SomeInterface"
	assert.Equal(t, "BaseGeometryTestCase extends SomeSuperClass implements SomeInterface", RemoveTemplate(str))
}


func TestRemoveButDoesNotEndWithTemplate(t *testing.T) {

	str := "BaseGeometryTestCase<GeometryCollection<Geometry>> extends SomeSuperClass<SomeTemplate extends SomeSuperTemplate> implements SomeInterface"
	assert.Equal(t, "BaseGeometryTestCase extends SomeSuperClass implements SomeInterface", RemoveTemplate(str))
}

func TestRemoveTemplateAtEnd(t *testing.T) {

	str := "BaseGeometryTestCase extends SomeSuperClass<SomeTemplate extends SomeSuperTemplate> implements SomeInterface<SomeInterfaceTemplate>"
	assert.Equal(t, "BaseGeometryTestCase extends SomeSuperClass implements SomeInterface", RemoveTemplate(str))
}

	
func TestRemoveTemplateMore(t *testing.T) {

	str := "public class QueryToFilterAdapter<Q extends Query>"
	assert.Equal(t, "public class QueryToFilterAdapter", RemoveTemplate(str))
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

