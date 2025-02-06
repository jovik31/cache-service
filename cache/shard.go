package cache

import (
	"sync"
)

type cacheShard struct {
	users   map[uint64]Status // map of hashed keys to user status
	lock    sync.RWMutex
	ttlHeap CacheShardHeap
	clock   uint64
}

func initShard() *cacheShard {
	return &cacheShard{

		users:   make(map[uint64]Status, shardSize), // init map of users
		ttlHeap: initCacheShardHeap(),               // init heap of shard
		clock:   0,                                  // init clock at 0 time so the first request has always a higher timestamp
	}
}

func (cs *cacheShard) set(hashedKey uint64, status Status, timestamp uint64) {
	cs.lock.Lock()

	cs.users[hashedKey] = status // set the status of the user

	// push the user to the heap
	cs.push(hashedKey, timestamp + statusTTL[status]) // already ensured that statusTTL[status] exists

	if cs.clock < timestamp {
		cs.clock = timestamp
	}

	cs.eviction() // evict users with lower timestamp

	cs.lock.Unlock()

}

func (cs *cacheShard) get(hashedKey uint64, timestamp uint64) (Status, error) {
	cs.lock.RLock()

	if cs.clock < timestamp {
		cs.clock = timestamp
	}
	cs.eviction() // process shard heap eviction

	// cache miss - no user
	status, ok := cs.users[hashedKey]
	if !ok {
		cs.lock.RUnlock()
		return status, ErrCacheMiss
	}
	cs.lock.RUnlock()
	return status, nil
}

func (cs *cacheShard) push(hashedKey uint64, timestamp uint64) {

	cs.ttlHeap.Push(&Element{UserID: hashedKey, timestamp: timestamp})
}

func (cs *cacheShard) eviction() {
	// Evict all the users with a lower timestamp than the current shard clock

	for {
		user := cs.ttlHeap.Peek()
		if user == nil {
			break
		}
		if user.timestamp >= cs.clock {
			break
		}
		user = cs.ttlHeap.Pop().(*Element) 
		delete(cs.users, user.UserID)
	}

}

func (cs *cacheShard) reset() {
	cs.lock.Lock()
	cs.users = make(map[uint64]Status, shardSize)
	cs.ttlHeap = initCacheShardHeap()
	cs.clock = 0
	cs.lock.Unlock()
}
