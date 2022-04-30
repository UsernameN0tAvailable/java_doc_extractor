package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// template tests

func TestSmth(t *testing.T) {

	in := "public abstract class CacheLoader<K, V>"

	out := RemoveTemplate(in)

	assert.Equal(t, "public abstract class CacheLoader", out)
}

