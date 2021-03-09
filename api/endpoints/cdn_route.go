package endpoints

import log "github.com/sirupsen/logrus"

type FileExtension string

const (
	PNG  FileExtension = "png"
	JPEG FileExtension = "jpg"
	WEBP FileExtension = "webp"
	GIF  FileExtension = "gif"
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
func (r CDNRoute) Compile(fileExtension FileExtension, args ...string) string {
	supported := false
	for _, supportedFileExtension := range r.supportedFileExtensions {
		if supportedFileExtension == fileExtension {
			supported = true
		}
	}
	if !supported {
		log.Infof("provided file extension: %s is not supported by discord on this endpoint!", fileExtension)
	}
	return r.Route.Compile(args...) + fileExtension.String()
}
