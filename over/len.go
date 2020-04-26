package over

import "reflect"

type Len interface {
	// Len len()
	Len() int64
}

var LenReflectType = reflect.TypeOf(new(Len)).Elem()

type LenImpl struct{}

func (l *LenImpl) Len() int64 {
	return 0
}
