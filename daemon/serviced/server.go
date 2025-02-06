package serviced

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api"
	"gitlab.com/sibsfps/spc/spc-1/logging"

	cache "gitlab.com/sibsfps/spc/spc-1/cache"
)

type ServerCache interface {
	api.APICacheInterface
	Config() cache.Config
}

type Server struct {
	log      logging.Logger
	cache    ServerCache
	stopping chan struct{}
}

func (s *Server) Initialize(cfg cache.Config) error {
	s.log = logging.Base()

	var serverCache ServerCache

	cacheConfig := cache.DefaultConfig()             // get default cache config
	thisCache := cache.MakeCache(s.log, cacheConfig) // probably needs error checking if failed
	serverCache = api.APICache{CacheNode: thisCache}

	s.cache = serverCache
	logging.RegisterExitHandler(s.cache.Reset)

	return nil
}

func makeListener() (net.Listener, error) {
	return net.Listen("tcp", ":8090")
}

func (s *Server) Start() {
	s.log.Info("Starting cache server")
	s.stopping = make(chan struct{})

	listener, err := makeListener()
	if err != nil {
		s.log.Fatalf("Failed to create listener: %v", err)
		return
	}

	addr := listener.Addr().String()
	server := http.Server{
		Addr: "127.0.0.1:8090",
	}

	e := api.NewRouter(s.log, s.cache, s.stopping, listener)
	errChan := make(chan error, 1)

	e.Logger = s.log.MakeEchoLogger()
	go func() {

		err := e.StartServer(&server)
		errChan <- err

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGHUP)

	fmt.Printf("Cache service running and accepting requests over HTTP on port %v. Press Ctrl-C to exit\n", addr)
	select {
	case err := <-errChan:
		if err != nil {
			s.log.Warn(err)
		} else {
			s.log.Info("Cache exited successfully")
		}
		s.Stop()
	case sig := <-c:
		s.log.Infof("Exiting on %v", sig)
		s.Stop()
		os.Exit(0)
	}
}
func (s *Server) Stop() {
}
