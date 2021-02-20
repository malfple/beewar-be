package loader

// PriorityQueue is a specific priority queue for grid engine. It has fixed types
// d = dist, is ordered descending. Smallest d on top.
// this implementation uses usual min-heap
type PriorityQueue struct {
	heap []struct {
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
	return pq.heap[0].d, pq.heap[0].p
}

// Push adds a data
func (pq *PriorityQueue) Push(d int, p Pos) {
	pq.heap = append(pq.heap, struct {
		d int
		p Pos
	}{d, p})
	pq.size++
	// fix heap from last element (the new element)
	now := pq.size - 1
	for now > 0 {
		parent := (now - 1) / 2
		// if now is smaller than parent, swap
		if pq.heap[now].d < pq.heap[parent].d {
			pq.heap[now], pq.heap[parent] = pq.heap[parent], pq.heap[now]
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
	pq.heap[0] = pq.heap[pq.size-1]
	pq.heap = pq.heap[:pq.size-1]
	pq.size--
	// fix heap from top
	now := 0
	for now*2+1 < pq.size { // while left child still exists
		// find smaller child between left or right
		smallChild := now*2 + 1
		if smallChild+1 < pq.size && pq.heap[smallChild+1].d < pq.heap[smallChild].d {
			smallChild++
		}
		// if child is smaller than now, swap
		if pq.heap[smallChild].d < pq.heap[now].d {
			pq.heap[smallChild], pq.heap[now] = pq.heap[now], pq.heap[smallChild]
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
