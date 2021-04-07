package endpoints

import log "github.com/sirupsen/logrus"

// FileExtension is the type of an image on Discord's CDN
type FileExtension string

// The available FileExtension(s)
const (
	PNG  FileExtension = "png"
	JPEG FileExtension = "jpg"
	WEBP FileExtension = "webp"
	GIF  FileExtension = "gif"
	BLANK  FileExtension = ""
)

func (f FileExtension) String() string {
	return string(f)
}

// CDNRoute is a route for interacting with images hosted on discord's CDN
type CDNRoute struct {
	Route
	supportedFileExtensions []FileExtension
}

// NewCDNRoute generates a new discord cdn route struct
func NewCDNRoute(url string, supportedFileExtensions ...FileExtension) CDNRoute {
	return CDNRoute{
		Route: Route{
			baseRoute:  CDN,
			route:      url,
			paramCount: countParams(url),
		},
		supportedFileExtensions: supportedFileExtensions,
	}
}

// Compile builds a full request URL based on provided arguments
func (r CDNRoute) Compile(fileExtension FileExtension, args ...interface{}) CompiledCDNRoute {
	supported := false
	for _, supportedFileExtension := range r.supportedFileExtensions {
		if supportedFileExtension == fileExtension {
			supported = true
		}
	}
	if !supported {
		log.Infof("provided file extension: %s is not supported by discord on this endpoint!", fileExtension)
	}
	compiledRoute := CompiledCDNRoute{
		CompiledRoute: r.Route.Compile(args...),
	}
	compiledRoute.CompiledRoute.route += fileExtension.String()
	return compiledRoute
}

// CompiledCDNRoute is CDNRoute compiled with all URL args
type CompiledCDNRoute struct {
	CompiledRoute
}
