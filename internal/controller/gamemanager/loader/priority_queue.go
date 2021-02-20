package loader

// PriorityQueue is a specific priority queue for grid engine. It has fixed types
// d = dist, is ordered descending. Smallest d on top.
// this implementation uses usual min-heap
type PriorityQueue struct {
	queue []struct {
		d int
		p Pos
	}
	size int
}

// NewPriorityQueue returns a new empty priority queue
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{size: 0}
}

// Top returns the top data
func (pq *PriorityQueue) Top() (int, Pos) {
	if pq.size == 0 {
		panic("Top from empty priority queue")
	}
	return pq.queue[0].d, pq.queue[0].p
}

// Push adds a data
func (pq *PriorityQueue) Push(d int, p Pos) {
	pq.queue = append(pq.queue, struct {
		d int
		p Pos
	}{d, p})
	pq.size++
	// fix heap from last element (the new element)
	now := pq.size - 1
	for now > 0 {
		parent := (now - 1) / 2
		// if now is smaller than parent, swap
		if pq.queue[now].d < pq.queue[parent].d {
			pq.queue[now], pq.queue[parent] = pq.queue[parent], pq.queue[now]
		} else { // no swap = heap fixed = no need to continue
			break
		}
		// it is guaranteed that the parent is always smaller than the sibling
		now = parent
	}
}

// Pop is... well... pop!
func (pq *PriorityQueue) Pop() {
	if pq.size == 0 {
		panic("Pop from empty priority queue")
	}
	pq.queue[0] = pq.queue[pq.size-1]
	pq.queue = pq.queue[:pq.size-1]
	pq.size--
	// fix heap from top
	now := 0
	for now*2+1 < pq.size { // while left child still exists
		// find smaller child between left or right
		smallChild := now*2 + 1
		if smallChild+1 < pq.size && pq.queue[smallChild+1].d < pq.queue[smallChild].d {
			smallChild++
		}
		// if child is smaller than now, swap
		if pq.queue[smallChild].d < pq.queue[now].d {
			pq.queue[smallChild], pq.queue[now] = pq.queue[now], pq.queue[smallChild]
		} else { // no swap = heap fixed = no need to continue
			break
		}
		now = smallChild
	}
}

// Empty checks if the priority queue is empty
func (pq *PriorityQueue) Empty() bool {
	return pq.size == 0
}
