package resolver

import (
	"errors"
	"reflect"
)

// Reflect implements Object
// It uses reflection to determine the value to resolve to
type Reflect struct {
	Target interface{}
}

// Resolve resolves a method from the Target
//
// The following priority list is executed top to bottom
// The method returns on the first matching line, regardless of error status
//
// * Method r.GetMyField
// * Method r.MyField
// * Field MyField
//
// If no match is found, return an error
func (r Reflect) Resolve(field string, args Args) (result interface{}, err error) {
	val := reflect.ValueOf(r.Target)
	typ := val.Type()

	if _, ok := typ.MethodByName("Get" + field); ok {
		m := val.MethodByName("Get" + field)
		return call(m, args)
	}
	if _, ok := typ.MethodByName(field); ok {
		m := val.MethodByName(field)
		return call(m, args)
	}
	if _, ok := typ.FieldByName(field); ok {
		return val.FieldByName(field).Interface(), nil
	}

	return nil, errors.New("missing field")
}
func call(m reflect.Value, args Args) (result interface{}, err error) {
	switch m.Type().NumIn() {
	case 0:
		return call0(m)
	case 1:
		return call1(m, args)
	default:
		return nil, errors.New("unexpected input args")
	}
}
func call0(m reflect.Value) (result interface{}, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New("panic while calling method")
		}
	}()

	return coerceOutput(m.Call([]reflect.Value{}))
}
func call1(m reflect.Value, args Args) (result interface{}, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New("panic while calling method")
		}
	}()

	return coerceOutput(m.Call([]reflect.Value{reflect.ValueOf(args)}))
}
func coerceOutput(out []reflect.Value) (result interface{}, err error) {
	if len(out) == 1 {
		return out[0].Interface(), nil
	}

	if len(out) == 2 {
		result = out[0].Interface()
		e := out[1].Interface()
		if e != nil {
			err = e.(error)
		}
		return
	}

	return nil, errors.New("wrong # of values in output")
}
