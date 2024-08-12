package hw04lrucache

import "sync"

type Key string

type cacheItem struct {
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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	newCacheItem := cacheItem{key: key, value: value}
	if item, ok := c.items[key]; ok {
		item.Value = newCacheItem
		c.queue.MoveToFront(item)
		return true
	}

	if c.queue.Len() == c.capacity {
		oldestItem := c.queue.Back()
		c.queue.Remove(oldestItem)
		delete(c.items, oldestItem.Value.(cacheItem).key)
	}

	newItem := c.queue.PushFront(newCacheItem)
	c.items[key] = newItem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
