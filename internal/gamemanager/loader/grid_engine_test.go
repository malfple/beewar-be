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
	{1, -1, 2, -1, -1, -1, -1, 8, -1, -1},
	{2, -1, 3, -1, -1, -1, -1, 9, -1, -1},
	{3, -1, 4, 5, -1, 8, 9, 10, -1, -1},
	{-1, -1, 5, 5, 6, 7, 8, 9, 10, -1},
	{-1, 6, 6, 6, 7, 8, 9, 10, -1, -1},
	{-1, -1, 7, 7, 7, 8, 9, 10, -1, -1},
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
	gl := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	gl.BFS(3, 1, 100)
	assert.Equal(t, expectedDist, gl.dist)

	gl.BFSReset(3, 1)
	assert.Equal(t, cleanDist, gl.dist)
}

func BenchmarkGridEngine_BFS(b *testing.B) {
	gl := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	for i := 0; i < b.N; i++ {
		gl.BFS(3, 1, 100)
		gl.BFSReset(3, 1)
	}
}

func TestGridEngine_ValidateMove(t *testing.T) {
	gl := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	assert.Equal(t, true, gl.ValidateMove(3, 1, 4, 7, 8))
	assert.Equal(t, cleanDist, gl.dist)
	assert.Equal(t, false, gl.ValidateMove(3, 1, 4, 7, 7))
	assert.Equal(t, cleanDist, gl.dist)
}
