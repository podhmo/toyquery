package shortcut

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/core"
)

// Create :
func Create(t *testing.T, c core.Client) (Client, func()) {
	return &client{Client: c}, func() { t.Helper(); noerror.Must(t, c.Close()) }
}

type client struct {
	Client core.Client
}

func (c *client) Session(t *testing.T) Session {
	t.Helper()
	ctx := context.Background()
	s, err := c.Client.Session(ctx)
	noerror.Must(t, err)
	return &session{session: s}
}

type session struct {
	session core.Session
}

// Session :
func (s *session) Session() core.Session {
	return s.session
}

// MustExec :
func (s *session) MustExec(t *testing.T, code string) {
	t.Helper()
	ctx := context.Background()
	s.session.MustExec(ctx, code)
}

// Table :
func (s *session) Table(t *testing.T, name string) Table {
	t.Helper()
	ctx := context.Background()
	tbl, err := s.session.Table(ctx, name)
	noerror.Must(t, err)
	return &table{Table: tbl}
}

type table struct {
	Table core.Table
}

// Count :
func (tbl *table) Count(t *testing.T) int {
	t.Helper()
	ctx := context.Background()
	c, err := tbl.Table.Count(ctx)
	noerror.Must(t, err)
	return c
}

// InsertByID :
func (tbl *table) InsertByID(t *testing.T, id ID, v interface{}) {
	t.Helper()
	ctx := context.Background()
	noerror.Must(t, tbl.Table.InsertByID(ctx, id, v))
}

// FindByID :
func (tbl *table) FindByID(t *testing.T, id ID, v interface{}) {
	t.Helper()
	ctx := context.Background()
	noerror.Must(t, tbl.Table.FindByID(ctx, id, v))
}
