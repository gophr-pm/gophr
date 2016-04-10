package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 1, abs(-1), "abs should work on negative values")
	assert.Equal(t, 2, abs(2), "abs should work on positive values")
}
