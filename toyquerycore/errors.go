package toyquerycore

import "fmt"

var (
	// ErrDatabaseNotFound :
	ErrDatabaseNotFound = fmt.Errorf("database not found")
	// ErrTableNotFound :
	ErrTableNotFound = fmt.Errorf("table not found")
	// ErrFieldNotFound :
	ErrFieldNotFound = fmt.Errorf("field not found")
	// ErrRecordNotFound :
	ErrRecordNotFound = fmt.Errorf("record not found")
)
