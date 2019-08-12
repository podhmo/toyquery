package suite

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
	"github.com/podhmo/toyquery"
	"github.com/podhmo/toyquery/shortcut"
)

// Env :
type Env struct {
	Connect func() toyquery.Client
	Setup   func(toyquery.Session)
}

// Simple
func Simple(t *testing.T, ctx context.Context, env *Env) {
	type dummy struct {
		ID    toyquery.ID `db:"id" json:"id"`
		Name  string      `db:"name" json:"name"`
		Value int         `db:"value" json:"value"`
	}

	c := env.Connect()
	defer noerror.Must(t, c.Close())

	s, err := c.Session(ctx)
	noerror.Must(t, err)
	defer noerror.Must(t, s.Close())

	dummies := []dummy{
		{ID: "1", Name: "foo", Value: 10},
		{ID: "2", Name: "bar", Value: 100},
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

	t.Run("find by id", func(t *testing.T) {
		var got dummy
		noerror.Must(t, table.GetByID(ctx, dummies[1].ID, &got))
		noerror.Should(t, noerror.DeepEqual(dummies[1]).Actual(got))
	})

	t.Run("find", func(t *testing.T) {
		var got dummy
		noerror.Must(t, table.Get(ctx, &got, toyquery.Where("id = ?", dummies[0].ID)))
		noerror.Should(t, noerror.DeepEqual(dummies[0]).Actual(got))

		noerror.Must(t, table.Get(ctx, &got, toyquery.Where("? <> id", dummies[0].ID)))
		noerror.Should(t, noerror.DeepEqual(dummies[1]).Actual(got))

		noerror.Must(t, table.Get(ctx, &got, toyquery.Where("? < value", dummies[0].Value)))
		noerror.Should(t, noerror.DeepEqual(dummies[1]).Actual(got))

		noerror.Must(t, table.Get(ctx, &got, toyquery.Where("value > ?", dummies[0].Value)))
		noerror.Should(t, noerror.DeepEqual(dummies[1]).Actual(got))

		noerror.Must(t, table.Get(ctx, &got, toyquery.Where("name = ?", dummies[0].Name)))
		noerror.Should(t, noerror.DeepEqual(dummies[0]).Actual(got))
	})
}

// Shortcut
func Shortcut(t *testing.T, env *Env) {
	type dummy struct {
		ID   toyquery.ID `db:"id"`
		Name string      `db:"name"`
	}

	c, teardown := shortcut.Create(t, env.Connect())
	defer teardown()
	s := c.Session()
	defer noerror.Must(t, s.Raw().Close()) // TODO: AddTeardown

	dummies := []dummy{
		{ID: "1", Name: "foo"},
		{ID: "2", Name: "bar"},
	}

	table := s.Table("person")

	if env.Setup != nil {
		t.Run("setup", func(t *testing.T) {
			env.Setup(s.Raw())
		})
	}

	t.Run("insert", func(t *testing.T) {
		{
			ob := dummies[0]
			table.InsertByID(t, ob.ID, &ob)
		}
		{
			ob := dummies[1]
			table.InsertByID(t, ob.ID, &ob)
		}
	})

	t.Run("count", func(t *testing.T) {
		noerror.Should(t,
			noerror.Equal(2).Actual(table.Count(t)),
		)
	})

	t.Run("find one", func(t *testing.T) {
		var got dummy
		table.GetByID(t, dummies[1].ID, &got)
		noerror.Should(t, noerror.DeepEqual(dummies[1]).Actual(got))
	})
}
