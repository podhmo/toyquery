package toyquerycore

import "io"

// TODO: rename

// Client :
type Client interface {
	io.Closer
	Must(func() error)

	SessionFactory
	Ping() error
}

// SessionFactory :
type SessionFactory interface {
	Session(uri string) (Session, error)
}

// Session :
type Session interface {
	io.Closer

	// TODO: add methods
	// tenatative
	FindByID(name string, id ID, val interface{}) error
	InsertByID(name string, id ID, val interface{}) error
	Count(name string) (int, error)
}

// ID :
type ID = string
