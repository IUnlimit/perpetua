package queue

import "sync"

// RingQueue ring buffer queue
type RingQueue struct {
	buffer []interface{}
	size   int
	head   int
	tail   int
	count  int
	lock   sync.Mutex
	cond   *sync.Cond
}

func NewRingQueue(size int) *RingQueue {
	buffer := make([]interface{}, size)
	return &RingQueue{
		buffer: buffer,
		size:   size,
		head:   0,
		tail:   0,
		count:  0,
		cond:   sync.NewCond(&sync.Mutex{}),
	}
}

func (q *RingQueue) Enqueue(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.count == q.size {
		q.cond.Wait()
	}

	q.buffer[q.tail] = item
	q.tail = (q.tail + 1) % q.size
	q.count++

	q.cond.Signal()
}

func (q *RingQueue) Dequeue() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	for q.count == 0 {
		q.cond.Wait()
	}

	item := q.buffer[q.head]
	q.head = (q.head + 1) % q.size
	q.count--

	q.cond.Signal()

	return item
}
