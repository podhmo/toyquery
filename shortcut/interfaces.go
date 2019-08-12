package shortcut

import (
	"testing"

	"github.com/podhmo/toyquery"
)

// ID :
type ID = string

// Client :
type Client interface {
	SessionFactory
	// Ping() error
}

// SessionFactory :
type SessionFactory interface {
	Session() Session
}

// Session :
type Session interface {
	Table(name string) Table
	MustExec(t *testing.T, code string) // ?
	Raw() toyquery.Session
}

// Table :
type Table interface {
	// TODO: add methods (tentative)
	FindByID(t *testing.T, id ID, val interface{})
	InsertByID(t *testing.T, id ID, val interface{})
	Count(t *testing.T) int
}
