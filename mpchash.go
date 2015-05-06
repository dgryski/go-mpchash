// Package mpchash implements a multi-probe consistent hashing strategy
/*
http://arxiv.org/pdf/1505.00062.pdf
*/
package mpchash

import (
	"math"
	"sort"
)

type Multi struct {
	buckets []string
	seeds   []uint64
	hashf   func(b []byte, s uint64) uint64

	bmap    map[uint64]string
	bhashes []uint64
}

func New(buckets []string, h func(b []byte, s uint64) uint64, seeds []uint64) *Multi {
	m := &Multi{
		buckets: buckets,
		hashf:   h,
		seeds:   seeds,
		bhashes: make([]uint64, len(buckets)),
		bmap:    make(map[uint64]string, len(buckets)),
	}

	for i, b := range buckets {
		h := m.hashf([]byte(b), 0)
		m.bhashes[i] = h
		m.bmap[h] = b
		m.buckets[i] = b
	}

	sort.Sort(uint64Slice(m.bhashes))

	return m
}

func (m *Multi) Hash(key string) string {

	bkey := []byte(key)

	minDistance := uint64(math.MaxUint64)
	var minIdx int

	for _, seed := range m.seeds {
		hash := m.hashf(bkey, seed)

		idx := sort.Search(len(m.bhashes), func(i int) bool { return m.bhashes[i] >= hash })

		// Means we have cycled back to the first replica.
		if idx == len(m.bhashes) {
			idx = 0
		}

		distance := m.bhashes[idx] - hash
		if distance < minDistance {
			minDistance = distance
			minIdx = idx
		}
	}

	return m.bmap[m.bhashes[minIdx]]
}

type uint64Slice []uint64

func (p uint64Slice) Len() int           { return len(p) }
func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
