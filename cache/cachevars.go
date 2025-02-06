package cache

type Status = uint64
type TTL = uint64

// User location status
const (
	UserUnavailableStatus Status = iota
	UserLocalStatus
	UserRemoteStatus
)

// Status dependent TTLs
const (
	UserUnavailableTTL TTL = (4 * 60 * 60) // 4 hours
	UserLocalRemoteTTL TTL = (8 * 60 * 60) // 8 hours
)

const (
	shardNum  uint64 = 256  // Number of shards
	shardSize uint64 = 2048 // Number of users per shard
)

var (
	statusTTL = map[Status]TTL{

		UserUnavailableStatus: UserUnavailableTTL,
		UserLocalStatus:       UserLocalRemoteTTL,
		UserRemoteStatus:      UserLocalRemoteTTL,
	}
)
