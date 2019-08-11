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
	Table(name string) (Table, error)
}

// Table :
type Table interface {
	// TODO: add methods (tentative)
	FindByID(id ID, val interface{}) error
	InsertByID(id ID, val interface{}) error
	Count() (int, error)
}

// ID :
type ID = string
