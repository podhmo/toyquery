package toyinmemory

import (
	"context"
	"testing"

	"github.com/podhmo/noerror"
)

func TestQuery(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ctx := context.Background()
	u := newUniverse()
	w := u.NewWorld("person")

	ps := []Person{
		{
			Name: "foo",
			Age:  20,
		},
		{
			Name: "bar",
			Age:  20,
		},
	}

	t.Run("insert", func(t *testing.T) {
		noerror.Must(t, w.InsertByID(ctx, ps[0].Name, ps[0]))
		noerror.Must(t, w.InsertByID(ctx, ps[1].Name, ps[1]))

		noerror.Should(t, noerror.Equal(2).ActualWithError(w.Count(ctx)))
	})

	t.Run("find", func(t *testing.T) {
		for _, p := range ps {
			p := p
			t.Run(p.Name, func(t *testing.T) {
				var got Person
				noerror.Must(t, w.GetByID(ctx, p.Name, &got))
				noerror.Should(t, noerror.DeepEqual(p).Actual(got))
			})
		}
	})
}
