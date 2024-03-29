package toyinmemory

import (
	"context"
	"fmt"
	"os"

	"github.com/podhmo/toyquery"
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

// Session :
func (c *Client) Session(ctx context.Context) (toyquery.Session, error) {
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
		Config: c,
		Universe: &Universe{
			Worlds: map[string]*World{},
		},
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
func (s *Session) Table(ctx context.Context, name string) (toyquery.Table, error) {
	u := s.Client.Universe
	if u == nil {
		return nil, toyquery.ErrTableNotFound
	}

	w, ok := u.Worlds[name]
	if !ok {
		if !s.Client.Config.AutoCreate {
			return nil, toyquery.ErrTableNotFound.New(name)
		}
		w = u.NewWorld(name)
	}
	return w, nil
}

// Exec
func (s *Session) Exec(ctx context.Context, code string) error {
	fmt.Fprintln(os.Stderr, "*****not supported *********************")
	fmt.Fprintln(os.Stderr, code)
	fmt.Fprintln(os.Stderr, "****************************************")
	return nil
}
