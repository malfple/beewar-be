package loader

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/otqee/otqee-be/internal/access/formatter"
	"testing"
)

const testHeight = 10
const testWidth = 10

var testTerrain = formatter.ModelToGameTerrain(testHeight, testWidth, []byte{
	1, 0, 1, 1, 1, 1, 1, 1, 0, 0,
	1, 0, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 1, 1, 1, 1, 0, 1, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 0, 1, 1, 1,
	1, 1, 1, 1, 0, 1, 1, 1, 1, 0,
	0, 1, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 1, 1, 1, 1, 1, 1, 1, 0, 1,
	0, 0, 1, 1, 1, 1, 1, 1, 0, 1,
})

var testUnits = formatter.ModelToGameUnit(testHeight, testWidth, []byte{
	5, 1, 1, 1, 10, 0,
	4, 1, 1, 3, 10, 0,
	3, 1, 1, 3, 10, 0,
	6, 1, 1, 3, 10, 0,
	7, 1, 1, 3, 10, 0,
	4, 8, 2, 1, 10, 0,
	3, 8, 2, 3, 10, 0,
	2, 8, 2, 3, 10, 0,
	5, 8, 2, 3, 10, 0,
	6, 8, 2, 3, 10, 0,
})

var expectedDist = [][]int{
	{-1, -1, 3, 4, 5, 6, 7, 8, -1, -1},
	{-1, -1, 2, 3, 4, 5, 6, 7, 8, -1},
	{-1, 1, 2, 3, 4, 5, 6, 7, -1, -1},
	{-1, 0, 1, 2, 3, -1, 6, 7, -1, -1},
	{1, 1, 2, -1, -1, -1, -1, 8, -1, -1},
	{2, 2, 2, -1, -1, -1, -1, 9, -1, -1},
	{3, 3, 3, 4, -1, 7, 8, 9, -1, -1},
	{-1, 4, 4, 4, 5, 6, 7, 8, 9, -1},
	{-1, 5, 5, 5, 6, 7, 8, 9, -1, -1},
	{-1, -1, 6, 6, 6, 7, 8, 9, -1, -1},
}

var cleanDist = [][]int{
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
}

func TestGridEngine_BFS(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	ge.BFS(3, 1, 100)
	assert.Equal(t, expectedDist, ge.dist)

	ge.BFSReset(3, 1)
	assert.Equal(t, cleanDist, ge.dist)
}

func BenchmarkGridEngine_BFS(b *testing.B) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	for i := 0; i < b.N; i++ {
		ge.BFS(3, 1, 100)
		ge.BFSReset(3, 1)
	}
}

func TestGridEngine_ValidateMoveNormal(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	assert.Equal(t, true, ge.ValidateMoveNormal(3, 1, 4, 7, 8))
	assert.Equal(t, cleanDist, ge.dist)
	assert.Equal(t, false, ge.ValidateMoveNormal(3, 1, 4, 7, 7))
	assert.Equal(t, cleanDist, ge.dist)
}
