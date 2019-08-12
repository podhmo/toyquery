package toysql_test

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/podhmo/toyquery/core"
	"github.com/podhmo/toyquery/suite"
	"github.com/podhmo/toyquery/toysql"
)

func TestIt(t *testing.T) {
	ctx := context.Background()
	env := &suite.Env{
		Connect: func() core.Client {
			return toysql.Connect(ctx, "sqlite3", ":memory:")
		},
		Setup: func(s core.Session) {
			schema := `CREATE TABLE person (
        id text primary key,
        name text);`
			s.MustExec(ctx, schema)
		},
	}

	suite.Simple(t, ctx, env)
}
