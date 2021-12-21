package route

import (
	"fmt"
	"net/url"
	"strings"
)

// NewCDNRoute generates a new discord cdn path struct
func NewCDNRoute(path string, supportedImageFormats ...ImageFormat) *CDNRoute {
	queryParams := []string{"size", "v"}

	params := map[string]struct{}{}
	for _, param := range queryParams {
		params[param] = struct{}{}
	}

	return &CDNRoute{
		basePath:              CDN,
		path:                  path,
		queryParams:           params,
		urlParamCount:         countURLParams(path),
		supportedImageFormats: supportedImageFormats,
	}
}

// NewCustomCDNRoute generates a new custom cdn path struct
//goland:noinspection GoUnusedExportedFunction
func NewCustomCDNRoute(basePath string, path string, supportedImageFormats ...ImageFormat) *CDNRoute {
	route := NewCDNRoute(path, supportedImageFormats...)
	route.basePath = basePath
	return route
}

// CDNRoute is a path for interacting with images hosted on discord's CDN
type CDNRoute struct {
	basePath              string
	path                  string
	queryParams           map[string]struct{}
	urlParamCount         int
	supportedImageFormats []ImageFormat
}

// Compile builds a full request URL based on provided arguments
func (r *CDNRoute) Compile(queryValues QueryValues, imageFormat ImageFormat, size int, params ...interface{}) (*CompiledCDNRoute, error) {
	supported := false
	for _, supportedFileExtension := range r.supportedImageFormats {
		if supportedFileExtension == imageFormat {
			supported = true
		}
	}
	if !supported {
		return nil, ErrImageFormatNotSupported(imageFormat)
	}
	if queryValues == nil {
		queryValues = QueryValues{}
	}
	queryValues["size"] = size

	path := r.path
	for _, param := range params {
		start := strings.Index(path, "{")
		end := strings.Index(path, "}")
		paramValue := fmt.Sprint(param)
		path = path[:start] + paramValue + path[end+1:]
	}

	if imageFormat.String() != "" {
		path += "." + imageFormat.String()
	}

	queryParamsStr := ""
	if queryValues != nil {
		query := url.Values{}
		for param, value := range queryValues {
			if _, ok := r.queryParams[param]; !ok {
				return nil, ErrUnexpectedQueryParam(param)
			}
			query.Add(param, fmt.Sprint(value))
		}
		if len(query) > 0 {
			queryParamsStr = query.Encode()
		}
	}

	return &CompiledCDNRoute{
		CDNRoute:    r,
		path:        path,
		queryParams: queryParamsStr,
	}, nil
}

// Path returns the request path used by the path
func (r *CDNRoute) Path() string {
	return r.path
}

// CompiledCDNRoute is CDNRoute compiled with all URL args
type CompiledCDNRoute struct {
	CDNRoute    *CDNRoute
	path        string
	queryParams string
}

// URL returns the full URL for the resource
func (r *CompiledCDNRoute) URL() string {
	u := r.CDNRoute.basePath + r.path
	if r.queryParams != "" {
		u += "?" + r.queryParams
	}
	return u
}
