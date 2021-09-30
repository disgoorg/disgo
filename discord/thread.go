package discord

type ThreadCreateWithMessage struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
}

type ThreadCreate struct {
	Name                string              `json:"name"`
	AutoArchiveDuration AutoArchiveDuration `json:"auto_archive_duration"`
	Type                ChannelType         `json:"type"`
	Invitable           bool                `json:"invitable"`
}

type GetThreads struct {
	Threads []Channel `json:"threads"`
	Members []Member  `json:"members"`
	HasMore bool      `json:"has_more"`
}
