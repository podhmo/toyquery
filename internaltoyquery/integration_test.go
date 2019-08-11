package internaltoyquery_test

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/internaltoyquery"
	"github.com/podhmo/toyquery/toyquerycore"
)

type dummy struct {
	ID   internaltoyquery.ID
	Name string
}

func TetstIt(t *testing.T) {
	ctx := context.Background()

	var c toyquerycore.Client
	defer noerror.Bind(t, &c).Actual(internaltoyquery.Connect(ctx)).Teardown()

	var s toyquerycore.Session
	defer noerror.Bind(t, &s).ActualWithError(c.Session(ctx, "db")).Teardown()

	dummies := []dummy{
		{ID: "1", Name: "foo"},
		{ID: "2", Name: "bar"},
	}

	var table toyquerycore.Table
	defer noerror.Bind(t, &table).ActualWithError(s.Table(ctx, "person")).Teardown()

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
