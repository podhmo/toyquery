package toysql

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/podhmo/toyquery/core"
)

// ID :
type ID = string // TODO : int ?

// Config :
type Config struct {
	DriverName string
	URI        string
}

// Client :
type Client struct {
	Config *Config

	mu sync.Mutex
	DB *sqlx.DB
}

// Close :
func (c *Client) Close() error {
	return nil
}

// Session :
func (c *Client) Session(ctx context.Context) (core.Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.DB == nil {
		db, err := sqlx.Open(c.Config.DriverName, c.Config.URI) // xxx
		if err != nil {
			return nil, core.ErrDatabaseNotFound.Wrap(err, "open")
		}
		c.DB = db
	}
	return &Session{
		Client: c,
	}, nil
}

// Connect :
func Connect(ctx context.Context, driver, uri string, options ...func(*Config)) *Client {
	c := &Config{
		DriverName: driver,
		URI:        uri,
	}
	for _, opt := range options {
		opt(c)
	}
	return &Client{
		Config: c,
	}
}

// Session :
type Session struct {
	Client *Client
}

// Close :
func (s *Session) Close() error {
	// remove from client?
	return nil
}

// Table :
func (s *Session) Table(ctx context.Context, name string) (core.Table, error) {
	tableName := strings.SplitN(strings.TrimSpace(name), " ", 2)[0]
	return &Table{session: s, Name: tableName}, nil
}

// MustExec
func (s *Session) MustExec(ctx context.Context, code string) {
	db := s.Client.DB
	db.MustExecContext(ctx, code)
}

// Table :
type Table struct {
	session *Session
	Name    string
}

// Count :
func (t *Table) Count(ctx context.Context) (int, error) {
	var c int

	// XXX: SQLInjection (t.Name is stripped, but)
	stmt := fmt.Sprintf(
		`SELECT count(*) FROM %s`,
		t.Name,
	)

	db := t.session.Client.DB
	if err := db.GetContext(ctx, &c, stmt); err != nil {
		return 0, core.ErrRecordNotFound.Wrap(err, t.Name)
	}
	return c, nil
}

// InsertByID :
func (t *Table) InsertByID(ctx context.Context, id ID, v interface{}) error {
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	m := t.session.Client.DB.Mapper
	fields := m.TypeMap(rt)

	names := make([]string, len(fields.Index))
	values := make([]interface{}, len(fields.Index))
	holders := make([]string, len(fields.Index))

	for i, f := range fields.Index {
		names[i] = f.Name
		values[i] = reflectx.FieldByIndexesReadOnly(rv, f.Index).Interface()
		holders[i] = "?"
	}

	// INSERT INTO <table name> (<column name>...) VALUES (?...);
	stmt := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		t.Name,
		strings.Join(names, ", "),
		strings.Join(holders, ", "),
	)

	db := t.session.Client.DB
	result, err := db.ExecContext(ctx, stmt, values...)
	_ = result // TODO: saved in session?
	return err
}

// FindByID :
func (t *Table) FindByID(ctx context.Context, id ID, dst interface{}) error {
	// SELECT (<column name>...) FROM <table name> WHERE <id>=?
	// TODO : id field
	return t.Find(ctx, dst, core.Where("id=?", id))
}

// Find :
func (t *Table) Find(ctx context.Context, dst interface{}, options ...func(*core.QOption)) error {
	q := &core.QOption{}
	for _, op := range options {
		op(q)
	}
	if len(q.Wheres) == 0 {
		return core.ErrEmptyCondition.New(t.Name)
	}

	var vals []interface{}
	var fmts []string
	for _, opt := range q.Wheres {
		fmts = append(fmts, opt.Fmt)
		vals = append(vals, opt.Vals...)
	}

	// XXX: SQLInjection (t.Name is stripped, but)
	stmt := fmt.Sprintf(
		`SELECT * FROM %s WHERE %s`,
		t.Name, strings.Join(fmts, " AND "),
	)

	db := t.session.Client.DB
	if err := db.GetContext(ctx, dst, stmt, vals...); err != nil {
		return core.ErrRecordNotFound.Wrap(err, t.Name)
	}
	return nil
}
