// +build go1.8

package packages

import (
	"reflect"
	"time"

	"github.com/dgrr/pako/env"
)

func timeGo18() {
	env.Packages["time"]["Until"] = reflect.ValueOf(time.Until)
}
