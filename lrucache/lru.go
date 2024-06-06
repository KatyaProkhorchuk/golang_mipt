//go:build !solution

package lrucache

import (
	"container/list"
)

type lruCache struct {
	lruList  *list.List // элементы кэша в понядке устаревания
	mp       map[int]*list.Element
	values   map[int]int
	capacity int
}

func (lru *lruCache) Update(key int) bool {
	if lru.capacity == 0 {
		return false
	}
	if node, ok := lru.mp[key]; ok {
		lru.lruList.MoveToFront(node)
		return true
	}
	if len(lru.mp) >= lru.capacity {
		deleteNode := lru.lruList.Back()
		deleteKey := deleteNode.Value.(int)
		delete(lru.mp, deleteKey)
		delete(lru.values, deleteKey)
		lru.lruList.Remove(deleteNode)
	}
	lru.mp[key] = lru.lruList.PushFront(key)
	return true
}

func (lru *lruCache) Set(key, value int) {
	if lru.Update(key) {
		lru.values[key] = value
	}
}

func (lru *lruCache) Clear() {
	lru.mp = make(map[int]*list.Element, lru.capacity)
	lru.values = make(map[int]int, lru.capacity)
	lru.lruList = list.New()
}

func (lru *lruCache) Range(f func(key, value int) bool) {
	for node := lru.lruList.Back(); node != nil; node = node.Prev() {
		key := node.Value.(int)
		if !f(key, lru.values[key]) {
			break
		}
	}
}

func (lru *lruCache) Get(key int) (int, bool) {
	val, ok := lru.values[key]
	ok = ok && lru.Update(key) // вконец
	return val, ok
}

func New(cap int) Cache {
	return &lruCache{
		capacity: cap,
		lruList:  list.New(),
		mp:       make(map[int]*list.Element, cap),
		values:   make(map[int]int, cap),
	}
}
