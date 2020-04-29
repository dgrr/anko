package over

import (
	"fmt"
	"reflect"
)

type String interface {
	// String()
	fmt.Stringer
}

var StringReflectType = reflect.TypeOf(new(String)).Elem()

type StringImpl struct{}

func (a *StringImpl) String() string {
	return ""
}
