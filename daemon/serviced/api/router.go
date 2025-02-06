package api

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cache "gitlab.com/sibsfps/spc/spc-1/cache"

	v1 "gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/common"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/service"

	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/lib/middlewares"
	"gitlab.com/sibsfps/spc/spc-1/logging"
)

const (
	BaseURL   = "v1"
	WorkerURL = "http://localhost:8080"
)

type APICache struct {
	*cache.CacheNode
}

type APICacheInterface interface {
	v1.CacheInterface
}

func NewRouter(logger logging.Logger, cache APICacheInterface, shutdown <-chan struct{}, listener net.Listener) *echo.Echo {

	e := echo.New()
	e.Logger = logger.MakeEchoLogger()
	e.Listener = listener
	e.HideBanner = true

	e.Pre(

		middleware.RemoveTrailingSlash(),
	)

	e.Use(
		middlewares.MakeLogger(logger),
	)

	v1Handler := v1.Handlers{
		CacheNode: cache,
		Log:       logger,
		Shutdown:  shutdown,
		//Client:    workerRC,
	}

	common.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)
	service.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)

	return e
}
