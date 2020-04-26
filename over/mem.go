package over

import "reflect"

// TODO: Implement
type Mem interface {
	// Make make()
	Make() (interface{}, error)
	// New new
	New() (interface{}, error)
}

var MemReflectType = reflect.TypeOf(new(Mem)).Elem()

type MemImpl struct{}

func (m *MemImpl) Make() (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (m *MemImpl) New() (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
