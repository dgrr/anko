package over

import "reflect"

type Set interface {
	// Set =
	Set(interface{}) error
}

var SetReflectType = reflect.TypeOf(new(Set)).Elem()

type SetImpl struct{}

func (a *SetImpl) Set(_ interface{}) error {
	return ErrMethodNotImplemented
}
