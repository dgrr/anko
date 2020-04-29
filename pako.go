// +build !appengine

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dgrr/pako/ast"
	"github.com/dgrr/pako/core"
	"github.com/dgrr/pako/env"
	_ "github.com/dgrr/pako/packages"
	"github.com/dgrr/pako/parser"
	"github.com/dgrr/pako/vm"
)

const version = "0.2.3"

var (
	flagExecute string
	file        string
	args        []string
	e           *env.Env
)

func main() {
	var exitCode int

	parseFlags()
	setupEnv()
	if flagExecute != "" || flag.NArg() > 0 {
		exitCode = runNonInteractive()
	} else {
		exitCode = runInteractive()
	}

	os.Exit(exitCode)
}

func parseFlags() {
	flagVersion := flag.Bool("v", false, "prints out the version and then exits")
	flag.StringVar(&flagExecute, "e", "", "execute the Anko code")
	flag.Parse()

	if *flagVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if flagExecute != "" || flag.NArg() < 1 {
		args = flag.Args()
		return
	}

	file = flag.Arg(0)
	args = flag.Args()[1:]
}

func setupEnv() {
	e = env.NewEnv()
	e.Define("args", args)
	e.Import = doImport
	core.Import(e)
}

func doImport(pkg string) (*env.Env, error) {
	e := env.NewEnv()
	pkg = strings.Replace(pkg, ".", string(os.PathSeparator), -1)

	d, err := ioutil.ReadFile(pkg + ".ank")
	if err != nil {
		return nil, err
	}

	_, err = vm.Execute(e, nil, string(d))
	if err != nil {
		return nil, err
	}

	return e, nil
}

func printCode(file string, err error) {
	var pos ast.Position
	switch e := err.(type) {
	case *vm.Error:
		pos = e.Pos
	case *parser.Error:
		pos = e.Pos
	default:
		fmt.Fprintf(os.Stderr, e.Error())
		return
	}
	printDefErr := func() {
		fmt.Fprintf(os.Stderr, "%d:%d %s\n", pos.Line, pos.Column, err)
	}

	d, e := ioutil.ReadFile(file)
	if e != nil {
		printDefErr()
		return
	}
	lines := strings.Split(string(d), "\n")
	for i, line := range lines {
		if i > pos.Line+5 {
			break
		}

		if i > pos.Line-15 {
			if i == pos.Line-1 {
				fmt.Print("* ")
			} else {
				fmt.Print("  ")
			}
			n := fmt.Sprintf("%d", i)
			r := 2 - len(n)
			if r > 0 {
				fmt.Print(strings.Repeat(" ", r))
			}
			fmt.Printf("%s: %s\n", n, line)
		}
	}
	fmt.Printf("%s\n", err)
}

func runNonInteractive() int {
	var source string
	if flagExecute != "" {
		source = flagExecute
	} else {
		sourceBytes, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("ReadFile error:", err)
			return 2
		}
		source = string(sourceBytes)
	}

	_, err := vm.Execute(e, nil, source)
	if err != nil {
		if flagExecute == "" {
			printCode(file, err)
		} else {
			if e, ok := err.(*vm.Error); ok {
				fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
			} else if e, ok := err.(*parser.Error); ok {
				fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		return 4
	}

	return 0
}

func runInteractive() int {
	var following bool
	var source string
	scanner := bufio.NewScanner(os.Stdin)

	parser.EnableErrorVerbose()

	for {
		if following {
			source += "\n"
			fmt.Print("  ")
		} else {
			fmt.Print("> ")
		}

		if !scanner.Scan() {
			break
		}
		source += scanner.Text()
		if source == "" {
			continue
		}
		if source == "quit()" {
			break
		}

		stmts, err := parser.ParseSrc(source)

		if e, ok := err.(*parser.Error); ok {
			es := e.Error()
			if strings.HasPrefix(es, "syntax error: unexpected") {
				if strings.HasPrefix(es, "syntax error: unexpected $end,") {
					following = true
					continue
				}
			} else {
				if e.Pos.Column == len(source) && !e.Fatal {
					fmt.Fprintln(os.Stderr, e)
					following = true
					continue
				}
				if e.Error() == "unexpected EOF" {
					following = true
					continue
				}
			}
		}

		following = false
		source = ""
		var v interface{}

		if err == nil {
			v, err = vm.Run(e, nil, stmts)
		}
		if err != nil {
			if e, ok := err.(*vm.Error); ok {
				fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
			} else if e, ok := err.(*parser.Error); ok {
				fmt.Fprintf(os.Stderr, "%d:%d %s\n", e.Pos.Line, e.Pos.Column, err)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}

		fmt.Printf("%#v\n", v)
	}

	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, "ReadString error:", err)
			return 12
		}
	}

	return 0
}
