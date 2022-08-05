package rest

import "github.com/disgoorg/snowflake/v2"

type Page[T any] struct {
	Before     snowflake.ID
	After      snowflake.ID
	Limit      int
	Data       []T
	Err        error
	f          func(before snowflake.ID, after snowflake.ID, limit int) ([]T, error)
	lastIDFunc func(data []T) snowflake.ID
}

func (p *Page[T]) Next() bool {
	if p.Err != nil {
		return false
	}

	if p.Data != nil {
		p.After = p.lastIDFunc(p.Data)
	}

	data, err := p.f(0, p.After, p.Limit)
	p.Data = data
	p.Err = err
	return err == nil
}

func (p *Page[T]) Previous() bool {
	if p.Err != nil {
		return false
	}

	if p.Data != nil {
		p.Before = p.lastIDFunc(p.Data)
	}

	data, err := p.f(p.Before, 0, p.Limit)
	p.Data = data
	p.Err = err
	return err == nil
}
