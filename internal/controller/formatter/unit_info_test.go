package formatter

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"testing"
)

// testWidth and testHeight declared in terrain_info_test.go

var testUnitInfo = []byte{
	0, 0, 1, 1, 10, 0,
	0, 1, 2, 1, 10, 0,
	2, 1, 1, 3, 8, 1,
	3, 0, 1, 4, 8, 0,
	3, 1, 1, 5, 10, 0,
	3, 2, 1, 6, 14, 1,
	4, 0, 2, 9, 4, 0,
}
var testUnitInfo2 = []byte{
	0, 0, 1, 1, 10, 0,
	0, 0, 2, 1, 10, 0,
}
var testUnitInfo3 = []byte{
	0, 0, 1, 1, 10, 0,
}
var testUnitInfo4 = []byte{
	100, 100, 1, 1, 10, 0,
}
var testUnitInfo5 = []byte{
	0, 0, 3, 1, 10, 0,
}
var testUnitInfo6 = []byte{
	0, 0, 1, 1, 10, 0, 0,
}
var testUnitQueen = objects.NewQueen(1, 10, 0)
var testUnitQueen2 = objects.NewQueen(2, 10, 0)
var testUnitInfantry = objects.NewInfantry(1, 8, 1)
var testUnitJetCrew = objects.NewJetCrew(1, 8, 0)
var testUnitWizard = objects.NewWizard(1, 10, 0)
var testUnitTank = objects.NewTank(1, 14, 1)
var testUnitMortar = objects.NewMortar(2, 4, 0)
var testUnits = [][]objects.Unit{
	{testUnitQueen, testUnitQueen2, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, testUnitInfantry, nil, nil, nil, nil, nil, nil, nil, nil},
	{testUnitJetCrew, testUnitWizard, testUnitTank, nil, nil, nil, nil, nil, nil, nil},
	{testUnitMortar, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

func TestValidateUnitInfo(t *testing.T) {
	err := ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo, false)
	assert.Equal(t, nil, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo2, false)
	assert.Equal(t, ErrMapUnitSamePosition, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo3, false)
	assert.Equal(t, ErrMapPlayerQueen, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo3, true)
	assert.Equal(t, nil, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo4, true)
	assert.Equal(t, ErrMapUnitOutsideMap, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo5, true)
	assert.Equal(t, ErrMapPlayerNotExist, err)
	err = ValidateUnitInfo(testHeight, testWidth, 2, testUnitInfo6, true)
	assert.Equal(t, ErrMapInvalidUnitInfo, err)
}

func TestConvertUnit(t *testing.T) {
	realUnits := ModelToGameUnit(testHeight, testWidth, testUnitInfo)
	assert.Equal(t, testUnits, realUnits)
	realUnitInfo := GameUnitToModel(testHeight, testWidth, realUnits)
	assert.Equal(t, testUnitInfo, realUnitInfo)
}
