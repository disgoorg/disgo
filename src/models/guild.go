package models

type Guild struct {
	ID      Snowflake
	Name    string
	Icon    *string
	OwnerID Snowflake
}
