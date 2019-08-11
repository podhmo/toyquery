package toyinmemoryquery_test

import (
	"context"
	"testing"

	"github.com/podhmo/toyquery/core"
	"github.com/podhmo/toyquery/suite"
	"github.com/podhmo/toyquery/toyinmemoryquery"
)

func TestIt(t *testing.T) {
	ctx := context.Background()
	env := &suite.Env{
		Connect: func() core.Client {
			return toyinmemoryquery.Connect(ctx)
		},
	}
	t.Run("simple", func(t *testing.T) {
		suite.Simple(t, ctx, env)
	})
	t.Run("shortcut", func(t *testing.T) {
		suite.Shortcut(t, env)
	})
}
