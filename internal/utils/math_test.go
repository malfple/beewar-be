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
