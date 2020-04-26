// +build !appengine

package packages

import (
	"os"
	"reflect"

	"github.com/dgrr/pako/env"
)

func osNotAppEngine() {
	env.Packages["os"]["Getppid"] = reflect.ValueOf(os.Getppid)
}
