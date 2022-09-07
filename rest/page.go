package rest

import (
	"errors"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var ErrNoMorePages = errors.New("no more pages")

type Page[T any] struct {
	getItemsFunc func(before snowflake.ID, after snowflake.ID) ([]T, error)
	getIDFunc    func(t T) snowflake.ID

	Items []T
	Err   error

	ID snowflake.ID
}

func (p *Page[T]) Next() bool {
	if p.Err != nil {
		return false
	}

	if len(p.Items) > 0 {
		p.ID = p.getIDFunc(p.Items[0])
	}

	p.Items, p.Err = p.getItemsFunc(0, p.ID)
	if p.Err == nil && len(p.Items) == 0 {
		p.Err = ErrNoMorePages
	}
	return p.Err == nil
}

func (p *Page[T]) Previous() bool {
	if p.Err != nil {
		return false
	}

	if len(p.Items) > 0 {
		p.ID = p.getIDFunc(p.Items[len(p.Items)-1])
	}

	p.Items, p.Err = p.getItemsFunc(p.ID, 0)
	if p.Err == nil && len(p.Items) == 0 {
		p.Err = ErrNoMorePages
	}
	return p.Err == nil
}

type AuditLogPage struct {
	getItems func(before snowflake.ID) (discord.AuditLog, error)

	discord.AuditLog
	Err error

	ID snowflake.ID
}

func (p *AuditLogPage) Previous() bool {
	if p.Err != nil {
		return false
	}

	if len(p.AuditLogEntries) > 0 {
		p.ID = p.AuditLogEntries[len(p.AuditLogEntries)-1].ID
	}

	p.AuditLog, p.Err = p.getItems(p.ID)
	if p.Err == nil && len(p.AuditLogEntries) == 0 {
		p.Err = ErrNoMorePages
	}
	return p.Err == nil
}
