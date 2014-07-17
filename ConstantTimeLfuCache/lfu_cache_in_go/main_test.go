package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddingAndRemovingAnElementFromFrequencyNode(t *testing.T) {
	fn := NewFrequencyNode(1)
	assert.Equal(t, 0, fn.Len())
	fn.Add("hello")
	assert.Equal(t, 1, fn.Len())
}

func TestRemovingNonExistentElementFromFrequencyNode(t *testing.T) {
	fn := NewFrequencyNode(1)
	fn.Add("hello")
	fn.Remove("Hello")
	assert.Equal(t, 1, fn.Len())
	fn.Remove("hello")
	assert.Equal(t, 0, fn.Len())
	fn.Remove("hello")
	assert.Equal(t, 0, fn.Len())
}

func TestInitializingACache(t *testing.T) {
	cache := NewLfuCache()
	head := cache.Head()
	fn := head.Value.(*FrequencyNode)
	assert.Equal(t, 0, fn.Value)
	assert.Equal(t, 0, fn.Len())
	assert.Nil(t, head.Next())
	assert.Nil(t, head.Prev())
}

func TestGettingANewNodeForTheCache(t *testing.T) {
	cache := NewLfuCache()
	head := cache.Head()
	node, err := cache.GetNewNode(1, head)
	fn := node.Value.(*FrequencyNode)
	assert.Nil(t, err)
	assert.Equal(t, 1, fn.Value)
	assert.Equal(t, head, node.Prev())
	assert.Nil(t, node.Next())
}

func TestRemovingANodeFromTheCache(t *testing.T) {
	cache := NewLfuCache()
	assert.Equal(t, 0, cache.FreqListLen())
	node, err := cache.GetNewNode(1, cache.Head())
	assert.Nil(t, err)
	assert.Equal(t, 1, cache.FreqListLen())
	cache.DeleteNode(node)
	assert.Equal(t, 0, cache.FreqListLen())
}

func TestFetchingDataForAKeyThatDoesNotExist(t *testing.T) {
	cache := NewLfuCache()
	data, err := cache.Fetch("hello")
	assert.Nil(t, data)
	assert.Equal(t, "'hello' does not exist", err.Error())
}

func TestInsertingDataIntoTheCache(t *testing.T) {
	cache := NewLfuCache()
	cache.Store("hello", "world")
	actual, _ := cache.Fetch("hello")
	assert.Equal(t, "world", actual)
}

func TestInsertingIntoCacheStoresItemInCorrectLocation(t *testing.T) {
	cache := NewLfuCache()
	cache.Store("hello", "world")
	freq := cache.FreqNodeFor(1)
	assert.Equal(t, 1, freq.Value)
	assert.Equal(t, 1, freq.Len())
	assert.True(t, freq.Contains("hello"))
}

func TestInsertingSameValueTwiceReturnsAnError(t *testing.T) {
	cache := NewLfuCache()
	err := cache.Store("hello", "world")
	assert.Nil(t, err)
	err = cache.Store("hello", "world")
	assert.NotNil(t, err)
	assert.Equal(t, "Data already in cache", err.Error())
}
