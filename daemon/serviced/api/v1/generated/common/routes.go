// Package common provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package common

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns OK if healthy.
	// (GET /health)
	HealthCheck(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// HealthCheck converts echo context to params.
func (w *ServerInterfaceWrapper) HealthCheck(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.HealthCheck(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/health", wrapper.HealthCheck)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/4xU32vqShD+V4a5F/oS1Pa+5a2UC0da6KFaDgfrwzYZzR6T2XR2okjJ/37YTWqtBqoP",
	"orvffPPr+/YdM1fVjonVY/qOPiuoMvHnNA/fOflMbK3WMaY4ZaU1CQjVQp5YLa9BC4Kdkw3JlYdpjgnq",
	"viZM0XZobBN8oreGvAbGWlxNopZiFpv78zS3ImYPbgU9g/8uoVWqIs+/QitM8Z/xZ1vjvqfxNA+V9LWZ",
	"kCL+txV5NVX9Xfj8AGzbBIXeGiuUY7o4okhiP8tDFvf6hzLtBuBrx55iz5dU+xEwVaqG6v5yfzbBx5gY",
	"dlaL03GB4RxmarTxmJxt47Ih+i78G2yf5HRcNmysZxga1exA7pgeV5gugko5qGeSnPT5zGZrbGleS8I2",
	"OeCuT3EPLjPlMeLmFPFElVPCdjkk3/mxSC60RNaIECsc1DGC/vPCv10jMCPZ2owgK122garxCq8ETZ0b",
	"pfxjddbD1pQNjV74we1I4FCKh9wBO+0jYs4vlKMBJ4ZdWF658z7mBQmBEQLdOVCz9iksMldVjpdRMgvf",
	"kS9feF5QQHzeF4bzsiuhR0G34CPsR/gA+CoYvIxtS/dQBGn6pqqM7MN7wHD7cwrqLo5Vq2Vo/M5kBcHK",
	"CfyKDvAQlBBaDoyY4JbEdwO4Hk1Gk7BsVxOb2mKK/8WjBGujRRTkOAt8mL63CY4LMqUW4XhN8WULTorc",
	"4enEH/H6rqBsg8EAnV0jzc1kMmDZ+1HIntPKNKWe3z/zht2O4X8RJ9FUnwN6Im2EPTzeg11BV9g+7t+s",
	"ffBctyhcRgH0h6cJ7iIGiPPaWY5jZFOFKfbRwT5fQ2aHHQxE9QvCdtn+DQAA///JZyG3ZwYAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
