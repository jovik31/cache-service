package cache

import (
	"gitlab.com/sibsfps/spc/spc-1/logging"
)

type Cache interface {
	Get(key string) (Status, error)
	Set(key string, value Status, timestamp uint64)
	Reset()
	Config() Config
}

type CacheNode struct {
	shards    []*cacheShard
	clock     uint64
	hash      Hasher
	shardMask uint64
	log       logging.Logger
	config    Config
}

type CacheStatus struct {
}

func MakeCache(log logging.Logger, config Config) (*CacheNode) {

	cache := &CacheNode{
		clock:     0,
		hash:      config.Hasher,
		shardMask: (config.Shards - 1),
		config:    config,
		log:       config.log.With("cache", "internal"),
	}

	cache.shards = make([]*cacheShard, shardNum)

	for i := 0; i < int(shardNum); i++ {
		cache.shards[i] = initShard()
	}
	return cache
}

func (c *CacheNode) getShard(hashedKey uint64) *cacheShard {

	return c.shards[hashedKey&c.shardMask]

}

func (c *CacheNode) Get(key string, requestTimestamp uint64) (value Status, err error) {

	hashedKey := c.hash.Hash(key)
	shard := c.getShard(hashedKey)
	value, err = shard.get(hashedKey, requestTimestamp)
	if err != nil {

		return value, err // cache miss propagation

	}
	return value, nil

}

func (c *CacheNode) Set(key string, value Status, requestTimestamp uint64) {

	hashedKey := c.hash.Hash(key)
	shard := c.getShard(hashedKey)
	shard.set(hashedKey, value, requestTimestamp)

}

func (c *CacheNode) Reset() {
	for _, shard := range c.shards {
		shard.reset()
	}

}

func (c *CacheNode) Config() Config {
	return c.config
}
