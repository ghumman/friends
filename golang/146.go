package golang
type Node struct {
    key, value int
    prev, next *Node
}

type LRUCache struct {
    capacity int
    cache map[int]*Node
    head *Node
    tail *Node
}


func Constructor(capacity int) LRUCache {
    head := &Node{}
    tail := &Node{}

    head.next = tail
    tail.prev = head

    return LRUCache{
        capacity: capacity,
        cache: make(map[int]*Node),
        head: head,
        tail: tail,
    }
}


func (lru *LRUCache) Get(key int) int {
    if node, ok := lru.cache[key]; ok {
        lru.moveToFront(node)
        return node.value
    }
    return -1
}


func (lru *LRUCache) Put(key int, value int)  {
    if node, ok := lru.cache[key]; ok {
        node.value = value
        lru.moveToFront(node)
        return
    }

    node := &Node{key: key, value: value}
    lru.cache[key] = node
    lru.addToFront(node)

    if len(lru.cache) > lru.capacity {
        lru.removeLRU()
    }
}

func (lru *LRUCache) addToFront(node *Node) {
    node.next = lru.head.next
    node.prev = lru.head

    lru.head.next.prev = node
    lru.head.next = node
}

func (lru *LRUCache) remove(node *Node) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (lru *LRUCache) moveToFront(node *Node) {
    lru.remove(node)
    lru.addToFront(node)
}

func(lru * LRUCache) removeLRU() {
    lruLRU := lru.tail.prev
    lru.remove(lruLRU)
    delete(lru.cache, lruLRU.key)
}


/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */