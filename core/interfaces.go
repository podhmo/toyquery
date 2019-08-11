package core

import (
	"context"
	"io"
)

// TODO: rename

// Client :
type Client interface {
	io.Closer
	Must(func() error)

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
	MustExec(ctx context.Context, code string)
}

// Table :
type Table interface {
	// TODO: add methods (tentative)
	FindByID(ctx context.Context, id ID, val interface{}) error
	InsertByID(ctx context.Context, id ID, val interface{}) error
	Count(ctx context.Context) (int, error)
}

// ID :
type ID = string
