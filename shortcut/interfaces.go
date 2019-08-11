package shortcut

import (
	"testing"

	"github.com/podhmo/toyquery/core"
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
	Session(t *testing.T) Session
}

// Session :
type Session interface {
	Table(t *testing.T, name string) Table
	MustExec(t *testing.T, code string) // ?
	Session() core.Session
}

// Table :
type Table interface {
	// TODO: add methods (tentative)
	FindByID(t *testing.T, id ID, val interface{})
	InsertByID(t *testing.T, id ID, val interface{})
	Count(t *testing.T) int
}
