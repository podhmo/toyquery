package internaltoyquery

import (
	"context"
	"encoding/json"

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
