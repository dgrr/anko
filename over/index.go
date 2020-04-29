package over

import (
	"errors"
	"reflect"
)

var (
	ErrOutOfRange       = errors.New("index out of range")
	ErrUnsupportedIndex = errors.New("bad index type")
)

type Index interface {
	// Index []
	Index(interface{}) (interface{}, error)
	Len
}

var IndexReflectType = reflect.TypeOf(new(Index)).Elem()

type IndexImpl struct{}

func (s *IndexImpl) Index(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
