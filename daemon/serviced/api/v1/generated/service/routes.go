// Package service provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package service

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
	// Process a query to the cache
	// (POST /cache)
	Cache(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Cache converts echo context to params.
func (w *ServerInterfaceWrapper) Cache(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Cache(ctx)
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

	router.POST(baseURL+"/cache", wrapper.Cache)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5xUTW/bRhD9K4NpgVxYffXGW5r2ICBA0thBUMg6rMiRuTW5Q+8OrQiC/nuxsxQlSwIk",
	"NIfA4s7nm/feDgtuWnbkJGC+w1BU1Bj9c17G/0sKhbetWHaY49wJPZMHT62nQE6sewapCDbsX8h/CDAv",
	"MUPZtoQ52hSN+wy/0WtHQWLF1nNLXixpF1uGyzYfvTdb4DX0FcKthlao0Tq/elpjjr+Mj2uN+53G8zJO",
	"0s9mYgv9bRsKYpr2VvrjELjfZ+jptbOeSswXJyUy3Wc5dOHVv1RIAiC07ALpzvdMe0iYCzXX5n73foHg",
	"F20MGyvVOVxgXAkPYqQLmF1c4z4QQ0q/Eds3OYfLxov1Fa5B9TAUZ0df1pgvIktdZM8kO9vzuzNvxtZm",
	"VRPusyFueh73mQtTn0bMziO+UcNCuF9eo+/jKUnulETReU9OYGDHCPp/T+4f7jw8kH+zBUFRc/ECTRcE",
	"VgRdWxqh8nA6G+DN1B2Nntxn3pCHYZQAJYNj6TO057uSoytKjLewbs2XezxW5AmMJ5ANg5jnkMOi4KZh",
	"t1TKLEIqvnxyjxXFiON7ZVxZpxH6KEgHPok9pF8J/hAFXuvaPhlFpGbomsb4bfQDBx+/zkH47lyxUsfF",
	"P5miIlizhx+qgACRCXHlWBEzfCMfEgDT0WQ0icfmlpxpLeb4u37KsDVSKSHHRaynLsbJzc6M6zAD2ACq",
	"iEAlrLZgoLZBoqX9GJT4Z1BczQlDULt7HTD6b5ofk34oyB9cbmPXgp2Q0wFM29Y2rTT++dvKOoVsh/TT",
	"NG2dZP13R+njEVL9NIAy1UlmEY9INsxxsVTdJinf9qrk7kqvC1rBaisE5Aouj0fCU08Q35GaRLI0nXk2",
	"mdy56c2GqWoE9xSUU0c+4pIQgWnS30CWiQLUP87OHqf/H7d+hCvAzSaTYXTU17Xp6iuc++5eHG8c/OU9",
	"+9R82Oar54JCAAOvem/hZE09q6LKoyP3UsJlzB5XZGqpMN/FX33IedNPKnwgV7ZsnUrOmSYCkCxBrfZ9",
	"ysOg1ytZhwn2y/1/AQAA///qDw5FkwgAAA==",
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
