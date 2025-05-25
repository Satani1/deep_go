package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type OrderedMap struct {
	root *Node
	size int
}

type Node struct {
	key         int
	value       int
	left, right *Node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	m.root = m.insert(key, value, m.root)
}

func (m *OrderedMap) insert(key, value int, node *Node) *Node {
	if node == nil {
		m.size++
		return &Node{key: key, value: value}
	}

	if key < node.key {
		node.left = m.insert(key, value, node.left)
	}

	if key > node.key {
		node.right = m.insert(key, value, node.right)
	}

	node.value = value
	return node
}
func (m *OrderedMap) Erase(key int) {
	m.root = m.deleteNode(key, m.root, m.size)
}

func (m *OrderedMap) deleteNode(key int, node *Node, size int) *Node {
	if node == nil {
		return nil
	}
	switch {
	case key < node.key:
		node.left = m.deleteNode(key, node.left, m.size)
	case key > node.key:
		node.right = m.deleteNode(key, node.right, m.size)
	default:
		m.size--
		if node.left == nil {
			return node.right
		}
		if node.right == nil {
			return node.left
		}

		succ := node.right
		for succ.left != nil {
			succ = succ.left
		}

		node.key, node.value = succ.key, succ.value

		node.right = m.deleteNode(succ.key, node.right, m.size)
	}
	return node
}

func (m *OrderedMap) Contains(key int) bool {
	return contains(key, m.root)
}

func contains(key int, node *Node) bool {
	if node == nil {
		return false
	}

	if key < node.key {
		return contains(key, node.left)
	}

	if key > node.key {
		return contains(key, node.right)
	}

	return true
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.root != nil {
		forEach(action, m.root)
	}
}

func forEach(action func(int, int), node *Node) {
	if node.left != nil {
		forEach(action, node.left)
	}

	action(node.key, node.value)

	if node.right != nil {
		forEach(action, node.right)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
