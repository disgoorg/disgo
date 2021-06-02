package api

import "errors"

var ErrDMChannelToYourself = errors.New("can't open a dm channel to yourself")

type SelfUser User

// OpenDMChannel creates a DMChannel between the user and the Disgo client
func (u *SelfUser) OpenDMChannel() (*DMChannel, error) {
	return nil, ErrDMChannelToYourself
}
