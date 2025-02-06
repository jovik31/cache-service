package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/sibsfps/spc/spc-1/logging"
)

func TestSetandGet(t *testing.T) {

	t.Parallel()

	log := logging.Base()
	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())
	cache.Set("testUser", UserLocalStatus, requestTimestamp)
	cachedVal, _ := cache.Get("testUser", requestTimestamp)

	assert.Equal(t, UserLocalStatus, cachedVal)
}

func TestCacheMiss(t *testing.T) {

	t.Parallel()

	log := logging.Base()
	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())
	_, err := cache.Get("nonExistingUser", requestTimestamp)

	assert.Equal(t, ErrCacheMiss, err)
}

func TestEvictionUnavailableTTL(t *testing.T) {

	t.Parallel()

	log := logging.Base()

	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())
	cache.Set("expiredUser", UserUnavailableStatus, requestTimestamp)
	_, err := cache.Get("expiredUser", requestTimestamp+UserUnavailableTTL+1) // Unavailable user TTL expired

	assert.Equal(t, ErrCacheMiss, err)
}

func TestEvictionLocalTTL(t *testing.T) {

	t.Parallel()

	log := logging.Base()

	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())

	cache.Set("localExpiredUser", UserLocalStatus, requestTimestamp)
	_, err := cache.Get("localExpiredUser", requestTimestamp+UserLocalRemoteTTL+1) // Unavailable user TTL expired

	assert.Equal(t, ErrCacheMiss, err)

}

func TestEvictionRemoteTTL(t *testing.T) {

	t.Parallel()

	log := logging.Base()

	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())

	cache.Set("remoteExpiredUser", UserRemoteStatus, requestTimestamp)
	_, err := cache.Get("remoteExpiredUser", requestTimestamp+UserLocalRemoteTTL+1) // Unavailable user TTL expired

	assert.Equal(t, ErrCacheMiss, err)

}

func BenchmarkSet(b *testing.B) {

	log := logging.Base()

	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())
	value := Status(1)
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("user_%d", i), value, requestTimestamp)
	}
}

func BenchmarkGet(b *testing.B) {

	b.StopTimer()
	log := logging.Base()
	requestTimestamp := uint64(1000000)

	cache := MakeCache(log, DefaultConfig())
	value := Status(0)
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("user_%d", i), value, requestTimestamp)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("user_%d", i), requestTimestamp)
	}
}
