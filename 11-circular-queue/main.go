package main

import (
	"errors"
	"sync"
)

type CircularQueue struct {
	data []interface{}
	capacity int
	head int
	tail int
	size int
	mu sync.Mutex
}

func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		data: make([]interface{}, capacity),
		capacity: capacity,
		head: 0,
		tail: 0,
		size: 0,
	}
}

func (q *CircularQueue) Enqueue(item interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == q.capacity {
		return errors.New("queue is full")
	}

	q.data[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return nil
}


func (q *CircularQueue) Dequeue() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.size == 0 {
		return nil, errors.New("queue is empty")
	}

	item := q.data[q.head]
	q.data[q.head] = nil 
	q.head = (q.head + 1) % q.capacity
	q.size--

	return item, nil
}