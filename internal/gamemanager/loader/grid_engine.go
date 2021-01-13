package loader

import (
	"gitlab.com/beewar/beewar-be/internal/access/formatter/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

/*
Terrain information in formatter/terrain_info.go

Normal move
BFS. Unit pass-through determined in objects/unit.go


The move pattern for each unit is defined in ValidateMove function
*/

// defines position
type pos struct {
	Y int
	X int
}

// K is the number of adjacent cells (hex cells)
const K = 6

var adjY = []int{0, 0, -1, 1, -1, 1}
var adjXEven = []int{-1, 1, 0, 0, 1, 1}
var adjXOdd = []int{-1, 1, 0, 0, -1, -1}

func getAdjList(y, x int) ([]int, []int) {
	if y&1 == 0 { // even row
		return adjY, adjXEven
	}
	return adjY, adjXOdd
}

// GridEngine is a game engine for movement and attack range calculations
// keep in mind that GridEngine should not edit the pointer to slices (treated as input only)
type GridEngine struct {
	Height   int
	Width    int
	Terrain  *[][]int
	Units    *[][]objects.Unit
	dist     [][]int // distance matrix, used temporarily
	posQueue []pos   // the queue used for bfs
}

// NewGridEngine returns a new grid engine
func NewGridEngine(height, width int, terrain *[][]int, units *[][]objects.Unit) *GridEngine {
	engine := &GridEngine{
		Height:  height,
		Width:   width,
		Terrain: terrain,
		Units:   units,
		dist:    make([][]int, height),
	}
	for i := 0; i < height; i++ {
		engine.dist[i] = make([]int, width)
		for j := 0; j < width; j++ {
			engine.dist[i][j] = -1
		}
	}

	return engine
}

// insideMap checks if position is inside map
func (ge *GridEngine) insideMap(y, x int) bool {
	return y >= 0 && y < ge.Height && x >= 0 && x < ge.Width
}

// BFS does a breadth first search starting on (y, x) and fills dist array up to the required steps.
// This is also constraint by the given owner and weight
// WARNING: this function does not do validation checks
func (ge *GridEngine) BFS(y, x, steps, owner, weight int) {
	ge.dist[y][x] = 0
	ge.posQueue = append(ge.posQueue, pos{y, x})
	for len(ge.posQueue) > 0 {
		now := ge.posQueue[0]
		ge.posQueue = ge.posQueue[1:]

		if ge.dist[now.Y][now.X] >= steps {
			continue
		}

		cy, cx := getAdjList(now.Y, now.X)
		for k := 0; k < K; k++ {
			ty := now.Y + cy[k]
			tx := now.X + cx[k]
			if !ge.insideMap(ty, tx) {
				continue
			}
			if ge.dist[ty][tx] != -1 {
				continue
			}
			if (*ge.Terrain)[ty][tx] != 1 {
				continue
			}
			if unit := (*ge.Units)[ty][tx]; unit != nil {
				if unit.GetUnitOwner() != owner {
					continue
				}
				if unit.GetWeight()+weight > 1 {
					continue
				}
			}
			ge.dist[ty][tx] = ge.dist[now.Y][now.X] + 1
			ge.posQueue = append(ge.posQueue, pos{ty, tx})
		}
	}
}

// BFSReset is similar to BFS but clears the dist array instead of filling it
// it has to be used at the same spot when doing BFS
// WARNING: this function does not do validation checks
func (ge *GridEngine) BFSReset(y, x int) {
	ge.dist[y][x] = -1
	ge.posQueue = append(ge.posQueue, pos{y, x})
	for len(ge.posQueue) > 0 {
		now := ge.posQueue[0]
		ge.posQueue = ge.posQueue[1:]

		cy, cx := getAdjList(now.Y, now.X)
		for k := 0; k < K; k++ {
			ty := now.Y + cy[k]
			tx := now.X + cx[k]
			if !ge.insideMap(ty, tx) {
				continue
			}
			if ge.dist[ty][tx] == -1 {
				continue
			}
			ge.dist[ty][tx] = -1
			ge.posQueue = append(ge.posQueue, pos{ty, tx})
		}
	}
}

// ValidateMoveNormal checks if a normal move from (y1, x1) to (y2, x2) with the required steps is valid
// WARNING: does not validate positions or if a unit exists. Only validates move with BFS
func (ge *GridEngine) ValidateMoveNormal(y1, x1, y2, x2, steps int) bool {
	var reach bool
	self := (*ge.Units)[y1][x1]
	ge.BFS(y1, x1, steps, self.GetUnitOwner(), self.GetWeight())
	reach = ge.dist[y2][x2] != -1
	ge.BFSReset(y1, x1)
	return reach
}

// ValidateMove checks if a unit move from (y1, x1) to (y2, x2) is valid
func (ge *GridEngine) ValidateMove(y1, x1, y2, x2 int) bool {
	if !ge.insideMap(y1, x1) {
		return false
	}
	if (*ge.Units)[y1][x1] == nil {
		return false
	}
	if !ge.insideMap(y2, x2) {
		return false
	}
	if (*ge.Units)[y2][x2] != nil {
		return false
	}

	switch unit := (*ge.Units)[y1][x1]; unit.GetMoveType() {
	case objects.MoveTypeGround:
		return ge.ValidateMoveNormal(y1, x1, y2, x2, unit.GetMoveRange())
	default:
		panic("panic validate move: unknown move type")
	}
	return false
}

// ValidateAttack checks if an attack from (y, x) to (yt, xt) is valid
// there has to be a unit at (yt, xt), but not necessarily at (y, x). Therefore the attacker unit is needed
func (ge *GridEngine) ValidateAttack(y, x, yt, xt int, attacker objects.Unit) (bool, int) {
	if !ge.insideMap(y, x) {
		return false, -1
	}
	if !ge.insideMap(yt, xt) {
		return false, -1
	}
	if (*ge.Units)[yt][xt] == nil {
		return false, -1
	}
	if attacker.GetUnitOwner() == (*ge.Units)[yt][xt].GetUnitOwner() {
		return false, -1
	}

	distBetween := utils.HexDistance(y, x, yt, xt)
	switch attacker.GetUnitType() {
	case objects.UnitTypeInfantry:
		return distBetween <= objects.UnitAttackRangeInfantry, distBetween
	default:
		panic("panic validate move: unknown unit type")
	}
	return false, -1
}
