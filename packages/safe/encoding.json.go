package safe

import (
	"encoding/json"
	"reflect"

	"github.com/dgrr/pako/env"
)

func init() {
	env.Packages["encoding/json"] = map[string]reflect.Value{
		"Marshal":   reflect.ValueOf(json.Marshal),
		"Unmarshal": reflect.ValueOf(json.Unmarshal),
	}
}
