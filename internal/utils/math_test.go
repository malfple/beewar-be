package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbsInt(t *testing.T) {
	assert.Equal(t, 69, AbsInt(69))
	assert.Equal(t, 69, AbsInt(-69))
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 96, MaxInt(69, 96))
	assert.Equal(t, 96, MaxInt(96, 69))
}

func TestCeilDivInt(t *testing.T) {
	assert.Equal(t, 5, CeilDivInt(10, 2))
	assert.Equal(t, 5, CeilDivInt(9, 2))
	assert.Equal(t, 4, CeilDivInt(8, 2))
	assert.Equal(t, 3, CeilDivInt(9, 3))
	assert.Equal(t, 3, CeilDivInt(7, 3))
	assert.Equal(t, 2, CeilDivInt(6, 3))
	assert.Equal(t, 1, CeilDivInt(1, 69))
}
