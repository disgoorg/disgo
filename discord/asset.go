package discord

import (
	"strings"

	"github.com/DisgoOrg/disgo/rest/route"
)

var DefaultCDNConfig = CDNConfig{
	Size:   0,
	Format: route.PNG,
}

type CDNConfig struct {
	Size   int
	Format route.ImageFormat
	V      int
}

// Apply applies the given ConfigOpt(s) to the Config
func (c *CDNConfig) Apply(opts []CDNOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

type CDNOpt func(config *CDNConfig)

//goland:noinspection GoUnusedExportedFunction
func WithSize(size int) CDNOpt {
	return func(config *CDNConfig) {
		config.Size = size
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithFormat(format route.ImageFormat) CDNOpt {
	return func(config *CDNConfig) {
		config.Format = format
	}
}

//goland:noinspection GoUnusedExportedFunction
func WithV(v int) CDNOpt {
	return func(config *CDNConfig) {
		config.V = v
	}
}

func formatAssetURL(cdnRoute *route.CDNRoute, opts []CDNOpt, params ...any) *string {
	var lastStringParam string
	lastParam := params[len(params)-1]
	if str, ok := lastParam.(string); ok {
		if str == "" {
			return nil
		}
		lastStringParam = str
	} else if ptrStr, ok := lastParam.(*string); ok {
		if ptrStr == nil {
			return nil
		}
		lastStringParam = *ptrStr
	}

	config := &DefaultCDNConfig
	config.Apply(opts)

	if strings.HasPrefix(lastStringParam, "a_") && !config.Format.CanBeAnimated() {
		config.Format = route.GIF
	}

	queryValues := route.QueryValues{}
	if config.V > 0 {
		queryValues["v"] = config.V
	}

	compiledRoute, err := cdnRoute.Compile(queryValues, config.Format, config.Size, params...)
	if err != nil {
		return nil
	}
	url := compiledRoute.URL()
	return &url
}
