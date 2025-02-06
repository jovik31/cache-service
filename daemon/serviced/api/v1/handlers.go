package v1

import (
	"gitlab.com/sibsfps/spc/spc-1/logging"
	
)

// TODO

type CacheInterface interface {
	CommonInterface
	ServiceInterface
}

type Handlers struct {
	CacheNode CacheInterface
	Log       logging.Logger
	Shutdown  <-chan struct{}
	//Client    RestClient
}
