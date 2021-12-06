package ringbuffer

import (
	"log"
	"sync"
)

type RingBuffer struct {
	array []int
	pos   int
	size  int
	m     sync.RWMutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{make([]int, size), -1, size, sync.RWMutex{}}
}

func (r *RingBuffer) Push(el int) {
	log.Println("RingBuffer Push")
	r.m.Lock()
	defer r.m.Unlock()
	if r.pos == r.size-1 {
		for i := 1; i <= r.size-1; i++ {
			r.array[i-1] = r.array[i]
		}
		r.array[r.pos] = el
	} else {
		r.pos++
		r.array[r.pos] = el
	}
}

func (r *RingBuffer) Get() []int {
	log.Println("RingBuffer Get")
	if r.pos < 0 {
		return nil
	}
	r.m.RLock()
	defer r.m.RUnlock()
	var output []int = r.array[:r.pos]
	r.pos = 0
	return output
}
