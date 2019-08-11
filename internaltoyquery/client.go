package internaltoyquery

import (
	"context"

	"github.com/podhmo/toyquery/toyquerycore"
)

// Config :
type Config struct {
	Name       string
	PrimaryKey string // default primarykey TODO: implementation
	AutoCreate bool
}

// WithName :
func WithName(name string) func(*Config) {
	return func(c *Config) {
		c.Name = name
	}
}

// Client :
type Client struct {
	Config   *Config
	Universe *Universe
}

// Close :
func (c *Client) Close() error {
	return nil
}

// Must :
func (c *Client) Must(run func() error) {
	if err := run(); err != nil {
		panic(err)
	}
}

// Session :
func (c *Client) Session(uri string) (toyquerycore.Session, error) {
	return &Session{
		Client: c,
	}, nil
}

// Connect :
func Connect(ctx context.Context, options ...func(*Config)) *Client {
	c := &Config{
		AutoCreate: true,
	}
	for _, opt := range options {
		opt(c)
	}
	return &Client{
		Config:   c,
		Universe: &Universe{},
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

// Count :
func (s *Session) Count(name string) (int, error) {
	if s.Client.Universe == nil {
		return 0, toyquerycore.ErrTableNotFound
	}

	w, ok := s.Client.Universe.Worlds[name]
	if !ok {
		return 0, toyquerycore.ErrTableNotFound.New(name)
	}
	return w.Count()
}

// FindByID :
func (s *Session) FindByID(name string, id ID, val interface{}) error {
	if s.Client.Universe == nil {
		return toyquerycore.ErrTableNotFound
	}

	w, ok := s.Client.Universe.Worlds[name]
	if !ok {
		return toyquerycore.ErrTableNotFound.New(name)
	}

	return w.FindByID(id, val)
}

func (s *Session) InsertByID(name string, id ID, val interface{}) error {
	u := s.Client.Universe
	if u == nil {
		return toyquerycore.ErrTableNotFound
	}

	w, ok := u.Worlds[name]
	if !ok {
		if !s.Client.Config.AutoCreate {
			return toyquerycore.ErrTableNotFound.New(name)
		}
		w = u.NewWorld(name)
	}
	return w.InsertByID(id, val)
}
