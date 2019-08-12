package toyinmemory

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/podhmo/toyquery/core"
)

// Not goroutine safe
//
// universe, world, object = database, table, record (document)

// ID :
type ID = string

func newUniverse() *Universe {
	return &Universe{
		Worlds: map[string]*World{},
	}
}

// Universe :
type Universe struct {
	Name   string
	Worlds map[string]*World
}

// NewWorld :
func (u *Universe) NewWorld(name string) *World {
	// TODO: conflict check?
	w := &World{
		Name:    name,
		Objects: map[ID]*Object{},
	}
	u.Worlds[name] = w
	return w
}

// World :
type World struct {
	Name    string
	Objects map[ID]*Object
}

// Describe :
func (w *World) Describe() string {
	return w.Name
}

// Count :
func (w *World) Count(ctx context.Context) (int, error) {
	return len(w.Objects), nil
}

// InsertByID :
func (w *World) InsertByID(ctx context.Context, id ID, src interface{}) error {
	var ob Object
	if err := Materialize(&ob, src); err != nil {
		return core.ErrSomethingWrong.Wrap(err, w.Describe())
	}
	w.Objects[id] = &ob
	return nil
}

// FindByID :
func (w *World) FindByID(ctx context.Context, id ID, dst interface{}) error {
	if ob, ok := w.Objects[id]; ok {
		return Unmaterialize(dst, ob)
	}
	return core.ErrRecordNotFound.New(w.Describe())
}

// Find :
func (w *World) Find(ctx context.Context, dst interface{}, options ...func(*core.QOption)) error {
	q := &core.QOption{}
	for _, op := range options {
		op(q)
	}
	if len(q.Wheres) == 0 {
		return core.ErrEmptyCondition.New(w.Describe())
	}

	exprs := make([]Expr, 0, len(q.Wheres))
	for _, opt := range q.Wheres {
		expr, err := Parse(opt.Fmt, opt.Vals[0])
		if err != nil {
			return core.ErrInvalidCondition.Wrap(err, fmt.Sprintf("%s with %v", opt.Fmt, opt.Vals[0]))
		}
		exprs = append(exprs, expr)
	}

	for _, ob := range w.Objects {
		ok := true
		for _, expr := range exprs {
			matched, err := Eval(expr, ob)
			if err != nil {
				return core.ErrSomethingWrong.Wrap(err, fmt.Sprintf("%v with %v", expr, ob))
			}
			if !matched {
				ok = matched
				break
			}
		}
		if ok {
			return Unmaterialize(dst, ob)
		}
	}
	return core.ErrRecordNotFound.New(w.Describe())

}

// Object :
type Object map[string]interface{}

// Unmaterialize :
func Unmaterialize(dst interface{}, src *Object) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return err
	}
	return nil
}

// Materialize :
func Materialize(dst *Object, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return err
	}
	return nil
}
