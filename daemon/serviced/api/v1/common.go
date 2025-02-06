package v1

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	cache "gitlab.com/sibsfps/spc/spc-1/cache"
)

// TODO

// HealthServer implements the ServerInterface
type CommonInterface interface {
	Config() cache.Config
}

// HealthCheck will return the status of the server
func (h *Handlers) HealthCheck(ctx echo.Context) error {

	w := ctx.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})

	return nil
}
