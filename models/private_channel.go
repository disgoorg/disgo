package models

type PrivateChannel struct {
	ID            Snowflake `json:"id"`
	LastMessageID Snowflake `json:"last_message_id"`
	Type          int    `json:"type"`
	Users    []User `json:"recipients"`
}

