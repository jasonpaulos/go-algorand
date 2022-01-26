// Package private provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package private

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/algorand/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Aborts a catchpoint catchup.
	// (DELETE /v2/catchup/{catchpoint})
	AbortCatchup(ctx echo.Context, catchpoint string) error
	// Starts a catchpoint catchup.
	// (POST /v2/catchup/{catchpoint})
	StartCatchup(ctx echo.Context, catchpoint string) error
	// Return a list of participation keys
	// (GET /v2/participation)
	GetParticipationKeys(ctx echo.Context) error
	// Add a participation key to the node
	// (POST /v2/participation)
	AddParticipationKey(ctx echo.Context) error
	// Delete a given participation key by ID
	// (DELETE /v2/participation/{participation-id})
	DeleteParticipationKeyByID(ctx echo.Context, participationId string) error
	// Get participation key info given a participation ID
	// (GET /v2/participation/{participation-id})
	GetParticipationKeyByID(ctx echo.Context, participationId string) error
	// Append state proof keys to a participation key
	// (POST /v2/participation/{participation-id})
	AppendKeys(ctx echo.Context, participationId string) error

	// (POST /v2/shutdown)
	ShutdownNode(ctx echo.Context, params ShutdownNodeParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AbortCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) AbortCatchup(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameter("simple", false, "catchpoint", ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AbortCatchup(ctx, catchpoint)
	return err
}

// StartCatchup converts echo context to params.
func (w *ServerInterfaceWrapper) StartCatchup(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "catchpoint" -------------
	var catchpoint string

	err = runtime.BindStyledParameter("simple", false, "catchpoint", ctx.Param("catchpoint"), &catchpoint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter catchpoint: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.StartCatchup(ctx, catchpoint)
	return err
}

// GetParticipationKeys converts echo context to params.
func (w *ServerInterfaceWrapper) GetParticipationKeys(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetParticipationKeys(ctx)
	return err
}

// AddParticipationKey converts echo context to params.
func (w *ServerInterfaceWrapper) AddParticipationKey(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddParticipationKey(ctx)
	return err
}

// DeleteParticipationKeyByID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteParticipationKeyByID(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameter("simple", false, "participation-id", ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteParticipationKeyByID(ctx, participationId)
	return err
}

// GetParticipationKeyByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetParticipationKeyByID(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameter("simple", false, "participation-id", ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetParticipationKeyByID(ctx, participationId)
	return err
}

// AppendKeys converts echo context to params.
func (w *ServerInterfaceWrapper) AppendKeys(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameter("simple", false, "participation-id", ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set("api_key.Scopes", []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AppendKeys(ctx, participationId)
	return err
}

// ShutdownNode converts echo context to params.
func (w *ServerInterfaceWrapper) ShutdownNode(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty":  true,
		"timeout": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error

	ctx.Set("api_key.Scopes", []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params ShutdownNodeParams
	// ------------- Optional query parameter "timeout" -------------
	if paramValue := ctx.QueryParam("timeout"); paramValue != "" {

	}

	err = runtime.BindQueryParameter("form", true, false, "timeout", ctx.QueryParams(), &params.Timeout)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter timeout: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ShutdownNode(ctx, params)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE("/v2/catchup/:catchpoint", wrapper.AbortCatchup, m...)
	router.POST("/v2/catchup/:catchpoint", wrapper.StartCatchup, m...)
	router.GET("/v2/participation", wrapper.GetParticipationKeys, m...)
	router.POST("/v2/participation", wrapper.AddParticipationKey, m...)
	router.DELETE("/v2/participation/:participation-id", wrapper.DeleteParticipationKeyByID, m...)
	router.GET("/v2/participation/:participation-id", wrapper.GetParticipationKeyByID, m...)
	router.POST("/v2/participation/:participation-id", wrapper.AppendKeys, m...)
	router.POST("/v2/shutdown", wrapper.ShutdownNode, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
