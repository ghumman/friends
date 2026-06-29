package main

import (
	"sync"
)

type Node struct {
	key int
	value interface{}
	next *Node
	prev *Node
}

type LRUCache struct {
	cache map[int]*Node
	capacity int
	head *Node
	tail *Node
	mu sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		panic("capacity must be greater than 0")
	}

	head := &Node{}
	tail := &Node{}
	
	head.next = tail
	tail.prev = head

	return &LRUCache{
		cache: make(map[int]*Node),
		capacity: capacity,
		head: head,
		tail: tail,
	}
}

func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *LRUCache) insertFront(node *Node) {
	node.next = c.head.next
	node.prev = c.head

	c.head.next.prev = node
	c.head.next = node
}

func (c *LRUCache) Get(key int) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.remove(node)
	c.insertFront(node)
	return node.value, true
}

func (c *LRUCache) Set(key int, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		node.value = value
		c.remove(node)
		c.insertFront(node)
		return
	}

	node := &Node{
		key: key,
		value: value,
	}
	c.cache[key] = node
	c.insertFront(node)

	if len(c.cache) > c.capacity {
		lru := c.tail.prev

		c.remove(lru)
		delete(c.cache, lru.key)
	}

}