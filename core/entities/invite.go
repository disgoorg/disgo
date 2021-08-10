package entities

import (
	"github.com/DisgoOrg/restclient"
)

// URL returns the invite URL in format like https://discord.gg/{code}
func (i *Invite) URL() string {
	url, err := restclient.InviteURL.Compile(nil, i.Code)
	if err != nil {
		return ""
	}
	return url.Route()
}
