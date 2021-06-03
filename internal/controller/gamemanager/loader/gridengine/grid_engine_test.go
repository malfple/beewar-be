package gridengine

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/beewar/beewar-be/internal/controller/formatter"
	"gitlab.com/beewar/beewar-be/internal/utils"
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

func TestGridEngine_FillMoveGround(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	self := (*ge.Units)[3][1]
	ge.FillMoveGround(3, 1, 100, self.GetUnitOwner(), self.GetWeight())
	assert.Equal(t, expectedDist, ge.Dist)

	ge.FillMoveGroundReset(3, 1)
	assert.Equal(t, cleanDist, ge.Dist)
}

func BenchmarkGridEngine_FillMoveGround(b *testing.B) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	self := (*ge.Units)[3][1]
	for i := 0; i < b.N; i++ {
		ge.FillMoveGround(3, 1, 100, self.GetUnitOwner(), self.GetWeight())
		ge.FillMoveGroundReset(3, 1)
	}
}

func TestGridEngine_ValidateMoveGround(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	assert.Equal(t, true, ge.ValidateMoveGround(3, 1, 4, 7, 8))
	assert.Equal(t, cleanDist, ge.Dist)
	assert.Equal(t, false, ge.ValidateMoveGround(3, 1, 4, 7, 7))
	assert.Equal(t, cleanDist, ge.Dist)
}

func TestGridEngine_ValidateMove(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain, &testUnits)

	assert.Equal(t, true, ge.ValidateMove(3, 1, 3, 4))
	assert.Equal(t, false, ge.ValidateMove(3, 1, 3, 5))
	assert.Equal(t, false, ge.ValidateMove(3, 1, 6, 1))
	assert.Equal(t, true, ge.ValidateMove(3, 1, 6, 2))
}

var testTerrain2 = formatter.ModelToGameTerrain(testHeight, testWidth, []byte{
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
})
var testUnits2 = formatter.ModelToGameUnit(testHeight, testWidth, []byte{})

func TestHexDistance(t *testing.T) {
	ge := NewGridEngine(testHeight, testWidth, &testTerrain2, &testUnits2)

	ge.FillMoveGround(5, 5, 100, 0, 0)

	for i := 0; i < testHeight; i++ {
		for j := 0; j < testWidth; j++ {
			assert.Equal(t, ge.Dist[i][j], utils.HexDistance(5, 5, i, j))
		}
	}
}
