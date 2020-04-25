package vm_test

import (
	"log"

	"github.com/dgrr/anko/env"
	_ "github.com/dgrr/anko/packages"
	"github.com/dgrr/anko/vm"
)

func Example_vmSort() {
	// _ "github.com/dgrr/anko/packages"

	e := env.NewEnv()

	script := `
import fmt
import sort
a = [5, 1.1, 3, "f", "2", "4.4"]
sortFuncs = make(sort.SortFuncsStruct)
sortFuncs.LenFunc = func() { return len(a) }
sortFuncs.LessFunc = fn(i, j) { return a[i] < a[j] }
sortFuncs.SwapFunc = fn(i, j) { temp = a[i]; a[i] = a[j]; a[j] = temp }
sort.Sort(sortFuncs)
fmt.Println(a)
`

	_, err := vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// [f 1.1 2 3 4.4 5]
}

func Example_vmRegexp() {
	// _ "github.com/dgrr/anko/packages"

	e := env.NewEnv()

	script := `
import fmt
import regexp

re = regexp.MustCompile("^simple$")
result = re.MatchString("simple")
fmt.Println(result)
fmt.Println("")

re = regexp.MustCompile("simple")
result = re.FindString("This is a simple sentence")
fmt.Println(result)
fmt.Println("")

re = regexp.MustCompile(",")
result = re.Split("a,b,c", -1)
fmt.Println(result)
fmt.Println("")

re = regexp.MustCompile("foo")
result = re.ReplaceAllString("foo", "bar")
fmt.Println(result)
`

	_, err := vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// true
	//
	// simple
	//
	// [a b c]
	//
	// bar
}

func Example_vmHttp() {
	// _ "github.com/dgrr/anko/packages"

	e := env.NewEnv()

	script := `
import fmt
import io
import io/ioutil
import net
import net/http
import time

fn handlerRoot(responseWriter, request) {
	io.WriteString(responseWriter, "Hello World :)")
}

serveMux = http.NewServeMux()
serveMux.HandleFunc("/", handlerRoot)
listener, err = net.Listen("tcp", ":8084")
if err != nil {
	fmt.Println(err)
	return
}
go http.Serve(listener, serveMux)

client = http.DefaultClient

response, err = client.Get("http://localhost:8084/")
if err != nil {
	fmt.Println(err)
	return
}

body, err = ioutil.ReadAll(response.Body)
if err != nil {
	fmt.Println(err)
}
response.Body.Close()

fmt.Printf("%s\n", body)
`

	_, err := vm.Execute(e, nil, script)
	if err != nil {
		log.Fatalf("execute error: %v\n", err)
	}

	// output:
	// Hello World :)
}
