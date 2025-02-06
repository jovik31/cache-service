package cache


const(

	// offset64 is the offset basis for the FNV hash
	offset64 = 14695981039346656037 // from the FNV-1a spec 64 bit offset parameter| Decimal representation
	 
	prime64 = 1099511628211 // from the FNV-1a spec 64 bit prime parameter | Decimal representation

)

type Hasher interface {
	Hash(key string) uint64
}


type fnv64a struct{}


func newHasher() Hasher {
	return fnv64a{}
}

func (f fnv64a) Hash(key string) uint64 {

	var hash uint64 = offset64

	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash

}