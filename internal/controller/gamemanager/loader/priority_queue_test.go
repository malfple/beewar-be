package loader

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	unsorted := []int{8, 3, 5, 4, 6, 1, 2, 7}

	pq := NewPriorityQueue()
	for _, v := range unsorted {
		pq.Push(v, Pos{v, v})
	}
	for i := range unsorted {
		assert.Equal(t, false, pq.Empty())
		d, p := pq.Top()
		pq.Pop()

		assert.Equal(t, i+1, d)
		assert.Equal(t, i+1, p.X)
		assert.Equal(t, i+1, p.Y)
	}
	assert.Equal(t, true, pq.Empty())
}
