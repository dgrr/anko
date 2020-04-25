package vm_test

import (
	"fmt"
	"log"

	"github.com/dgrr/anko/env"
	"github.com/dgrr/anko/vm"
)

func Example_vmFunctions() {
	e := env.NewEnv()

	err := e.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
fn a(b) {
	println(b)
}
a("b")

a = fn(b) {
	println(b)
}
a("b")

fn(b) {
	println(b)
}("b")

fn a() {
	return "a"
}
println(a())

println("")


fn fib(n) {
	if (n <= 1) {
		return n
	}
	return fib(n - 1) + fib(n - 2) 
}

println(fib(8))
 
 fn sum(n...) {
	 t = 0
	 for a in n {
		 t += a
	 }
	 return t
 }
println(sum(1, 2, 3, 4))

fn add(a, b) {
	return a + b
} 
println(add([1, 2]...))
`

	_, err = vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// b
	// b
	// b
	// a
	//
	// 21
	// 10
	// 3

}

func Example_vmFunctionsScope() {
	e := env.NewEnv()

	err := e.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
a = 1
fn () {
	a = 2
}()
println(a)

var a = 1
fn () {
	a = 2
}()
println(a)

a = 1
fn () {
	var a = 2
}()
println(a)

var a = 1
fn () {
	var a = 2
}()
println(a)
`

	_, err = vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// 2
	// 2
	// 1
	// 1

}

func testFunc1(a interface{}) int {
	b, ok := a.([]interface{})
	if ok {
		return len(b)
	}
	return 0
}

func Example_vmFunctionsOutside() {

	/*
	   // the following function would be uncommented
	   func testFunc1(a interface{}) int {
	   	b, ok := a.([]interface{})
	   	if ok {
	   		return len(b)
	   	}
	   	return 0
	   }
	*/

	e := env.NewEnv()

	err := e.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	err = e.Define("addString", func(a string, b string) string { return a + b })
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}
	// uses the function that would be declared above
	err = e.Define("aFunc", testFunc1)
	if err != nil {
		log.Fatalf("define error: %v\n", err)
	}

	script := `
a = addString("a", "b")
println(a)

a = aFunc([1, 2, 3])
println(a) 
`

	_, err = vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// ab
	// 3

}
