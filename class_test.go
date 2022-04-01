package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)



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
}

