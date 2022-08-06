package rest

import (
	"errors"

	"github.com/disgoorg/snowflake/v2"
)

var ErrNoMorePages = errors.New("no more pages")

type Page[T any] struct {
	getItems  func(before snowflake.ID, after snowflake.ID, limit int) ([]T, error)
	getIDFunc func(t T) snowflake.ID

	Items []T
	Err   error

	Before snowflake.ID
	After  snowflake.ID
	Limit  int
}

func (p *Page[T]) Next() bool {
	if p.Err != nil {
		return false
	}
	if len(p.Items) != p.Limit {
		p.Err = ErrNoMorePages
		return false
	}

	if len(p.Items) > 0 {
		p.After = p.getIDFunc(p.Items[0])
	}

	p.Items, p.Err = p.getItems(0, p.After, p.Limit)
	return p.Err == nil
}

func (p *Page[T]) Previous() bool {
	if p.Err != nil {
		return false
	}
	if len(p.Items) != p.Limit {
		p.Err = ErrNoMorePages
		return false
	}

	if len(p.Items) > 0 {
		p.Before = p.getIDFunc(p.Items[len(p.Items)-1])
	}

	p.Items, p.Err = p.getItems(p.Before, 0, p.Limit)
	return p.Err == nil
}
