package over

import "reflect"

type Binary interface {
	// Binary && or || or if
	Binary() bool
}

var BinaryReflectType = reflect.TypeOf(new(Binary)).Elem()

type BinaryImpl struct{}

func (b *BinaryImpl) Binary() bool {
	return false
}
