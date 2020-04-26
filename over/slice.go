package over

import "reflect"

// TODO: Implement
type Slice interface {
	// Index []
	Index(interface{}) (interface{}, error)
}

var SliceReflectType = reflect.TypeOf(new(Slice)).Elem()

type SliceImpl struct{}

func (s *SliceImpl) Index(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
