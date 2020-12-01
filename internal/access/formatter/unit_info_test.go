package formatter

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/otqee/otqee-be/internal/gamemanager/loader/objects/units"
	"testing"
)

// testWidth and testHeight declared in terrain_info_test.go

var testUnitInfo = []byte{
	0, 0, 1, 1, 10, 0,
	2, 1, 1, 3, 8, 1,
}
var testUnitYou = &units.You{
	Owner: 1,
	HP:    10,
	State: 0,
}
var testUnitInfantry = &units.Infantry{
	Owner: 1,
	HP:    8,
	State: 1,
}
var testUnits = [][]units.Unit{
	{testUnitYou, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, testUnitInfantry, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

func TestValidateUnitInfo(t *testing.T) {
	err := ValidateUnitInfo(testWidth, testHeight, testUnitInfo)
	assert.Equal(t, nil, err)
}

func TestConvertUnit(t *testing.T) {
	realUnits := ModelToGameUnit(testWidth, testHeight, testUnitInfo)
	assert.Equal(t, testUnits, realUnits)
	realUnitInfo := GameUnitToModel(testWidth, testHeight, realUnits)
	assert.Equal(t, testUnitInfo, realUnitInfo)
}
