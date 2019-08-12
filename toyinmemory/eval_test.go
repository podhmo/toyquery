package toyinmemory

import (
	"fmt"
	"testing"

	"github.com/podhmo/noerror"
)

func TestEval(t *testing.T) {
	type Value struct {
		ID    string
		Value int    `db:"v"`
	}

	type C struct {
		expr    string
		literal interface{}
		value   Value
		ok      bool
	}

	cases := []C{
		{expr: "v = ?", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "v = ?", literal: 10, value: Value{Value: 11}, ok: false},
		{expr: "v = ?", literal: 11, value: Value{Value: 10}, ok: false},
		{expr: "? = v", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "? = v", literal: 10, value: Value{Value: 11}, ok: false},
		{expr: "? = v", literal: 11, value: Value{Value: 10}, ok: false},

		{expr: "v <> ?", literal: 10, value: Value{Value: 10}, ok: !true},
		{expr: "v <> ?", literal: 10, value: Value{Value: 11}, ok: !false},
		{expr: "v <> ?", literal: 11, value: Value{Value: 10}, ok: !false},
		{expr: "? <> v", literal: 10, value: Value{Value: 10}, ok: !true},
		{expr: "? <> v", literal: 10, value: Value{Value: 11}, ok: !false},
		{expr: "? <> v", literal: 11, value: Value{Value: 10}, ok: !false},

		{expr: "v > ?", literal: 10, value: Value{Value: 9}, ok: false},
		{expr: "v > ?", literal: 10, value: Value{Value: 10}, ok: false},
		{expr: "v > ?", literal: 10, value: Value{Value: 11}, ok: true},
		{expr: "? < v", literal: 10, value: Value{Value: 9}, ok: false},
		{expr: "? < v", literal: 10, value: Value{Value: 10}, ok: false},
		{expr: "? < v", literal: 10, value: Value{Value: 11}, ok: true},

		{expr: "v < ?", literal: 10, value: Value{Value: 9}, ok: true},
		{expr: "v < ?", literal: 10, value: Value{Value: 10}, ok: false},
		{expr: "v < ?", literal: 10, value: Value{Value: 11}, ok: false},
		{expr: "? > v", literal: 10, value: Value{Value: 9}, ok: true},
		{expr: "? > v", literal: 10, value: Value{Value: 10}, ok: false},
		{expr: "? > v", literal: 10, value: Value{Value: 11}, ok: false},

		{expr: "v >= ?", literal: 10, value: Value{Value: 9}, ok: false},
		{expr: "v >= ?", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "v >= ?", literal: 10, value: Value{Value: 11}, ok: true},
		{expr: "? <= v", literal: 10, value: Value{Value: 9}, ok: false},
		{expr: "? <= v", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "? <= v", literal: 10, value: Value{Value: 11}, ok: true},

		{expr: "v <= ?", literal: 10, value: Value{Value: 9}, ok: true},
		{expr: "v <= ?", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "v <= ?", literal: 10, value: Value{Value: 11}, ok: false},
		{expr: "? >= v", literal: 10, value: Value{Value: 9}, ok: true},
		{expr: "? >= v", literal: 10, value: Value{Value: 10}, ok: true},
		{expr: "? >= v", literal: 10, value: Value{Value: 11}, ok: false},

		{expr: "? = ID", literal: "x", value: Value{ID: "x"}, ok: true},
		{expr: "? = ID", literal: "x", value: Value{ID: "y"}, ok: false},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			expr, err := Parse(c.expr, c.literal)
			noerror.Must(t, err)
			noerror.Should(t,
				noerror.Equal(c.ok).ActualWithError(Eval(expr, c.value)),
				fmt.Sprintf("%s with %v in %+v", c.expr, c.literal, c.value),
			)
		})
	}
}
