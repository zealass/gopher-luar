package luar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/yuin/gopher-lua"
)

type luaStructWrapper struct {
	L      *lua.LState
	Struct interface{}
}

func (w *luaStructWrapper) Index(key lua.LValue) (lua.LValue, error) {
	ref := reflect.ValueOf(w.Struct)
	refType := ref.Type()

	// Check for method
	keyLString, ok := key.(lua.LString)
	if ok {
		if method, ok := refType.MethodByName(string(keyLString)); ok {
			return New(w.L, method.Func.Interface()), nil
		}
	}

	// Check for field
	if field := ref.FieldByName(string(keyLString)); field.IsValid() {
		if !field.CanInterface() {
			return nil, errors.New("cannot interface field " + string(keyLString))
		}
		if val := New(w.L, field.Interface()); val != nil {
			return val, nil
		}
	}

	return lua.LNil, nil
}

func (w *luaStructWrapper) NewIndex(key lua.LValue, value lua.LValue) error {
	ref := reflect.ValueOf(w.Struct)

	keyLString, ok := key.(lua.LString)
	if !ok {
		return errors.New("invalid non-string key")
	}

	keyString := string(keyLString)
	field := ref.FieldByName(keyString)
	if !field.IsValid() {
		return errors.New("unknown field " + keyString)
	}
	if !field.CanSet() {
		return errors.New("cannot set field " + keyString)
	}
	field.Set(lValueToReflect(value, field.Type()))
	return nil
}

func (w *luaStructWrapper) Len() (lua.LValue, error) {
	return nil, errors.New("cannot # struct")
}

func (w *luaStructWrapper) Call(...lua.LValue) ([]lua.LValue, error) {
	return nil, errors.New("cannot call struct")
}

func (w *luaStructWrapper) String() (string, error) {
	if stringer, ok := w.Struct.(fmt.Stringer); ok {
		return stringer.String(), nil
	}
	return reflect.ValueOf(w.Struct).String(), nil
}

func (w *luaStructWrapper) Equals(v luaWrapper) (lua.LValue, error) {
	return nil, errors.New("cannot compare struct")
}

func (w *luaStructWrapper) Unwrap() interface{} {
	return w.Struct
}
