package over

import "reflect"

type Add interface {
	// Add + or ++
	Add(interface{}) (interface{}, error)
	// Sub - or --
	Sub(interface{}) (interface{}, error)
	// Or |
	Or(interface{}) (interface{}, error)
}

var AddReflectType = reflect.TypeOf(new(Add)).Elem()

type AddImpl struct{}

func (a *AddImpl) Add(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *AddImpl) Sub(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *AddImpl) Or(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
