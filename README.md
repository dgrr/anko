# Pako

[![GoDoc Reference](https://godoc.org/github.com/dgrr/pako/vm?status.svg)](http://godoc.org/github.com/dgrr/pako/vm)
[![Build Status](https://travis-ci.org/dgrr/pako.svg?branch=master)](https://travis-ci.org/dgrr/pako)
[![Go Report Card](https://goreportcard.com/badge/github.com/dgrr/pako)](https://goreportcard.com/report/github.com/dgrr/pako)

Pako is a scriptable interpreter written in Go.

# BEWARE! THIS IS A FORK OF THE ORIGINAL ANKO. THIS WILL BE A DIFFERENT LANGUAGE.

# Pako features

- Import directives with local imports. Controlled by [Env.Import](https://godoc.org/github.com/dgrr/pako/env#ImportFrom)
- Load function controlled by [Env.LoadFrom](https://godoc.org/github.com/dgrr/pako/env#LoadFrom)
- `?` operator to ignore errors.
- Strict lvalue and rvalue assignment.
- Value unpacking/ignoring. Using ()
- [Operator override](https://github.com/dgrr/pako/tree/master/_example/programs/override)

## Usage Example - Embedded

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/dgrr/pako/env"
	"github.com/dgrr/pako/vm"
)

func main() {
	env.Packages["my_pkg"] = map[string]reflect.Value{
		"print": reflect.ValueOf(fmt.Println),
		"getStr": reflect.ValueOf(func(n int64) (string, error) {
			if n < 0 {
				return "", errors.New("Cannot represent negative numbers")
			}
			return strconv.FormatInt(n, 10), nil
		}),
	}
	e := env.NewEnv()

	script := `import my_pkg as pkg
             str = pkg.getStr(20)?
             pkg.print("My string representation:", str)
             pkg.print("Now should panic!!!")
             pkg.getStr(-1)?`

	_, err := vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("Execute error: %v\n", err)
	}
	// output:
	// My string representation: 20
	// Now should panic!!!
}
```

More examples are located in the GoDoc:

https://godoc.org/github.com/dgrr/pako/vm


## Usage Example - Command Line

### Building
```
go get github.com/dgrr/pako
go install github.com/dgrr/pako
```

### Running an Pako script file named script.ank
```
./pako script.pak
```

## Pako Script Quick Start
```
// declare variables
x = 1
y = x + 1

// print using outside the script defined println function
println(x + y) // 3

// if else statement
if x < 1 || y < 1 {
	println(x)
} else if x < 1 && y < 1 {
	println(y)
} else {
	println(x + y)
}

// slice
a = []interface{1, 2, 3}
println(a) // [1 2 3]
println(a[0]) // 1

// map
a = map[interface]interface{"x": 1}
println(a) // map[x:1]
a.b = 2
a["c"] = 3
println(a["b"]) // 2
println(a.c) // 3

// struct
a = make(struct {
	A int64,
	B float64
})
a.A = 4
a.B = 5.5
println(a.A) // 4
println(a.B) // 5.5

// function
fn a (x) {
	println(x + 1)
}
a(5) // 6
```


## Please note that the master branch is not stable

The master branch language and API may change at any time.

To mitigate breaking changes, please use tagged branches. New tagged branches will be created for breaking changes.


## Original author

Yasuhiro Matsumoto (a.k.a mattn)

## Contributors

### Code Contributors

This project exists thanks to all the people who contribute. [[Contribute](CONTRIBUTING.md)].
<a href="https://github.com/dgrr/pako/graphs/contributors"><img src="https://opencollective.com/mattn-anko/contributors.svg?width=890&button=false" /></a>

### Financial Contributors

Become a financial contributor and help us sustain our community. [[Contribute](https://opencollective.com/mattn-anko/contribute)]

#### Individuals

<a href="https://opencollective.com/mattn-anko"><img src="https://opencollective.com/mattn-anko/individuals.svg?width=890"></a>

#### Organizations

Support this project with your organization. Your logo will show up here with a link to your website. [[Contribute](https://opencollective.com/mattn-anko/contribute)]

<a href="https://opencollective.com/mattn-anko/organization/0/website"><img src="https://opencollective.com/mattn-anko/organization/0/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/1/website"><img src="https://opencollective.com/mattn-anko/organization/1/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/2/website"><img src="https://opencollective.com/mattn-anko/organization/2/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/3/website"><img src="https://opencollective.com/mattn-anko/organization/3/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/4/website"><img src="https://opencollective.com/mattn-anko/organization/4/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/5/website"><img src="https://opencollective.com/mattn-anko/organization/5/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/6/website"><img src="https://opencollective.com/mattn-anko/organization/6/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/7/website"><img src="https://opencollective.com/mattn-anko/organization/7/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/8/website"><img src="https://opencollective.com/mattn-anko/organization/8/avatar.svg"></a>
<a href="https://opencollective.com/mattn-anko/organization/9/website"><img src="https://opencollective.com/mattn-anko/organization/9/avatar.svg"></a>
