package toyinmemory

import (
	"fmt"
	"reflect"

	"github.com/google/shlex"
	"github.com/jmoiron/sqlx/reflectx"
)

// Op :: { '=' '>' '<' '>=' '<=' '<>' }
// Val :: { '?' <int> <string> <float> }
// support `Val Op Val` only

type Expr struct {
	Op    string
	Name  string
	Value interface{}
}

// Parse
func Parse(s string, val interface{}) (Expr, error) {
	tokens, err := shlex.Split(s)
	if err != nil {
		return Expr{}, err
	}
	if len(tokens) != 3 {
		return Expr{}, fmt.Errorf("invalid input %q", s)
	}

	supported := false
	for _, x := range []string{"=", ">", "<", ">=", "<=", "<>"} {
		if x == tokens[1] {
			supported = true
			break
		}
	}
	if !supported {
		return Expr{}, fmt.Errorf("unsupported operator %q", tokens[1])
	}

	if tokens[0] == "?" {
		// ? <op> val
		var op string
		switch tokens[1] {
		case "=":
			op = "="
		case "<>":
			op = "<>"
		case ">":
			op = "<"
		case ">=":
			op = "<="
		case "<":
			op = ">"
		case "<=":
			op = ">="
		}
		return Expr{Value: val, Op: op, Name: tokens[2]}, nil
	} else if tokens[2] == "?" {
		// val <op> ?
		return Expr{Name: tokens[0], Op: tokens[1], Value: val}, nil
	}
	return Expr{}, fmt.Errorf("invalid input %q", s)
}

// Eval :
func Eval(expr Expr, env interface{}) (bool, error) {
	switch expr.Op {
	case "=":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return reflect.DeepEqual(x, y), nil })
	case "<>":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return !reflect.DeepEqual(x, y), nil })
	case ">":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return GreaterThan(x, y) })
	case "<":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return GreaterThan(y, x) })
	case ">=":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return GreaterThanEqual(x, y) })
	case "<=":
		return eval(expr, env, func(x, y interface{}) (bool, error) { return GreaterThanEqual(y, x) })
	default:
		return false, fmt.Errorf("unsupported operator %q", expr.Op)
	}
}

// xxx : see db tag ?
var m = reflectx.NewMapperFunc("db", func(s string) string { return s })

// eval :
func eval(expr Expr, env interface{}, check func(x, y interface{}) (bool, error)) (bool, error) {
	switch env := env.(type) {
	case map[string]interface{}:
		x := env[expr.Name]
		y := expr.Value
		return check(x, y)
	case Object:
		x := env[expr.Name]
		y := expr.Value
		return check(x, y)
	case *Object:
		x := (*env)[expr.Name]
		y := expr.Value
		return check(x, y)
	default:
		x := m.FieldByName(reflect.ValueOf(env), expr.Name).Interface()
		y := expr.Value
		return check(x, y)
	}
}

// GreaterThan is x > y
func GreaterThan(x interface{}, y interface{}) (bool, error) {
	rx := reflect.ValueOf(x)
	ry := reflect.ValueOf(y)
	// if reflect.TypeOf(x).Kind() != reflect.TypeOf(y).Kind() {
	// 	return false, fmt.Errorf("mismatch %s <-> %s", rx.Kind(), ry.Kind())
	// }
	switch rx.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rx.Int() > ry.Convert(rx.Type()).Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rx.Uint() > ry.Convert(rx.Type()).Uint(), nil
	case reflect.Float32, reflect.Float64:
		return rx.Float() > ry.Convert(rx.Type()).Float(), nil
	case reflect.String:
		return rx.String() > ry.Convert(rx.Type()).String(), nil
	default:
		return false, fmt.Errorf("unexpected %s", rx.Kind())
	}
}

// GreaterThanEqual is x >= y
func GreaterThanEqual(x interface{}, y interface{}) (bool, error) {
	rx := reflect.ValueOf(x)
	ry := reflect.ValueOf(y)

	switch rx.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rx.Int() >= ry.Convert(rx.Type()).Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rx.Uint() >= ry.Convert(rx.Type()).Uint(), nil
	case reflect.Float32, reflect.Float64:
		return rx.Float() >= ry.Convert(rx.Type()).Float(), nil
	case reflect.String:
		return rx.String() >= ry.Convert(rx.Type()).String(), nil
	case reflect.Bool:
		return rx.Bool() == ry.Convert(rx.Type()).Bool(), nil
	default:
		return reflect.DeepEqual(rx, ry), nil
	}
}
