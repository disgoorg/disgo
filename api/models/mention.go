package models

type Mentionable interface {
	Mention() string
}
