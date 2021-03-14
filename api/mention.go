package api

// Mentionable is a struct for Mention parsing and AllowedMentions
type Mentionable interface {
	Mention() string
}
