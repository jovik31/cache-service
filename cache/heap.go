package cache

import (
	"container/heap"
)

// Heap element representing a user in the shard
type Element struct {
	UserID    uint64
	timestamp uint64 // Unix timestamp - priority
	index     int
}

type CacheShardHeap []*Element

func (ch CacheShardHeap) Len() int { return len(ch) }

func (ch CacheShardHeap) Less(i, j int) bool {

	return ch[i].timestamp < ch[j].timestamp
}

func (ch CacheShardHeap) Swap(i, j int) {

	ch[i], ch[j] = ch[j], ch[i]
	ch[i].index = i
	ch[j].index = j
}

func (ch CacheShardHeap) Peek() *Element {

	if len(ch) == 0 {
		return nil
	}
	return ch[0]
}

func (ch *CacheShardHeap) Push(x any) {

	n := len(*ch)
	item := x.(*Element)
	item.index = n
	*ch = append(*ch, item)

}

func (ch *CacheShardHeap) Pop() any {
	old := *ch
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*ch = old[0 : n-1]
	return item
}

func (ch *CacheShardHeap) Update(item *Element, timestamp uint64) {
	item.timestamp = timestamp
	heap.Fix(ch, item.index)
}

func initCacheShardHeap() CacheShardHeap {
	csh := make(CacheShardHeap, 0)
	heap.Init(&csh)
	return csh
}
