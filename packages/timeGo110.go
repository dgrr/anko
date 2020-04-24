// +build go1.10

package packages

import (
	"reflect"
	"time"

	"github.com/dgrr/anko/env"
)

func timeGo110() {
	env.Packages["time"]["LoadLocationFromTZData"] = reflect.ValueOf(time.LoadLocationFromTZData)
}
