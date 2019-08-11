package toyquerycore

import "io"

// TODO: rename

// Client :
type Client interface {
	io.Closer
	SessionFactory
	Ping() error
}

// SessionFactory :
type SessionFactory interface {
	Session(uri string) (Session, error)
}

// Session :
type Session interface {
	// TODO: add methods
	io.Closer
}
