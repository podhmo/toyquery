package toysqlquery_test

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/toysqlquery"
)

type dummy struct {
	ID   toysqlquery.ID `db:"id"`
	Name string         `db:"name"`
}

func TestIt(t *testing.T) {
	ctx := context.Background()

	c := toysqlquery.Connect(ctx, "sqlite3", ":memory:")
	defer noerror.Must(t, c.Close())

	s, err := c.Session(ctx)
	noerror.Must(t, err)
	defer noerror.Must(t, s.Close())

	dummies := []dummy{
		{ID: "1", Name: "foo"},
		{ID: "2", Name: "bar"},
	}

	table, err := s.Table(ctx, "person")
	noerror.Must(t, err)

	t.Run("create table", func(t *testing.T) {
		schema := `CREATE TABLE person (
        id text primary key,
        name text);`
		s.MustExec(ctx, schema)
	})

	t.Run("insert", func(t *testing.T) {
		{
			ob := dummies[0]
			noerror.Must(t, table.InsertByID(ctx, ob.ID, &ob))
		}
		{
			ob := dummies[1]
			noerror.Must(t, table.InsertByID(ctx, ob.ID, &ob))
		}
	})

	t.Run("count", func(t *testing.T) {
		noerror.Should(t,
			noerror.Equal(2).ActualWithError(table.Count(ctx)),
		)
	})
}
