package toysql_test

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/podhmo/toyquery"
	"github.com/podhmo/toyquery/suite"
	"github.com/podhmo/toyquery/toysql"
)

func TestIt(t *testing.T) {
	ctx := context.Background()
	env := &suite.Env{
		Connect: func() toyquery.Client {
			return toysql.Connect(ctx, "sqlite3", ":memory:")
		},
		Setup: func(s toyquery.Session) {
			schema := `CREATE TABLE person (
        id text primary key not null,
        name text not null,
        value integer not null);`
			s.MustExec(ctx, schema)
		},
	}

	// todo: drop db
	t.Run("simple", func(t *testing.T) {
		suite.Simple(t, ctx, env)
	})
}
