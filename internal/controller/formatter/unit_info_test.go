package formatter

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"testing"
)

// testWidth and testHeight declared in terrain_info_test.go

var testUnitInfo = []byte{
	0, 0, 1, 1, 10, 0,
	2, 1, 1, 3, 8, 1,
}
var testUnitInfo2 = []byte{
	2, 1, 1, 1, 10, 0,
	2, 1, 1, 3, 8, 1,
}
var testUnitQueen = &objects.Queen{
	Owner: 1,
	HP:    10,
	State: 0,
}
var testUnitInfantry = &objects.Infantry{
	Owner: 1,
	HP:    8,
	State: 1,
}
var testUnits = [][]objects.Unit{
	{testUnitQueen, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, testUnitInfantry, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

func TestValidateUnitInfo(t *testing.T) {
	err := ValidateUnitInfo(testHeight, testWidth, testUnitInfo)
	assert.Equal(t, nil, err)
	err = ValidateUnitInfo(testHeight, testWidth, testUnitInfo2)
	assert.Equal(t, ErrMapUnitSamePosition, err)
}

func TestConvertUnit(t *testing.T) {
	realUnits := ModelToGameUnit(testHeight, testWidth, testUnitInfo)
	assert.Equal(t, testUnits, realUnits)
	realUnitInfo := GameUnitToModel(testHeight, testWidth, realUnits)
	assert.Equal(t, testUnitInfo, realUnitInfo)
}
