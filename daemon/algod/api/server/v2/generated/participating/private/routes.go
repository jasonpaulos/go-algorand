// Package private provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package private

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	. "github.com/algorand/go-algorand/daemon/algod/api/server/v2/generated/model"
	"github.com/algorand/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
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
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetParticipationKeys converts echo context to params.
func (w *ServerInterfaceWrapper) GetParticipationKeys(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetParticipationKeys(ctx)
	return err
}

// AddParticipationKey converts echo context to params.
func (w *ServerInterfaceWrapper) AddParticipationKey(ctx echo.Context) error {
	var err error

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddParticipationKey(ctx)
	return err
}

// DeleteParticipationKeyByID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteParticipationKeyByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "participation-id", runtime.ParamLocationPath, ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteParticipationKeyByID(ctx, participationId)
	return err
}

// GetParticipationKeyByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetParticipationKeyByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "participation-id", runtime.ParamLocationPath, ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetParticipationKeyByID(ctx, participationId)
	return err
}

// AppendKeys converts echo context to params.
func (w *ServerInterfaceWrapper) AppendKeys(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "participation-id" -------------
	var participationId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "participation-id", runtime.ParamLocationPath, ctx.Param("participation-id"), &participationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter participation-id: %s", err))
	}

	ctx.Set(Api_keyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AppendKeys(ctx, participationId)
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
func RegisterHandlers(router EchoRouter, si ServerInterface, m ...echo.MiddlewareFunc) {
	RegisterHandlersWithBaseURL(router, si, "", m...)
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v2/participation", wrapper.GetParticipationKeys, m...)
	router.POST(baseURL+"/v2/participation", wrapper.AddParticipationKey, m...)
	router.DELETE(baseURL+"/v2/participation/:participation-id", wrapper.DeleteParticipationKeyByID, m...)
	router.GET(baseURL+"/v2/participation/:participation-id", wrapper.GetParticipationKeyByID, m...)
	router.POST(baseURL+"/v2/participation/:participation-id", wrapper.AppendKeys, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+x9+3PcNtLgv4Ka76vy44Yz8iPZtapS38mWk9XF8bosJXv32b4EQ/bMYEUCXACUZuLT",
	"/36FBkCCJDhDPSJvqvyTrSEejUaj0W98nqSiKAUHrtXk8POkpJIWoEHiXzRNRcV1wjLzVwYqlazUTPDJ",
	"of9GlJaMrybTCTO/llSvJ9MJpwU0bUz/6UTCvyomIZscalnBdKLSNRTUDKy3pWldj7RJViJxQxzZIU6O",
	"J1c7PtAsk6BUH8q/83xLGE/zKgOiJeWKpuaTIpdMr4leM0VcZ8I4ERyIWBK9bjUmSwZ5pmZ+kf+qQG6D",
	"VbrJh5d01YCYSJFDH85XolgwDh4qqIGqN4RoQTJYYqM11cTMYGD1DbUgCqhM12Qp5B5QLRAhvMCrYnL4",
	"YaKAZyBxt1JgF/jfpQT4HRJN5Qr05NM0trilBploVkSWduKwL0FVuVYE2+IaV+wCODG9ZuSnSmmyAEI5",
	"ef/9K/Ls2bMXZiEF1RoyR2SDq2pmD9dku08OJxnV4D/3aY3mKyEpz5K6/fvvX+H8p26BY1tRpSB+WI7M",
	"F3JyPLQA3zFCQoxrWOE+tKjf9IgciubnBSyFhJF7Yhvf6aaE83/RXUmpTtelYFxH9oXgV2I/R3lY0H0X",
	"D6sBaLUvDaakGfTDQfLi0+cn0ycHV//x4Sj5b/fnN8+uRi7/VT3uHgxEG6aVlMDTbbKSQPG0rCnv4+O9",
	"owe1FlWekTW9wM2nBbJ615eYvpZ1XtC8MnTCUimO8pVQhDoyymBJq1wTPzGpeG7YlBnNUTthipRSXLAM",
	"sqnhvpdrlq5JSpUdAtuRS5bnhgYrBdkQrcVXt+MwXYUoMXDdCB+4oH9fZDTr2oMJ2CA3SNJcKEi02HM9",
	"+RuH8oyEF0pzV6nrXVbkbA0EJzcf7GWLuOOGpvN8SzTua0aoIpT4q2lK2JJsRUUucXNydo793WoM1gpi",
	"kIab07pHzeEdQl8PGRHkLYTIgXJEnj93fZTxJVtVEhS5XINeuztPgioFV0DE4p+QarPt/+v072+JkOQn",
	"UIqu4B1NzwnwVGTDe+wmjd3g/1TCbHihViVNz+PXdc4KFgH5J7phRVUQXhULkGa//P2gBZGgK8mHALIj",
	"7qGzgm76k57Jiqe4uc20LUHNkBJTZU63M3KyJAXdfHcwdeAoQvOclMAzxldEb/igkGbm3g9eIkXFsxEy",
	"jDYbFtyaqoSULRlkpB5lByRumn3wMH49eBrJKgDHDzIITj3LHnA4bCI0Y46u+UJKuoKAZGbkZ8e58KsW",
	"58BrBkcWW/xUSrhgolJ1pwEYcerd4jUXGpJSwpJFaOzUocNwD9vGsdfCCTip4JoyDpnhvAi00GA50SBM",
	"wYS7lZn+Fb2gCr59PnSBN19H7v5SdHd9546P2m1slNgjGbkXzVd3YONiU6v/COUvnFuxVWJ/7m0kW52Z",
	"q2TJcrxm/mn2z6OhUsgEWojwF49iK051JeHwI39s/iIJOdWUZ1Rm5pfC/vRTlWt2ylbmp9z+9EasWHrK",
	"VgPIrGGNalPYrbD/mPHi7FhvokrDGyHOqzJcUNrSShdbcnI8tMl2zOsS5lGtyoZaxdnGaxrX7aE39UYO",
	"ADmIu5KahuewlWCgpekS/9kskZ7oUv5u/inL3PTW5TKGWkPH7r5F24CzGRyVZc5SapD43n02Xw0TAKsl",
	"0KbFHC/Uw88BiKUUJUjN7KC0LJNcpDRPlKYaR/pPCcvJ4eQ/5o1xZW67q3kw+RvT6xQ7GXnUyjgJLctr",
	"jPHOyDVqB7MwDBo/IZuwbA8lIsbtJhpSYoYF53BBuZ41+kiLH9QH+IObqcG3FWUsvjv61SDCiW24AGXF",
	"W9vwgSIB6gmilSBaUdpc5WJR//DwqCwbDOL3o7K0+EDREBhKXbBhSqtHuHzanKRwnpPjGfkhHBvlbMHz",
	"rbkcrKhh7oalu7XcLVYbjtwamhEfKILbKeTMbI1Hg5Hh74LiUGdYi9xIPXtpxTT+m2sbkpn5fVTnPweJ",
	"hbgdJi7UohzmrAKDvwSay8MO5fQJx9lyZuSo2/dmZGNGiRPMjWhl537acXfgsUbhpaSlBdB9sXcp46iB",
	"2UYW1lty05GMLgpzcIYDWkOobnzW9p6HKCRICh0YXuYiPf8bVes7OPMLP1b/+OE0ZA00A0nWVK1nk5iU",
	"ER6vZrQxR8w0RO2dLIKpZvUS72p5e5aWUU2DpTl442KJRT32Q6YHMqK7/B3/Q3NiPpuzbVi/HXZGzpCB",
	"KXucnQchM6q8VRDsTKYBmhgEKaz2TozWfS0oXzWTx/dp1B69tgYDt0NuEbhDYnPnx+Cl2MRgeCk2vSMg",
	"NqDugj7MOChGaijUCPiOHWQC99+hj0pJt30k49hjkGwWaERXhaeBhze+maWxvB4thLwZ9+mwFU4aezKh",
	"ZtSA+U47SMKmVZk4UozYpGyDzkCNC2830+gOH8NYCwunmv4BWFBm1LvAQnugu8aCKEqWwx2Q/jrK9BdU",
	"wbOn5PRvR988efrr02++NSRZSrGStCCLrQZFHjrdjCi9zeFRf2WoHVW5jo/+7XNvhWyPGxtHiUqmUNCy",
	"P5S1bloRyDYjpl0fa20046prAMcczjMwnNyinVjDvQHtmCkjYRWLO9mMIYRlzSwZcZBksJeYrru8Zppt",
	"uES5ldVdqLIgpZAR+xoeMS1SkScXIBUTEVfJO9eCuBZevC27v1toySVVxMyNpt+Ko0ARoSy94eP5vh36",
	"bMMb3Ozk/Ha9kdW5ecfsSxv53pKoSAky0RtOMlhUq5YmtJSiIJRk2BHv6DdstdaByPJOCrG881s7Okts",
	"SfjBCny56dMX+96KDIzaXak7YO/NYA32DOWEOKMLUWlCCRcZoI5eqTjjH3D0oocJHWM6vEv02spwCzD6",
	"YEors9qqJOj26dFi0zGhqaWiBFGjBuzitUPDtrLTWSdiLoFmRk8ETsTCGZ+dWRwXSdFnpT3rdNdORHNu",
	"wVVKkYJSRr+3Wtte0Hw7S5Z6B54QcAS4noUoQZZU3hrY84u9cJ7DNkEPqyIPf/xFPfoC8Gqhab4Hsdgm",
	"ht5ahXAehj7U46bfRXDdyUOyoxKI531GXzEMIgcNQyi8Fk4G968LUW8Xb4+WC5Bo6/9DKd5PcjsCqkH9",
	"g+n9ttBW5UDckBOdz1iBliBOuVCQCp6p6GA5VTrZx5ZNo5Z8b1YQcMIYJ8aBB6yRb6jS1j/FeIZqtb1O",
	"cB5rpjRTDAM8KOKYkX/x0k1/7NTcg1xVqhZ1VFWWQmrIYmvgsNkx11vY1HOJZTB2LU9pQSoF+0YewlIw",
	"vkOWXYlFENW1Gdc5cPuLQ2Onuee3UVS2gGgQsQuQU98qwG4YOzEACFMNoi3hMNWhnDpgYzpRWpSl4RY6",
	"qXjdbwhNp7b1kf65adsnLqqbezsTYGbXHiYH+aXFrI2aWVOjlOHIpKDnRvZAFcs60vowm8OYKMZTSHZR",
	"vjmWp6ZVeAT2HNIB7dbF5QWzdQ5Hh36jRDdIBHt2YWjBA6r2Oyo1S1mJkuKPsL1zwbk7QdQATDLQlBn1",
	"L/hghegy7E+sZ7Q75s0E6VFaUR/8nloUWU7OFF4YbeDPYYueoHc25OYsCNS5A00gMqo53ZQTBNQ78o0A",
	"EzaBDU11vjXXnF7DllyCBKKqRcG0tjFUbUVBizIJB4hanHbM6MyrNlzF78AYe+8pDhUsr78V04mVqHbD",
	"d9YRq1rocJJUKUQ+wtPWQ0YUglGeOFIKs+vMhez5uC5PSS0gnRCDtvWaeT5QLTTjCsj/ERVJKUeBtdJQ",
	"3whCIpvF69fMYC6wek7nc2swBDkUYOVw/PL4cXfhjx+7PWeKLOHSx7mahl10PH6MWvA7oXTrcN2BCcYc",
	"t5MIb0dTnLkonAzX5Sn7fT5u5DE7+a4zeG2/M2dKKUe4Zvm3ZgCdk7kZs/aQRsb5u3DcUVa2YOjYunHf",
	"MeDgj7HRNEPHoOtPHLhpm49DnlojX+XbO+DTdiAioZSg8FSFeomyX8UyDIV2x05tlYaib7qxXX8dEGze",
	"e7GgJ2UKnjMOSSE4bKPZP4zDT/gx1tue7IHOyGOH+nbFphb8HbDa84yhwtviF3c7IOV3dYjCHWx+d9yO",
	"1S4MAketFPKSUJLmDHVWwZWWVao/copScXCWI64cL+sP60mvfJO4YhbRm9xQHzlFN14tK0fNz0uIaMHf",
	"A3h1SVWrFSjdkQ+WAB+5a8U4qTjTOFdh9iuxG1aCRH/KzLYs6JYsaY5q3e8gBVlUun1jYqyq0kbrsiZE",
	"Mw0Ry4+capKD0UB/Yvxsg8P5kFBPMxz0pZDnNRZm0fOwAg6KqSTucvrBfsVoALf8tYsMwMQh+9kancz4",
	"TUDrVkMrGeb/Pvyvww9HyX/T5PeD5MX/mH/6/Pzq0ePej0+vvvvu/7V/enb13aP/+s/YTnnYY5GUDvKT",
	"YydNnhyjyNBYnXqw35vFoWA8iRLZ2RpIwTgG5Hdoizw0go8noEeNWc/t+keuN9wQ0gXNWUb1zcihy+J6",
	"Z9Gejg7VtDaio0D6tX6KRUesRFLS9Bw9tpMV0+tqMUtFMfdS9Hwlaol6nlEoBMdv2ZyWbK5KSOcXT/Zc",
	"6bfgVyTCrjpM9sYCQd/fG49+RoOqC2jGk7esuCWKSjmjLgb3eb+bWE7rCHeb2XpIMPx5Tb3T2P359Jtv",
	"J9MmbLn+bjR1+/VT5EywbBMLTs9gE5PU3FHDI/ZAkZJuFeg4H0LYoy5G65cKhy3AiPhqzcr75zlKs0Wc",
	"V/qQKafxbfgJt7FM5iSieXbrrD5ief9wawmQQanXsYy3lsyBrZrdBOi4zEopLoBPCZvBrKtxZStQ3tmZ",
	"A11i5hWaGMWYEND6HFhC81QRYD1cyCi1JkY/KCY7vn81nTgxQt25ZO8GjsHVnbO2xfq/tSAPfnh9RuaO",
	"9aoHNk/CDh1EtkcsGS54s+VMNdzM5vnaRJGP/CM/hiXjzHw//Mgzqul8QRVL1bxSIF/SnPIUZitBDn08",
	"6DHV9CPvyWyDqfhBJC4pq0XOUnIeytYNedr0yv4IHz9+MBz/48dPPb9SXxJ2U0X5i50guWR6LSqduPyx",
	"RMIllVkEdFXnD+HINvtz16xT4sa2rNjlp7nx4zyPlqXq5hH0l1+WuVl+QIbKRcmbLSNKC+mlGiPqWGhw",
	"f98KdzFIeumTDysFivxW0PID4/oTST5WBwfPgLQC639zwoOhyW0JLZvXjfIcuvYuXLjVkGCjJU1KugIV",
	"Xb4GWuLuo+RdoHU1zwl2awX0+4AlHKpZgMfH8AZYOK4dnIyLO7W9fCGA+BLwE24htjHiRuO0uOl+BSH+",
	"N96uTppAb5cqvU7M2Y6uShkS9ztT5wevjJDlPUmKrbg5BC6VegEkXUN6DhlmdUJR6u201d07K53I6lkH",
	"Uzb72QboYooemgcXQKoyo06op3zbzZVSoLVPEHsP57A9E02G33WSo9q5OmrooCKlBtKlIdbw2Loxupvv",
	"POKYn1CWPuUFY589WRzWdOH7DB9kK/LewSGOEUUrl2QIEVRGEGGJfwAFN1ioGe9WpB9bntFXFvbmiyRL",
	"e95PXJNGDXPO63A1mCJjvxeApRTEpSILauR24aoA2HyUgItViq5gQEIOLbQjsz5aVl0cZN+9F73pxLJ7",
	"ofXumyjItnFi1hylFDBfDKmgMtMJWfAzWScArmBGsLiPQ9giRzGpju2wTIfKlqXcVisZAi1OwCB5I3B4",
	"MNoYCSWbNVW+QAHWcfBneZQM8AfmV+3Kqj0JvO1BsYY6Z9bz3O457WmXLrfWJ9T6LNpQtRyREWskfAzw",
	"i22H4CgAZZDDyi7cNvaE0uR6NRtk4Pj7cpkzDiSJOe6pUiJltsJEc824OcDIx48JscZkMnqEGBkHYKNz",
	"Cwcmb0V4NvnqOkByl6tG/djoFgv+hnhYrQ1lMyKPKA0LZ3wgaNJzAOqiPer7qxNzhMMQxqfEsLkLmhs2",
	"5zS+ZpBecieKrZ1UTudefTQkzu6w5duL5VprslfRTVYTykwe6LhAtwPihdgkNq4+KvEuNgtD79HoPozy",
	"jx1Mm0b7QJGF2KDLHq8WG022B5ZhODwYgYa/YQrpFfsN3eYWmF3T7pamYlSokGScOa8mlyFxYszUAxLM",
	"ELk8DDJjbwRAx9jR1JBzyu9eJbUtnvQv8+ZWmzYVH3zgdOz4Dx2h6C4N4K9vhalzWd91JZaonaLteW6n",
	"8QYiZIzoDZvou3v6TiUFOaBSkLSEqOQ85gQ0ug3gjXPquwXGC0wWpnz7KAhnkLBiSkNjjjcXs/cv3bd5",
	"kmKNEiGWw6vTpVya9b0Xor6mbBI8dmwt895XcCE0JEsmlU7QlxFdgmn0vUKl+nvTNC4rtQMmbLkulsV5",
	"A057DtskY3kVp1c374/HZtq3NUtU1QL5LeMEaLomCywvFw2j2jG1jbTbueA3dsFv6J2td9xpME3NxNKQ",
	"S3uOP8m56HDeXewgQoAx4ujv2iBKdzBIlH2OIdexDMhAbrKHMzMNZ7usr73DlPmx9wagWCiG7yg7UnQt",
	"gcFg5yoYuomMWMJ0UJ2tn9UzcAZoWbJs07GF2lEHNWZ6LYOHL3vRwQLurhtsDwYCu2cssFiCalc4aQR8",
	"W2evlWA8G4WZs3YdkpAhhFMx5avE9hFVJx7sw9UZ0PxH2P5i2uJyJlfTye1MpzFcuxH34Ppdvb1RPKOT",
	"35rSWp6Qa6KclqUUFzRPnIF5iDSluHCkic29PfqeWV3cjHn2+ujNOwf+1XSS5kBlUosKg6vCduWfZlW2",
	"mMrAAfFVKI3O52V2K0oGm19XgAiN0pdrcBX/Amm0V5qocTgER9EZqZfxWKO9JmfnG7FL3OEjgbJ2kTTm",
	"O+shaXtF6AVlubebeWgH4oJwcePqW0W5QjjArb0rgZMsuVN20zvd8dPRUNcenhTOtaMmYWHLbioieNeF",
	"bkRINMchqRYUCwtZq0ifOfGqQEtConKWxm2sfKEMcXDrOzONCTYeEEbNiBUbcMXyigVjmWZqhKLbATKY",
	"I4pMX6RqCHcL4eqlV5z9qwLCMuDafJJ4KjsHFSs5OWt7/zo1skN/LjewtdA3w99GxgiLanVvPARit4AR",
	"eup64B7XKrNfaG2RMj8ELolrOPzDGXtX4g5nvaMPR802DHLd9riF5c37/M8Qhi2Fub+2uldeXXWvgTmi",
	"tdKZSpZS/A5xPQ/V40jWgS8jxjDK5Xfgs0jyVpfF1NadpuR7M/vgdg9JN6EVqh2kMED1uPOBWw7rGXkL",
	"NeV2q23p4lasW5xgwvjUuR2/IRgHcy+mN6eXCxor9mSEDAPTUeMAbtnStSC+s8e9M/szV9ltRgJfct2W",
	"2Xy8EmSTENTP7b+hwGCnHS0qNJIBUm0oE0yt/y9XIjJMxS8ptxWwTT97lFxvBdb4ZXpdConZtCpu9s8g",
	"ZQXN45JDlvZNvBlbMVv/uVIQFBh2A9nC+ZaKXJFm62JvUHOyJAfToIS5242MXTDFFjlgiye2xYIq5OS1",
	"IaruYpYHXK8VNn86ovm64pmETK+VRawSpBbqUL2pnVcL0JcAnBxguycvyEN02yl2AY8MFt39PDl88gKN",
	"rvaPg9gF4Aq97+ImGbKTfzh2Eqdj9FvaMQzjdqPOormh9nWOYca14zTZrmPOErZ0vG7/WSoopyuIR4oU",
	"e2CyfXE30ZDWwQvPbGl5paXYEqbj84Omhj8NxLEb9mfBIKkoCqYL59xRojD01FQPtpP64Wydelf4zcPl",
	"P6KPtPQuoo4Seb9GU3u/xVaNnuy3tIA2WqeE2hTqnDXRC74cJTnxFRqwEl5dAM/ixsxllo5iDgYzLEkp",
	"GdeoWFR6mfyVpGsqaWrY32wI3GTx7fNI9b92FSp+PcDvHe8SFMiLOOrlANl7GcL1JQ+54ElhOEr2qMkb",
	"CU7loDM37rYb8h3uHnqsUGZGSQbJrWqRGw049a0Ij+8Y8JakWK/nWvR47ZXdO2VWMk4etDI79PP7N07K",
	"KISMlV1qjruTOCRoyeACY/fim2TGvOVeyHzULtwG+i/refAiZyCW+bMcUwReioh26itS1pZ0F6sesQ4M",
	"HVPzwZDBwg01Je3qf/fv9PPG577zyXzxsOIfXWC/8JYikv0KBjYxqEwa3c6s/h74vyl5KTZjN7VzQvzG",
	"/hugJoqSiuXZL01+Z6fwq6Q8XUf9WQvT8dfmiYp6cfZ+ilY3WlPOIY8OZ2XBX73MGJFq/ynGzlMwPrJt",
	"txatXW5ncQ3gbTA9UH5Cg16mczNBiNV2wlsdUJ2vREZwnqaUTsM9+zWMg0qT/6pA6VjyEH6wQV1otzT6",
	"ri10SIBnqC3OyA/2ibk1kFalD9TSWFHltmoEZCuQzqBelbmg2ZSYcc5eH70hdlbbxxZat4UWV6iktFfR",
	"sVcFVcLGhQf7munx1IXx4+yOpTarVhoL7yhNizKWZmpanPkGmMsa2vBRfQmxMyPHVnNUXi+xkxh6WDJZ",
	"GI2rHs3KLkgT5j9a03SNKlmLpQ6T/PgKoZ4qVfAqT11dvy6dhefOwO2KhNoaoVMijN58yZR9WQwuoJ3Z",
	"Wqd5O5OAz3RtL09WnFtKicoeu8oQ3ATtHjgbqOHN/FHIOoi/pkBuC+xet2DqKfaK1qLpVl/tPcdjsxvr",
	"qun+xciUcsFZipVgYleze6VsjA9sRNGcrpHVH3F3QiOHK1rztQ6Tc1gcrALrGaFDXN8IH3w1m2qpw/6p",
	"8TmsNdVkBVo5zgbZ1JcudnZAxhW4Umj4YF3AJ4Vs+RWRQ0Zd1Unt0rgmGWFazIBi97359tap/Rgvfs44",
	"CvgObS403Vrq8BElbbQCpslKgHLraecGqw+mzwzTZDPYfJr5R5dwDOuWM8u2Puj+UEfeI+08wKbtK9PW",
	"FkVpfm5FINtJj8rSTTpc2DoqD+gNH0RwxLOYeNdOgNx6/HC0HeS2M5QE71NDaHCBjmgo8R7uEUZd5Lnz",
	"gIARWi1FYQtiQ7iitRAYj4DxhnFongSLXBBp9ErAjcHzOtBPpZJqKwKO4mlnQHP0PscYmtLO9XDboTob",
	"jCjBNfo5hrexqU89wDjqBo3gRvm2fonMUHcgTLzCJxAdIvvVplGqckJUhhkFnfrTMcZhGLevcN++APrH",
	"oC8T2e5aUntyrnMTDSWJLqpsBTqhWRarIfkSvxL8SrIKJQfYQFrVNfjKkqRYXaVdbqZPbW6iVHBVFTvm",
	"8g1uOV0qYnL0W5xA+ZSJZvAZQfZrWO/x63fvX786Ont9bO8LRVRls0SNzC2hMAxxRk640mBE50oB+S1E",
	"42/Y77fOguNgBnXnI0Qb1r73hIi5Most/hurkzdMQC5W5NrRij4wBDteW7xvj9QTzs3RSxRbJeMxgVff",
	"7dHRTH2z89j0v9MDmYtVG5B7rmCxixmHexRjw6/N/RYWeOgVf7Q3YF1/AWMDhX8tCLXbOnO4zTzxxu1V",
	"g0SfVP0ayW47yfC7IlO8owcihIO6HdSKAdbJORQnnA6GtVPtEuw0JTs55WDSkg0ysulJ9lHsqIF3KLDI",
	"xhWZz73e4wTYnjqAY+9EqI9Y6wP0ow+HJSVlzoPfMIs+Zl3g/LBVc9ehaza4uwgXjj5oWIw/7jBcQqcp",
	"m4PXQCkUawrWxl59GBkudYYPNwQlgPpj+ViFC0i1EeoDH6wEuE5BIDNZ8EbN11I6A+pHHVXmKujsKpvT",
	"L028h9n0MluC7Cxb1nU2vkjMUR1pg/5/fCVmBdw9E9OOWR8dObtcQqrZxZ5Mon8YLbXJUpl6PdY+9xYk",
	"FrE6EtM/w39N9boBaFeiz054gtJytwZnKI/gHLYPFGlRQ7TO7NTzvJvUIEAMIHdIDIkIFfNkW8Obcy4y",
	"VVMGYsFHjtju0FRzGizwH+TF3XAuT5KEhrlyO6a8EDHNfdRcpuu1MkgxqHAo2ahfYntYEDrGiuaqfnyn",
	"fmc/0GrISb/S26WrgYB5X7Wt2VdDAOV/80medpacnUP4BAFa9i+pzHyLqKrqteBkx33UyxDy5aG7QC/r",
	"mVkT59fPCYnUDsJozjQXivFVMhQS2w6tC99+xQACvA6wdjnCtQTpnmpBE3IuFCRa+LjAXXDsQoV7p/Qm",
	"SFCD9foscINVNN43ZUKwAirFqhnUBUeECzR6KzXQyaCYx/Ccu5D9yn73SRC+AuYIjdzRa7K3GoeP8GSq",
	"h8SQ6pfE3Zb7kytuovUyzu1TYypW2YMbVIbW41KKrErtBR0ejMbGMLZuzg5WElUY0/4qe7J/jlWk3gSp",
	"auewnVv5O11T3pTzah9rK0LZNQSp4Z3dvlODQFz3yVd2Aas7gfNLKtXTSSlEngyYi0/6BUq6Z+CcpeeQ",
	"EXN3+NiogSL/5CFaKWt/4OV66wtylCVwyB7NCDFqeVHqrXcNtmvtdibnD/Su+Tc4a1bZmkFO35995PGw",
	"PqzmI2/J3/wwu7maAsP8bjmVHWRP+YvNQHEUSS8jT16MfdE44qzrPkPQEJWFIial3DAXetT57uv8EdIP",
	"6vDv1n7CUglNDJa0piOUlrxBpyu8/NRYhMa9COA77AEvVIqDNwE8N3LgfOFAqZ9qpARLGaSE1vL36dn+",
	"Ie6aLwVbpDCy3izTFq6xTvb2vgRGFPWqtk3E8dw3YWBdBMGxVkzf9KHQlIglZ0PCMedSXtD8/s0XWDDj",
	"CPHhHraKLzTUf0MkW1Sqm0UrvKGj5g503bubmr9Dc8s/wOxR1AbshnJ21PotBl9CEkuj0ZzkonmTBYck",
	"lzimNRo/+ZYsXKR1KSFlinWSUC59Ncxa3cPi0M17Z7v1y33r/EXoW5CxUxBESd42lfW0wPuhgbA5ol+Y",
	"qQyc3CiVx6ivRxYR/MV4VJjyvOe6OG9Zk22l0k40h5Bwx1blwI19TatyP5l77PJwHXjpVAr66xx9W7dw",
	"G7mom7WNdYn0kbur/NoYT0a8qqLpjq4UixAsSUoQVPLbk9+IhCW+OSDI48c4wePHU9f0t6ftz+Y4P34c",
	"FePuzYnSevrdzRujmF+Gov9shNtAoGlnPyqWZ/sIoxU23Lz/gYGxv7rEgS/yAsmv1p7aP6qudvt13Lfd",
	"TUDERNbamjyYKggIHhEL7LrNoo/zK0gryfQW6xl48xv7NVon6ofaYu88PnUGrLv7tDiHuiJGY9+vlL9d",
	"fxD2Mf/CyNToPNf4GNzrDS3KHNxB+e7B4i/w7K/Ps4NnT/6y+OvBNwcpPP/mxcEBffGcPnnx7Ak8/es3",
	"zw/gyfLbF4un2dPnTxfPnz7/9psX6bPnTxbPv33xlweGDxmQLaATnz03+d/4TE9y9O4kOTPANjihJavf",
	"gDRk7F8IoCmeRCgoyyeH/qf/6U/YLBVFM7z/deKScyZrrUt1OJ9fXl7Owi7zFRr0Ei2qdD338/Tf3nt3",
	"UgdY24Rv3FEbO2tIATfVkcIRfnv/+vSMHL07mTUEMzmcHMwOZk/wZa0SOC3Z5HDyDH/C07PGfZ87Ypsc",
	"fr6aTuZroDn6v8wfBWjJUv9JXdLVCuTMPZVgfrp4OveixPyzM2Ze7fo2D6uOzj+3bL7Znp5YlXD+2Sfb",
	"727dymZ3tu6gw0godjWbLzCHZ2xTUEHj4aXYV77nn1FEHvx97hIb4h9RVbFnYO4dI/GWLSx91hsDa6eH",
	"e0R2/rl51TkAywb2z+0rZs3PvYrWK4hmGGCsP931PCnSriX7kwy5ke69torlMa2xEkn66cHBn+Ph1efX",
	"BHSnJaQVBxMB5iXNiE/2wLmf3N/cJxz9ooZDEcuBEYLn9wdBuxbpj7Alb4Um36OqcDWdfHOfO3HCjeBC",
	"c4Itg1ID/SPyMz/n4pL7lubqroqCyu3o46PpSqGpTrIL6gSnoDz15BPahm3gYvuoHWVZj+itCANKvxTZ",
	"dgfGCrUqXXBug7RGgmPcLKGvAvbf++q9jnoOW2I9Z95C6l4Hb2QrLSu4uiVP+NM+5PqVp3zlKdJO/+z+",
	"pj8FecFSIGdQlEJSyfIt+ZnXqVU35nFHWRaNSmof/b08zmjHqchgBTxxDCxZiGzry0e1JjgHq6z1BJn5",
	"53YNWCu4TTLIQUcjLszv9etc/UUstuTkuCfh2G5dzvtyi02D2qqHHz5bbceI8o0y0gWxxxnDsp5d3vQp",
	"zjV3kb1ZyEpoYrGQuUV9ZURfGdGthJvRh2eMfBPVPmziMu3d2VOfgxyrPkF1H5QxOsoXPb53svF9/Sem",
	"79joLshI8MGGIXfR/JVFfGURt2MRP0DkMOKpdUwjQnTX04fGMgwMbMm6Ly2gwd83r3IqiYKxZo4jHNEZ",
	"N+6Da9y3UhfFldXpKG8eo4ls4N3qeV9Z3leW9+dheUf7GU1bMLm1ZnQO24KWtT6k1pXOxGVg/0dYbDxO",
	"32pdv/3W+nt+SZlOlkK6XAGsRNrvrIHmc1dIofNrkxTY+4KZjsGPgYU7/uu8LvQc/dh1HcS+OtO5b9T4",
	"BkNfG/Lu2sv24ZPhu1gn0LH1xnV0OJ9jgO1aKD2fXE0/d9xK4cdP9R5/ri8Dt9dXn67+fwAAAP//ROgk",
	"pRK6AAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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