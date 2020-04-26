package over

import "reflect"

// TODO: Add support for non-pointers?

type Mem interface {
	// New new
	New() (interface{}, error)
}

var MemReflectType = reflect.TypeOf(new(Mem)).Elem()

type MemImpl struct{}

func (m *MemImpl) New() (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
