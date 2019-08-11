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

// Table :
func (s *Session) Table(name string) (toyquerycore.Table, error) {
	u := s.Client.Universe
	if u == nil {
		return nil, toyquerycore.ErrTableNotFound
	}

	w, ok := u.Worlds[name]
	if !ok {
		if !s.Client.Config.AutoCreate {
			return nil, toyquerycore.ErrTableNotFound.New(name)
		}
		w = u.NewWorld(name)
	}
	return w, nil
}
