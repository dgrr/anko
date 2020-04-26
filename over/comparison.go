package over

import "reflect"

type Comparison interface {
	// Equals ==
	Equals(interface{}) error
	// Distinct !=
	Distinct(interface{}) error
	// Less <
	Less(interface{}) error
	// LessEquals <=
	LessEquals(interface{}) error
	// Greater >
	Greater(interface{}) error
	// GreaterEquals >=
	GreaterEquals(interface{}) error
}

var ComparisonReflectType = reflect.TypeOf(new(Comparison)).Elem()

type ComparisonImpl struct{}

func (c *ComparisonImpl) Equals(_ interface{}) error {
	return ErrMethodNotImplemented
}

func (c *ComparisonImpl) Distinct(_ interface{}) error {
	return ErrMethodNotImplemented
}

func (c *ComparisonImpl) Less(_ interface{}) error {
	return ErrMethodNotImplemented
}

func (c *ComparisonImpl) LessEquals(_ interface{}) error {
	return ErrMethodNotImplemented
}

func (c *ComparisonImpl) Greater(_ interface{}) error {
	return ErrMethodNotImplemented
}

func (c *ComparisonImpl) GreaterEquals(_ interface{}) error {
	return ErrMethodNotImplemented
}
