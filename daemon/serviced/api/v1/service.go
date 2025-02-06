package v1

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gitlab.com/sibsfps/spc/spc-1/cache"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/model"
)

//TODO

// Service implements the ServerInterface
/*
Service will have the following dependencies:
1. Logger
2. Cache client
3. REST Client to worker database
4. Configs

*/
type ServiceInterface interface {
	Get(key string, requestTimestamp uint64) (value cache.Status, err error)
	Set(key string, value cache.Status, requestTimestamp uint64)
	Reset()
	Config() cache.Config
	//Load() error    Not implemented
}

func (v1 *Handlers) Cache(ctx echo.Context) error {

	/* - REFACTOR using REST Client for worker service
	- get all user ids from cache
	- add the misses to a slice
	- add the hits to a slice
	- get user statuses from worker service
	- set user statuses in cache
	- return user statuses

	*/

	log := v1.Log
	req := ctx.Request()
	if req == nil {
		log.Error("request can't be nil")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorf("could not read body: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	var request model.Request
	if err := json.Unmarshal(reqBytes, &request); err != nil {
		log.Errorf("could not unmarshal request: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	// get all user ids from request
	ids := request.Ids
	//missedIds := []int{}
	hitIds := map[int]cache.Status{}

	// get all user statuses from cache
	for _, id := range ids {
		cacheVal, err := v1.CacheNode.Get(strconv.Itoa(id), uint64(request.Timestamp))
		if err != nil {
			if err == cache.ErrCacheMiss {
				//missedIds = append(missedIds, id)
			} else {
				log.Errorf("could not get cache value for id %d: %v", id, err)
				return ctx.NoContent(http.StatusInternalServerError)
			}
		}
		hitIds[id] = cacheVal
	}

	/* REFACTOR using REST Client for worker service
	for _, id := range missedIds {

		/* REFACTOR using REST Client for worker service
		workerStatus, err := v1.WorkerClient.Get(id)
		v1.CacheNode.Set(strconv.Itoa(id), cache.Status(1), uint64(request.Timestamp))
		hitIds[id] = cache.Status(workerStatus)

	}
	var response model.Response
	for k, v := range hitIds {
		var responseItem model.ResponseItem
		responseItem.Id = model.Id(k)
		responseItem.Status = model.Status(v)
		response = append(response, responseItem)
	}

	respBytes, err := json.Marshal(response)
	*/

	return nil
}
