package suite

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery/core"
)

// Env :
type Env struct {
	Connect func() core.Client
	Setup   func(core.Session)
}

// Simple
func Simple(t *testing.T, ctx context.Context, env *Env) {
	type dummy struct {
		ID   core.ID `db:"id"`
		Name string  `db:"name"`
	}

	c := env.Connect()
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

	if env.Setup != nil {
		t.Run("setup", func(t *testing.T) {
			env.Setup(s)
		})
	}

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
