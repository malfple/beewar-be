package formatter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testWidth = 10
const testHeight = 10

var testTerrainInfo = []byte{
	1, 2, 3, 4, 5, 0, 0, 0, 0, 0,
	1, 2, 1, 2, 1, 0, 0, 0, 0, 0,
	5, 4, 3, 2, 1, 0, 0, 0, 0, 0,
	5, 5, 5, 5, 5, 0, 0, 0, 0, 0,
	3, 2, 3, 3, 2, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}
var testTerrains = [][]uint8{
	{1, 2, 3, 4, 5, 0, 0, 0, 0, 0},
	{1, 2, 1, 2, 1, 0, 0, 0, 0, 0},
	{5, 4, 3, 2, 1, 0, 0, 0, 0, 0},
	{5, 5, 5, 5, 5, 0, 0, 0, 0, 0},
	{3, 2, 3, 3, 2, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func TestValidateTerrainInfo(t *testing.T) {
	err := ValidateTerrainInfo(testWidth, testHeight, testTerrainInfo)
	assert.Equal(t, nil, err)
}

func TestConvertTerrain(t *testing.T) {
	realTerrains := ModelToGameTerrain(testWidth, testHeight, testTerrainInfo)
	assert.Equal(t, testTerrains, realTerrains)
	realTerrainInfo := GameTerrainToModel(testWidth, testHeight, realTerrains)
	assert.Equal(t, testTerrainInfo, realTerrainInfo)
}
