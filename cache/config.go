package cache

import "gitlab.com/sibsfps/spc/spc-1/logging"

type Config struct {
	Shards uint64
	Hasher Hasher
	log    logging.Logger
}

func DefaultConfig() Config {
	return Config{
		Shards: shardNum,
		Hasher: newHasher(),
		log:    logging.NewLogger(),
	}
}
