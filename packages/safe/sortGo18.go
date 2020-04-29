// +build go1.8

package safe

import (
	"reflect"
	"sort"

	"github.com/dgrr/pako/env"
)

func sortGo18() {
	env.Packages["sort"]["Slice"] = reflect.ValueOf(sort.Slice)
	env.Packages["sort"]["SliceIsSorted"] = reflect.ValueOf(sort.SliceIsSorted)
	env.Packages["sort"]["SliceStable"] = reflect.ValueOf(sort.SliceStable)
}
