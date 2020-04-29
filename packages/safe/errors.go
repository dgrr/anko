package safe

import (
	"errors"
	"reflect"

	"github.com/dgrr/pako/env"
)

func init() {
	env.Packages["errors"] = map[string]reflect.Value{
		"New": reflect.ValueOf(errors.New),
	}
}
