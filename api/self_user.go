package api

import (
	"errors"

	"github.com/DisgoOrg/restclient"
)

// ErrDMChannelToYourself occurs when opening a DMChannel to yourself
var ErrDMChannelToYourself = restclient.NewRestError(nil, errors.New("can't open a dm channel to yourself"))

// SelfUser represents the current logged in User
type SelfUser struct {
	*User
}

// Update updates the SelfUser with the given payload
func (u *SelfUser) Update(updateSelfUser UpdateSelfUser) (*SelfUser, restclient.RestError) {
	return u.Disgo.RestClient().UpdateSelfUser(updateSelfUser)
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *SelfUser) OpenDMChannel() (*DMChannel, restclient.RestError) {
	return nil, ErrDMChannelToYourself
}

// UpdateSelfUser is the payload used to update the SelfUser
type UpdateSelfUser struct {
	Username string `json:"username"`
	Avatar   interface{}
}
