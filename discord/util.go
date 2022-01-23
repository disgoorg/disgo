package discord

import (
	"strings"

	"github.com/DisgoOrg/disgo/rest/route"
	"github.com/DisgoOrg/snowflake"
)

func FormatAssetURL(cdnRoute *route.CDNRoute, entityId snowflake.Snowflake, assetId *string, size int) *string {
	if assetId == nil {
		return nil
	}
	format := route.PNG
	if strings.HasPrefix(*assetId, "a_") {
		format = route.GIF
	}
	compiledRoute, err := cdnRoute.Compile(nil, format, size, entityId, *assetId)
	if err != nil {
		return nil
	}
	url := compiledRoute.URL()
	return &url
}
