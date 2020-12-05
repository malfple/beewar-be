package loader

import "gitlab.com/otqee/otqee-be/internal/access/formatter/objects"

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
	Width    int
	Height   int
	Terrain  *[][]int
	Units    *[][]objects.Unit
	dist     [][]int // distance matrix, used temporarily
	posQueue []pos   // the queue used for bfs
}

// NewGridEngine returns a new grid engine
func NewGridEngine(width, height int, terrain *[][]int, units *[][]objects.Unit) *GridEngine {
	engine := &GridEngine{
		Width:   width,
		Height:  height,
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

// BFS does a breadth first search starting on (y, x) and fills dist array up to the required steps.
func (ge *GridEngine) BFS(y, x, steps int) {
	if y < 0 || y >= ge.Height || x < 0 || x > ge.Width {
		return
	}
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
			if ty < 0 || ty >= ge.Height || tx < 0 || tx > ge.Width {
				continue
			}
			if (*ge.Units)[ty][tx] != nil {
				continue
			}
			if (*ge.Terrain)[ty][tx] != 1 {
				continue
			}
			if ge.dist[ty][tx] != -1 {
				continue
			}
			ge.dist[ty][tx] = ge.dist[now.Y][now.X] + 1
			ge.posQueue = append(ge.posQueue, pos{ty, tx})
		}
	}
}

// BFSReset is similar to BFS but clears the dist array instead of filling it
func (ge *GridEngine) BFSReset(y, x int) {
	if y < 0 || y >= ge.Height || x < 0 || x > ge.Width {
		return
	}
	ge.dist[y][x] = -1
	ge.posQueue = append(ge.posQueue, pos{y, x})
	for len(ge.posQueue) > 0 {
		now := ge.posQueue[0]
		ge.posQueue = ge.posQueue[1:]

		cy, cx := getAdjList(now.Y, now.X)
		for k := 0; k < K; k++ {
			ty := now.Y + cy[k]
			tx := now.X + cx[k]
			if ty < 0 || ty >= ge.Height || tx < 0 || tx > ge.Width {
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

// ValidateMove checks if a move from (y1, x1) to (y2, x2) with the required steps is valid
func (ge *GridEngine) ValidateMove(y1, x1, y2, x2, steps int) bool {
	if y1 < 0 || y1 >= ge.Height || x1 < 0 || x1 >= ge.Width {
		return false
	}
	if (*ge.Units)[y1][x1] == nil {
		return false
	}
	if y2 < 0 || y2 >= ge.Height || x2 < 0 || x2 >= ge.Width {
		return false
	}
	if (*ge.Units)[y2][x2] != nil {
		return false
	}
	ge.BFS(y1, x1, steps)
	reach := ge.dist[y2][x2] != -1
	ge.BFSReset(y1, x1)
	return reach
}
