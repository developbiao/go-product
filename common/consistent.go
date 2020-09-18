package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Declare unsigned slice
type units []uint32

// Return slice length
func (x units) Len() int {
	return len(x)
}

// Compare tow number
func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

// Swap value
func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// Declare empty hash circle error
var errEmpty = errors.New("Hash circle not data")

// Declare consistent struct save hash information
type Consistent struct {
	// hash circle, key is hash value, value save node information
	circle map[uint32]string

	// already sorted node hash slice
	sortedHashes units

	// virtual node amount LSB bash
	VirtualNode int
	// map write/read lock
	sync.RWMutex
}

// Create consistent algorithm
func NewConsistent() *Consistent {
	return &Consistent{
		// init variable
		circle: make(map[uint32]string),
		// default virtual node is 20
		VirtualNode: 20,
	}
}

// Automatic generator key
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

// Get Hash key position
func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		// declare array length is 64
		var srcatch [64]byte
		// copy key to array
		copy(srcatch[:], key)
		// use IEEE CRC sum check
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// Update sort for quick search
func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	// judge slice capacity is so big well be reset it
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	// Add hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	// sort all node hash value
	sort.Sort(hashes)
	// reassign
	c.sortedHashes = hashes
}

// Add node to hash circle
func (c *Consistent) Add(element string) {
	// lock
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		// compute node key and add node to hash circel
		c.circle[c.hashKey(c.generateKey(element, i))] = element

		// Update sort
		c.updateSortedHashes()
	}
}

// Delete a node
func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)

}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

// Search by clockwise neighbour
func (c *Consistent) search(key uint32) int {
	// query algorithm
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}

	// Use binary search
	i := sort.Search(len(c.sortedHashes), f)

	// if out of range, set i= 0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// Get server node by identify
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return "", errEmpty
	}

	// compute hash value
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
