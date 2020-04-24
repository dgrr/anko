// +build !appengine

package packages

import (
	"os"
	"reflect"

	"github.com/dgrr/anko/env"
)

func osNotAppEngine() {
	env.Packages["os"]["Getppid"] = reflect.ValueOf(os.Getppid)
}
