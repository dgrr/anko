package env

import (
	"fmt"
	"reflect"
	"strings"
)

// DefineType defines type in current scope.
func (e *Env) DefineType(symbol string, aType interface{}) error {
	var reflectType reflect.Type
	if aType == nil {
		reflectType = NilType
	} else {
		var ok bool
		reflectType, ok = aType.(reflect.Type)
		if !ok {
			reflectType = reflect.TypeOf(aType)
		}
	}

	return e.DefineReflectType(symbol, reflectType)
}

// DefineReflectType defines type in current scope.
func (e *Env) DefineReflectType(symbol string, reflectType reflect.Type) error {
	if strings.Contains(symbol, ".") {
		return ErrSymbolContainsDot
	}

	e.rwMutex.Lock()
	if e.types == nil {
		e.types = make(map[string]reflect.Type)
	}
	e.types[symbol] = reflectType
	e.rwMutex.Unlock()

	return nil
}

// DefineGlobalType defines type in global scope.
func (e *Env) DefineGlobalType(symbol string, aType interface{}) error {
	for e.parent != nil {
		e = e.parent
	}
	return e.DefineType(symbol, aType)
}

// DefineGlobalReflectType defines type in global scope.
func (e *Env) DefineGlobalReflectType(symbol string, reflectType reflect.Type) error {
	for e.parent != nil {
		e = e.parent
	}
	return e.DefineReflectType(symbol, reflectType)
}

// Type returns reflect type from the scope where symbol is frist found.
func (e *Env) Type(symbol string) (reflect.Type, error) {
	e.rwMutex.RLock()
	reflectType, ok := e.types[symbol]
	e.rwMutex.RUnlock()
	if ok {
		return reflectType, nil
	}

	if e.externalLookup != nil {
		var err error
		reflectType, err = e.externalLookup.Type(symbol)
		if err == nil {
			return reflectType, nil
		}
	}

	if e.parent == nil {
		reflectType, ok = basicTypes[symbol]
		if ok {
			return reflectType, nil
		}
		return NilType, fmt.Errorf("undefined type '%s'", symbol)
	}

	return e.parent.Type(symbol)
}

// Method is used to define methods in structs.
//
// The symbol should be the name of the struct + "." + function name.
func (e *Env) Method(symbol string) (reflect.Value, error) {
	e.rwMutex.RLock()
	v, ok := e.methods[symbol]
	e.rwMutex.RUnlock()
	if ok {
		return v, nil
	}

	if e.parent == nil {
		return reflect.Value{}, fmt.Errorf("undefined method '%s'", symbol)
	}

	return e.parent.Method(symbol)
}

// DefineMethod will define a method in a struct.
func (e *Env) DefineMethod(symbol string, v reflect.Value) {
	e.rwMutex.Lock()
	if e.methods == nil {
		e.methods = make(map[string]reflect.Value)
	}
	e.methods[symbol] = v
	e.rwMutex.Unlock()
}
