package internaltoyquery_test

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/internaltoyquery"
)

type dummy struct {
	ID   internaltoyquery.ID
	Name string
}

func TetstIt(t *testing.T) {
	ctx := context.Background()
	c := internaltoyquery.Connect(ctx)

	defer c.Must(c.Close)
	s, err := c.Session("db")
	noerror.Must(t, err)
	defer c.Must(c.Close)

	dummies := []dummy{
		{ID: "1", Name: "foo"},
		{ID: "2", Name: "bar"},
	}

	dbname := "person"

	t.Run("insert", func(t *testing.T) {
		{
			ob := dummies[0]
			noerror.Must(t, s.InsertByID(dbname, ob.ID, &ob))
		}
		{
			ob := dummies[1]
			noerror.Must(t, s.InsertByID(dbname, ob.ID, &ob))
		}
	})

	t.Run("count", func(t *testing.T) {
		noerror.Should(t,
			noerror.Equal(2).ActualWithError(s.Count(dbname)),
		)
	})
}
