package gridengine

import (
	"gitlab.com/beewar/beewar-be/internal/controller/gamemanager/objects"
	"gitlab.com/beewar/beewar-be/internal/utils"
)

/*
Terrain information in formatter/terrain_info.go

Normal move
FillMoveGround. Unit pass-through determined in objects/unit.go


The move pattern for each unit is defined in ValidateMove function
*/

// Pos defines position
type Pos struct {
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
	Height  int
	Width   int
	Terrain *[][]int
	Units   *[][]objects.Unit
	Dist    [][]int        // distance matrix, used temporarily
	pqueue  *PriorityQueue // the priority_queue used for dijkstra
	queue   []Pos          // the queue for bfs (to reset the dijkstra)
}

// NewGridEngine returns a new grid engine
func NewGridEngine(height, width int, terrain *[][]int, units *[][]objects.Unit) *GridEngine {
	engine := &GridEngine{
		Height:  height,
		Width:   width,
		Terrain: terrain,
		Units:   units,
		Dist:    make([][]int, height),
		pqueue:  NewPriorityQueue(),
	}
	for i := 0; i < height; i++ {
		engine.Dist[i] = make([]int, width)
		for j := 0; j < width; j++ {
			engine.Dist[i][j] = -1
		}
	}

	return engine
}

// insideMap checks if position is inside map
func (ge *GridEngine) insideMap(y, x int) bool {
	return y >= 0 && y < ge.Height && x >= 0 && x < ge.Width
}

// FillMoveGround does a Dijkstra starting on (y, x) and fills dist array up to the required steps.
// This is also constraint by the given owner and weight.
// WARNING: this function does not do validation checks.
func (ge *GridEngine) FillMoveGround(y, x, steps, owner, weight int) {
	ge.pqueue.Push(0, Pos{y, x})
	for !ge.pqueue.Empty() {
		d, now := ge.pqueue.Top()
		ge.pqueue.Pop()

		if ge.Dist[now.Y][now.X] != -1 { // already visited
			continue
		}
		ge.Dist[now.Y][now.X] = d
		if d == steps { // end of steps. break to optimize
			continue
		}

		cy, cx := getAdjList(now.Y, now.X)
		for k := 0; k < K; k++ {
			ty := now.Y + cy[k]
			tx := now.X + cx[k]
			if !ge.insideMap(ty, tx) {
				continue
			}
			if ge.Dist[ty][tx] != -1 {
				continue
			}
			if unit := (*ge.Units)[ty][tx]; unit != nil {
				if unit.GetOwner() != owner {
					continue
				}
				if unit.UnitWeight()+weight > 1 {
					continue
				}
			}
			if dnext := d + CalcMoveCost((*ge.Terrain)[ty][tx], weight); dnext <= steps {
				ge.pqueue.Push(dnext, Pos{ty, tx})
			}
		}
	}
}

// FillMoveGroundReset is similar to FillMoveGround but clears the dist array instead of filling it.
// It has to be used at the same spot when doing FillMoveGround.
// It also uses BFS instead of Dijkstra (because we don't need dijkstra to clear).
// WARNING: this function does not do validation checks.
func (ge *GridEngine) FillMoveGroundReset(y, x int) {
	ge.Dist[y][x] = -1
	ge.queue = append(ge.queue, Pos{y, x})
	for len(ge.queue) > 0 {
		now := ge.queue[0]
		ge.queue = ge.queue[1:]

		cy, cx := getAdjList(now.Y, now.X)
		for k := 0; k < K; k++ {
			ty := now.Y + cy[k]
			tx := now.X + cx[k]
			if !ge.insideMap(ty, tx) {
				continue
			}
			if ge.Dist[ty][tx] == -1 {
				continue
			}
			ge.Dist[ty][tx] = -1
			ge.queue = append(ge.queue, Pos{ty, tx})
		}
	}
}

// ValidateMoveGround checks if a normal move from (y1, x1) to (y2, x2) with the required steps is valid
// WARNING: does not validate positions or if a unit exists. Only validates move with FillMoveGround
func (ge *GridEngine) ValidateMoveGround(y1, x1, y2, x2, steps int) bool {
	var reach bool
	self := (*ge.Units)[y1][x1]
	ge.FillMoveGround(y1, x1, steps, self.GetOwner(), self.UnitWeight())
	reach = ge.Dist[y2][x2] != -1
	ge.FillMoveGroundReset(y1, x1)
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

	switch unit := (*ge.Units)[y1][x1]; unit.UnitMoveType() {
	case objects.MoveTypeGround:
		return ge.ValidateMoveGround(y1, x1, y2, x2, unit.UnitMoveRange())
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
	if attacker.GetOwner() == (*ge.Units)[yt][xt].GetOwner() {
		return false, -1
	}

	distBetween := utils.HexDistance(y, x, yt, xt)
	switch attacker.UnitAttackType() {
	case objects.AttackTypeNone:
		return false, -1
	case objects.AttackTypeGround:
		return distBetween <= attacker.UnitAttackRange(), distBetween
	default:
		panic("panic validate attack: unknown attack type")
	}
	return false, -1
}
