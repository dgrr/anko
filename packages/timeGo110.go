// +build go1.10

package packages

import (
	"reflect"
	"time"

	"github.com/dgrr/pako/env"
)

func timeGo110() {
	env.Packages["time"]["LoadLocationFromTZData"] = reflect.ValueOf(time.LoadLocationFromTZData)
}
