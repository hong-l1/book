package data

import (
	"fmt"
	"sync"
)

type BoundedQueue struct {
	Queue    []int
	capacity int
	*sync.Cond
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	return &BoundedQueue{
		Queue:    make([]int, 0, capacity),
		capacity: capacity,
		Cond:     sync.NewCond(&sync.Mutex{}),
	}
}
func (q *BoundedQueue) Enqueue(item int) {
	q.L.Lock()
	defer q.L.Unlock()
	for len(q.Queue) == q.capacity {
		q.Wait()
	}
	q.Queue = append(q.Queue, item)
	fmt.Printf("Enqueued: %d\n", item)
	q.Cond.Signal()
}
func (q *BoundedQueue) Dequeue() int {
	q.L.Lock()
	defer q.L.Unlock()
	for len(q.Queue) == 0 {
		q.Wait()
	}
	item := q.Queue[0]
	q.Queue = q.Queue[1:]
	fmt.Printf("Dequeue: %d\n", item)
	q.Cond.Signal()
	return item
}
