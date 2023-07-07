package hw04lrucache

import "sync"

type Key string

// We will store this record in the queue - it allows to remove the item from dictionary on queue overflow
// Decided not to store this in the actual map as it is not needed
type record struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

// Set puts a value for the key into the cache and moves it to the front of the queue
// reduces the size of the cache if it is over the limit by removing item from the bottom of the queue (the least accessed one)
func (c *lruCache) Set(key Key, value interface{}) bool {
	// create record to store in the queue
	r := record{key, value}
	// is element key in the cache
	c.mu.Lock()
	if el, ok := c.items[key]; ok {
		// if yes - update the value and move to queue start
		el.Value = r
		c.queue.MoveToFront(el)
		c.mu.Unlock()
		return true
	}
	// if element not in the cache - add the key to dictionary and add to start of queue
	el := c.queue.PushFront(r)
	c.items[key] = el

	//    in case cache size is greater then capacity - remove the last element and its key from the dictionary
	if c.capacity > 0 && c.queue.Len() > c.capacity {
		el := c.queue.Back()
		// remove the key from the dictionary
		delete(c.items, el.Value.(record).key)
		c.queue.Remove(el)
	}
	c.mu.Unlock()
	// return if the value was in cache
	return false
}

// Get returns the value for a given key, and moves the element with this key to the front of the queue
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	// if the key in dictionary then move this element to queue start and return its value and true
	if el, ok := c.items[key]; ok {
		c.queue.MoveToFront(el)
		c.mu.Unlock()
		return el.Value.(record).value, true
	}
	c.mu.Unlock()
	// if the key is not in dictionary return nil and false
	return nil, false
}

// Clear removes all elements from the cache by creating new pointers and counting on GC
func (c *lruCache) Clear() {
	// We count on the fact that Go cleares data not referenced any more so we simply create new lists for cache to work
	// without the actual clearing.
	c.mu.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mu.Unlock()
}

// NewCache creates a new cache and returns it
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
