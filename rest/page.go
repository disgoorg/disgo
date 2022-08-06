package rest

import "github.com/disgoorg/snowflake/v2"

type Page[T any] struct {
	Before    snowflake.ID
	After     snowflake.ID
	Limit     int
	Data      []T
	Err       error
	f         func(before snowflake.ID, after snowflake.ID, limit int) ([]T, error)
	getIDFunc func(t T) snowflake.ID
}

func (p *Page[T]) Next() bool {
	if p.Err != nil {
		return false
	}

	if len(p.Data) > 0 {
		p.After = p.getIDFunc(p.Data[0])
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

	if len(p.Data) > 0 {
		p.Before = p.getIDFunc(p.Data[len(p.Data)-1])
	}

	data, err := p.f(p.Before, 0, p.Limit)
	p.Data = data
	p.Err = err
	return err == nil
}
