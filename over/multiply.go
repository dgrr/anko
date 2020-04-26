package over

import "reflect"

type Multiply interface {
	// Mul *
	Mul(interface{}) (interface{}, error)
	// Div /
	Div(interface{}) (interface{}, error)
	// Mod %
	Mod(interface{}) (interface{}, error)
	// Left <<
	Left(interface{}) (interface{}, error)
	// Right >>
	Right(interface{}) (interface{}, error)
	// And &
	And(interface{}) (interface{}, error)
}

var MultiplyReflectType = reflect.TypeOf(new(Multiply)).Elem()

type MultiplyImpl struct{}

func (a *MultiplyImpl) Mul(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *MultiplyImpl) Div(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *MultiplyImpl) Mod(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *MultiplyImpl) Left(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *MultiplyImpl) Right(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}

func (a *MultiplyImpl) And(_ interface{}) (interface{}, error) {
	return nil, ErrMethodNotImplemented
}
