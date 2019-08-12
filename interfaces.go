package toyquery

import (
	"context"
	"io"
)

// TODO: rename

// Client :
type Client interface {
	io.Closer
	SessionFactory
	// Ping() error
}

// SessionFactory :
type SessionFactory interface {
	Session(ctx context.Context) (Session, error)
}

// Session :
type Session interface {
	io.Closer
	Table(ctx context.Context, name string) (Table, error)
	Exec(ctx context.Context, code string) error
}

// Table :
type Table interface {
	// TODO: add methods (tentative)
	GetByID(ctx context.Context, id ID, val interface{}) error
	Get(ctx context.Context, val interface{}, options ...func(*QOption)) error
	InsertByID(ctx context.Context, id ID, val interface{}) error
	Count(ctx context.Context, options ...func(*QOption)) (int, error)
}

// ID :
type ID = string
