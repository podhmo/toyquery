package toyinmemoryquery_test

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/core"
	"github.com/podhmo/toyquery/toyinmemoryquery"
)

type dummy struct {
	ID   toyinmemoryquery.ID
	Name string
}

func TestIt(t *testing.T) {
	ctx := context.Background()

	var c core.Client
	defer noerror.Bind(t, &c).Actual(toyinmemoryquery.Connect(ctx)).Teardown()

	var s core.Session
	defer noerror.Bind(t, &s).ActualWithError(c.Session(ctx)).Teardown()

	dummies := []dummy{
		{ID: "1", Name: "foo"},
		{ID: "2", Name: "bar"},
	}

	var table core.Table
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