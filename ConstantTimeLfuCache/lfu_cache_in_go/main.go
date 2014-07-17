package main

import (
	"container/list"
	"errors"
	"fmt"
)

type FrequencyNode struct {
	Value    int
	elements map[string]bool
	count    int
}

func NewFrequencyNode(value int) *FrequencyNode {
	elements := make(map[string]bool)
	return &FrequencyNode{value, elements, 0}
}

func (f *FrequencyNode) Len() int {
	return f.count
}

func (f *FrequencyNode) Contains(key string) bool {
	return f.elements[key]
}

func (f *FrequencyNode) Add(key string) {
	f.elements[key] = true
	f.count++
}

func (f *FrequencyNode) Remove(key string) {
	if f.elements[key] {
		f.count--
	}
	f.elements[key] = false
}

type CacheItem struct {
	Data   interface{}
	Parent *list.Element
}

func NewCacheItem(data interface{}, parent *list.Element) *CacheItem {
	return &CacheItem{data, parent}
}

type LfuCache struct {
	dataLookup    map[string]*CacheItem
	frequencyHead *list.Element
	freqList      *list.List
}

func NewLfuCache() LfuCache {
	cache := LfuCache{
		dataLookup: make(map[string]*CacheItem),
		freqList:   list.New(),
	}
	cache.frequencyHead = cache.freqList.PushFront(NewFrequencyNode(0))
	return cache
}

func (l LfuCache) Head() *list.Element {
	return l.frequencyHead
}

func (l LfuCache) FreqListLen() int {
	return l.freqList.Len() - 1 // ignore the first node
}

func (l LfuCache) GetNewNode(value int, parent *list.Element) (*list.Element, error) {
	if parent == nil {
		return nil, errors.New("Cannot insert node because parent is nil")
	}
	node := l.freqList.InsertAfter(NewFrequencyNode(value), parent)
	return node, nil
}

func (l LfuCache) DeleteNode(node *list.Element) {
	l.freqList.Remove(node)
}

func (l LfuCache) FreqNodeFor(value int) *FrequencyNode {
	node := l.frequencyHead
	for i := value; i > 0; i-- {
		if node.Next() == nil {
			break
		}
		node = node.Next()
	}

	return node.Value.(*FrequencyNode)
}

func (l LfuCache) Fetch(key string) (data interface{}, err error) {
	var freq, nextFreq *FrequencyNode
	var node, nextNode *list.Element
	elem := l.dataLookup[key]
	if elem == nil {
		return nil, errors.New(fmt.Sprintf("'%s' does not exist", key))
	}
	node = elem.Parent
	freq = node.Value.(*FrequencyNode)
	if nextNode = node.Next(); nextNode != nil {
		nextFreq = nextNode.Value.(*FrequencyNode)
	}

	if nextNode == nil || nextFreq.Value != freq.Value+1 {
		nextNode, err = l.GetNewNode(freq.Value+1, node)
		if err != nil {
			return nil, err
		}
		nextFreq = nextNode.Value.(*FrequencyNode)
	}

	nextFreq.Add(key)
	elem.Parent = nextNode

	freq.Remove(key)
	if freq.Len() == 0 {
		l.DeleteNode(node)
	}
	return elem.Data, nil
}

func (l LfuCache) Store(key string, value interface{}) error {
	var node *list.Element
	var freq *FrequencyNode

	if l.dataLookup[key] != nil {
		return errors.New("Data already in cache")
	}

	node = l.frequencyHead.Next()
	if node != nil {
		freq = node.Value.(*FrequencyNode)
	}

	if node == nil || freq.Value != 1 {
		node, _ = l.GetNewNode(1, l.frequencyHead)
		freq = node.Value.(*FrequencyNode)
	}

	freq.Add(key)
	l.dataLookup[key] = NewCacheItem(value, node)
	return nil
}
