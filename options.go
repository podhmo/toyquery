package toyquery

// QOption :
type QOption struct {
	Wheres []QOptionItem
}

// QOptionItem :
type QOptionItem struct {
	Fmt  string
	Vals []interface{}
}

// Where :
func Where(fmt string, vals ...interface{}) func(*QOption) {
	return func(opt *QOption) {
		// assertion ?
		opt.Wheres = append(opt.Wheres, QOptionItem{Fmt: fmt, Vals: vals})
	}
}
