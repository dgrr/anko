package over

import "reflect"

type Assign interface {
	// Assign =
	Assign(interface{}) error
}

var AssignReflectType = reflect.TypeOf(new(Assign)).Elem()

type AssignImpl struct{}

func (a *AssignImpl) Assign(_ interface{}) error {
	return ErrMethodNotImplemented
}
