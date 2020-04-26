package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"github.com/dgrr/pako/core"
	"github.com/dgrr/pako/env"
	"github.com/dgrr/pako/over"
	"github.com/dgrr/pako/parser"
	"github.com/dgrr/pako/vm"

	_ "github.com/dgrr/pako/packages"
)

type Something struct {
	over.IndexImpl
	over.AddImpl
	over.MultiplyImpl
	over.AssignImpl
	over.ComparisonImpl
	over.BinaryImpl
	over.MemImpl
	n int
}

func (s *Something) Index(v interface{}) (interface{}, error) {
	switch n := v.(type) {
	case int:
		if n == 0 {
			return s.n, nil
		}
	case int64:
		if n == 0 {
			return s.n, nil
		}
	}

	return nil, over.ErrOutOfRange
}

func (s *Something) Assign(v interface{}) error {
	switch n := v.(type) {
	case int:
		s.n = n
	case int64:
		s.n = int(n)
	case *Something:
		s.n = s.n
	default:
		return over.ErrNoMatch
	}

	return nil
}

func (s *Something) Add(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n + int(n)
	case *Something:
		r = s.n + n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Sub(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n - int(n)
	case *Something:
		r = s.n - n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Or(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n | int(n)
	case *Something:
		r = s.n | n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Mul(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n * int(n)
	case *Something:
		r = s.n * n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Div(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n / int(n)
	case *Something:
		r = s.n / n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Mod(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n % int(n)
	case *Something:
		r = s.n % n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Left(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n << int(n)
	case *Something:
		r = s.n << n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Right(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n >> int(n)
	case *Something:
		r = s.n >> n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) New() (interface{}, error) {
	fmt.Println("NEW")
	return &Something{}, nil
}

func (s *Something) And(v interface{}) (r interface{}, err error) {
	switch n := v.(type) {
	case int64:
		r = s.n & int(n)
	case *Something:
		r = s.n & n.n
	default:
		err = over.ErrNoMatch
	}

	return
}

func (s *Something) Equals(v interface{}) error {
	switch n := v.(type) {
	case int64:
		if s.n == int(n) {
			return nil
		}
	case *Something:
		if s.n == n.n {
			return nil
		}
	}
	return over.ErrNoMatch
}

func (s *Something) Binary() bool {
	return s.n != 0
}

func (s *Something) String() string {
	return strconv.FormatInt(int64(s.n), 10)
}

func main() {
	e := env.NewEnv()
	e.Define("S", &Something{n: 20})
	e.DefineType("Some", reflect.TypeOf(Something{}))
	core.Import(e)

	d, err := ioutil.ReadFile("./override.pak")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = vm.Execute(e, nil, string(d))
	if err != nil {
		switch e := err.(type) {
		case *parser.Error:
			log.Fatalf("%s:%d:%d %s\n", e.Filename, e.Pos.Line, e.Pos.Column, e.Message)
		case *vm.Error:
			log.Fatalf("%s at %d\n", e.Message, e.Pos.Line)
		default:
			log.Fatalf("Execute error: %v\n", err)
		}
	}
}
