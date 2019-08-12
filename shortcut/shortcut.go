package shortcut

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery"
)

// Create :
func Create(t *testing.T, c toyquery.Client) (Client, func()) {
	return &client{t: t, Client: c}, func() { t.Helper(); noerror.Must(t, c.Close()) }
}

type client struct {
	t      *testing.T
	Client toyquery.Client
}

func (c *client) Session() Session {
	c.t.Helper()
	ctx := context.Background()
	s, err := c.Client.Session(ctx)
	noerror.Must(c.t, err)
	return &session{t: c.t, session: s}
}

type session struct {
	t       *testing.T
	session toyquery.Session
}

// Raw :
func (s *session) Raw() toyquery.Session {
	return s.session
}

// MustExec :
func (s *session) MustExec(t *testing.T, code string) {
	t.Helper()
	ctx := context.Background()
	noerror.Must(t, s.session.Exec(ctx, code))
}

// Table :
func (s *session) Table(name string) Table {
	s.t.Helper()
	ctx := context.Background()
	tbl, err := s.session.Table(ctx, name)
	noerror.Must(s.t, err)
	return &table{Table: tbl}
}

type table struct {
	Table toyquery.Table
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

// GetByID :
func (tbl *table) GetByID(t *testing.T, id ID, v interface{}) {
	t.Helper()
	ctx := context.Background()
	noerror.Must(t, tbl.Table.GetByID(ctx, id, v))
}
