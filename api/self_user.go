package api

import "errors"

// ErrDMChannelToYourself occurs when opening a DMChannel to yourself
var ErrDMChannelToYourself = errors.New("can't open a dm channel to yourself")

// SelfUser represents the current logged in User
type SelfUser struct {
	*User
}

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *SelfUser) OpenDMChannel() (*DMChannel, error) {
	return nil, ErrDMChannelToYourself
}
