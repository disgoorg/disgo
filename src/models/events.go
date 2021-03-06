package models

type Event struct {
	ResponseNumber int
}

type ReadyEvent struct {
	Event
	User User
}
