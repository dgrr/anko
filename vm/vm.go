package vm

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/dgrr/pako/ast"
	"github.com/dgrr/pako/env"
	"github.com/dgrr/pako/over"
)

// Options provides options to run VM with
type Options struct {
	// DisableGo disables the usage of `go` statements.
	DisableGo      bool
	MaxMemoryUsage int64
	Debug          bool // run in Debug mode
}

type (
	// Error is a VM run error.
	Error struct {
		Message string
		Pos     ast.Position
	}

	// runInfo provides run incoming and outgoing information
	runInfoStruct struct {
		recv     reflect.Value
		memUsage int64
		vmTypes  []string

		// incoming
		ctx      context.Context
		env      *env.Env
		options  *Options
		stmt     ast.Stmt
		expr     ast.Expr
		operator ast.Operator

		// outgoing
		rv      reflect.Value
		callErr int
		err     error
	}

	vmStruct struct {
		v reflect.Value
	}
)

func (runInfo *runInfoStruct) runDecls(ctx context.Context) {
	mStmt := runInfo.stmt
	defer func() {
		runInfo.stmt = mStmt
	}()

	if stmts, ok := runInfo.stmt.(*ast.StmtsStmt); ok {
		for i := 0; runInfo.err == nil && i < len(stmts.Stmts); i++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			stmt := stmts.Stmts[i]
			remove := false
			switch st := stmt.(type) {
			case *ast.StructStmt: // struct declaration
				runInfo.stmt = st
				runInfo.runSingleStmt() // struct declared if no nil is returned
				remove = true
			case *ast.ExprStmt: // can be a function
				if fn, ok := st.Expr.(*ast.FuncExpr); ok && fn.Recv != "" {
					styp, err := runInfo.env.Type(fn.Recv)
					if err != nil {
						runInfo.err = newStringError(fn, fn.Recv+" not declared at this point")
					} else {
						runInfo.expr = fn
						ftyp := runInfo.funcExpr()
						runInfo.env.DefineMethod(fn.Recv+"."+fn.Name, runInfo.rv)

						ns := make([]reflect.StructField, styp.NumField()+1)
						for i := 0; i < styp.NumField(); i++ {
							ns[i] = styp.Field(i)
						}
						ns[len(ns)-1] = reflect.StructField{
							Name: fn.Name,
							Type: ftyp,
						}

						runInfo.rv = nilValue
						runInfo.env.DefineReflectType(fn.Recv, reflect.StructOf(ns))
						remove = true
					}
				}
			}
			if remove {
				stmts.Stmts = append(stmts.Stmts[:i], stmts.Stmts[i+1:]...)
				i--
			}
		}
	}
}

func getUnderlyingType(v reflect.Value) interface{} {
	for v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}
	return v.Interface()
}

func getSizeOf(v reflect.Value) int64 {
	for v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}
	x := int(v.Type().Size())
	switch v.Kind() {
	case reflect.Slice, reflect.Array,
		reflect.Chan, reflect.String, reflect.Map:
		x = 0
		if v.Kind() == reflect.Map {
			it := v.MapRange()
			for it.Next() {
				vl := it.Key()
				for vl.Type() == reflectValueType {
					vl = v.Interface().(reflect.Value)
				}
				x += int(vl.Type().Size())
				vl = it.Value()
				for vl.Type() == reflectValueType {
					vl = v.Interface().(reflect.Value)
				}
				x += int(vl.Type().Size())
			}
		} else {
			for i := 0; i < v.Len(); i++ {
				vl := v.Index(i)
				for vl.Type() == reflectValueType {
					vl = v.Interface().(reflect.Value)
				}
				x += int(vl.Type().Size())
			}
		}
	}

	return int64(x)
}

func (runInfo *runInfoStruct) allocMemUsage(v reflect.Value) {
	//runInfo.memUsage += getSizeOf(v)
}

func (runInfo *runInfoStruct) deallocMemUsage(v reflect.Value) {
	//runInfo.memUsage -= getSizeOf(v)
}

var (
	nilType            = reflect.TypeOf(nil)
	stringType         = reflect.TypeOf("a")
	byteType           = reflect.TypeOf(byte('a'))
	runeType           = reflect.TypeOf('a')
	interfaceType      = reflect.ValueOf([]interface{}{int64(1)}).Index(0).Type()
	interfaceSliceType = reflect.TypeOf([]interface{}{})
	reflectValueType   = reflect.TypeOf(reflect.Value{})
	errorType          = reflect.ValueOf([]error{nil}).Index(0).Type()
	vmErrorType        = reflect.TypeOf(&Error{})
	contextType        = reflect.TypeOf((*context.Context)(nil)).Elem()

	nilValue                  = reflect.New(reflect.TypeOf((*interface{})(nil)).Elem()).Elem()
	trueValue                 = reflect.ValueOf(true)
	falseValue                = reflect.ValueOf(false)
	zeroValue                 = reflect.Value{}
	reflectValueNilValue      = reflect.ValueOf(nilValue)
	reflectValueErrorNilValue = reflect.ValueOf(reflect.New(errorType).Elem())

	errInvalidTypeConversion = fmt.Errorf("invalid type conversion")

	// ErrBreak when there is an unexpected break statement
	ErrBreak = errors.New("unexpected break statement")
	// ErrContinue when there is an unexpected continue statement
	ErrContinue = errors.New("unexpected continue statement")
	// ErrReturn when there is an unexpected return statement
	ErrReturn = errors.New("unexpected return statement")
	// ErrInterrupt when execution has been interrupted
	ErrInterrupt = errors.New("execution interrupted")
)

// Error returns the VM error message.
func (e *Error) Error() string {
	return e.Message
}

// newError makes VM error from error
func newError(pos ast.Pos, err error) error {
	if err == nil {
		return nil
	}
	if pos == nil {
		return &Error{Message: err.Error(), Pos: ast.Position{Line: 1, Column: 1}}
	}
	return &Error{Message: err.Error(), Pos: pos.Position()}
}

// newStringError makes VM error from string
func newStringError(pos ast.Pos, err string) error {
	if err == "" {
		return nil
	}
	if pos == nil {
		return &Error{Message: err, Pos: ast.Position{Line: 1, Column: 1}}
	}
	return &Error{Message: err, Pos: pos.Position()}
}

// recoverFunc generic recover function
func recoverFunc(runInfo *runInfoStruct) {
	recoverInterface := recover()
	if recoverInterface == nil {
		return
	}
	switch value := recoverInterface.(type) {
	case *Error:
		runInfo.err = value
	case error:
		runInfo.err = value
	default:
		runInfo.err = fmt.Errorf("%v", recoverInterface)
	}
}

func isNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		// from reflect IsNil:
		// Note that IsNil is not always equivalent to a regular comparison with nil in Go.
		// For example, if v was created by calling ValueOf with an uninitialized interface variable i,
		// i==nil will be true but v.IsNil will panic as v will be the zero Value.
		return v.IsNil()
	default:
		return false
	}
}

func isNum(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// equal returns true when lhsV and rhsV is same value.
func equal(lhsV, rhsV reflect.Value) bool {
	lhsIsNil, rhsIsNil := isNil(lhsV), isNil(rhsV)
	if lhsIsNil && rhsIsNil {
		return true
	}
	if (!lhsIsNil && rhsIsNil) || (lhsIsNil && !rhsIsNil) {
		return false
	}
	if lhsV.Kind() == reflect.Interface || lhsV.Kind() == reflect.Ptr {
		lhsV = lhsV.Elem()
	}
	if rhsV.Kind() == reflect.Interface || rhsV.Kind() == reflect.Ptr {
		rhsV = rhsV.Elem()
	}
	if lhsV.Type().Implements(over.ComparisonReflectType) {
		lhv := lhsV.Interface().(over.Comparison)
		v := getUnderlyingType(rhsV)
		return lhv.Equals(v) == nil
	}

	// Compare a string and a number.
	// This will attempt to convert the string to a number,
	// while leaving the other side alone. Code further
	// down takes care of converting ints and floats as needed.
	if isNum(lhsV) && rhsV.Kind() == reflect.String {
		rhsF, err := tryToFloat64(rhsV)
		if err != nil {
			// Couldn't convert RHS to a float, they can't be compared.
			return false
		}
		rhsV = reflect.ValueOf(rhsF)
	} else if lhsV.Kind() == reflect.String && isNum(rhsV) {
		// If the LHS is a string formatted as an int, try that before trying float
		lhsI, err := tryToInt64(lhsV)
		if err != nil {
			// if LHS is a float, e.g. "1.2", we need to set lhsV to a float64
			lhsF, err := tryToFloat64(lhsV)
			if err != nil {
				return false
			}
			lhsV = reflect.ValueOf(lhsF)
		} else {
			lhsV = reflect.ValueOf(lhsI)
		}
	}

	if isNum(lhsV) && isNum(rhsV) {
		return fmt.Sprintf("%v", lhsV) == fmt.Sprintf("%v", rhsV)
	}

	// Try to compare bools to strings and numbers
	if lhsV.Kind() == reflect.Bool || rhsV.Kind() == reflect.Bool {
		lhsB, err := tryToBool(lhsV)
		if err != nil {
			return false
		}
		rhsB, err := tryToBool(rhsV)
		if err != nil {
			return false
		}
		return lhsB == rhsB
	}

	return reflect.DeepEqual(lhsV.Interface(), rhsV.Interface())
}

func getMapIndex(key reflect.Value, aMap reflect.Value) reflect.Value {
	if aMap.IsNil() {
		return nilValue
	}

	var err error
	key, err = convertReflectValueToType(key, aMap.Type().Key())
	if err != nil {
		return nilValue
	}

	// From reflect MapIndex:
	// It returns the zero Value if key is not found in the map or if v represents a nil map.
	value := aMap.MapIndex(key)
	if !value.IsValid() {
		return nilValue
	}

	if aMap.Type().Elem() == interfaceType && !value.IsNil() {
		value = reflect.ValueOf(value.Interface())
	}

	return value
}

// appendSlice appends rhs to lhs
// function assumes lhsV and rhsV are slice or array
func appendSlice(expr ast.Expr, lhsV reflect.Value, rhsV reflect.Value) (reflect.Value, error) {
	lhsT := lhsV.Type().Elem()
	rhsT := rhsV.Type().Elem()

	if lhsT == rhsT {
		return reflect.AppendSlice(lhsV, rhsV), nil
	}

	if rhsT.ConvertibleTo(lhsT) {
		for i := 0; i < rhsV.Len(); i++ {
			lhsV = reflect.Append(lhsV, rhsV.Index(i).Convert(lhsT))
		}
		return lhsV, nil
	}

	leftHasSubArray := lhsT.Kind() == reflect.Slice || lhsT.Kind() == reflect.Array
	rightHasSubArray := rhsT.Kind() == reflect.Slice || rhsT.Kind() == reflect.Array

	if leftHasSubArray != rightHasSubArray && lhsT != interfaceType && rhsT != interfaceType {
		return nilValue, newStringError(expr, "invalid type conversion")
	}

	if !leftHasSubArray && !rightHasSubArray {
		for i := 0; i < rhsV.Len(); i++ {
			value := rhsV.Index(i)
			if rhsT == interfaceType {
				value = value.Elem()
			}
			if lhsT == value.Type() {
				lhsV = reflect.Append(lhsV, value)
			} else if value.Type().ConvertibleTo(lhsT) {
				lhsV = reflect.Append(lhsV, value.Convert(lhsT))
			} else {
				return nilValue, newStringError(expr, "invalid type conversion")
			}
		}
		return lhsV, nil
	}

	if (leftHasSubArray || lhsT == interfaceType) && (rightHasSubArray || rhsT == interfaceType) {
		for i := 0; i < rhsV.Len(); i++ {
			value := rhsV.Index(i)
			if rhsT == interfaceType {
				value = value.Elem()
				if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
					return nilValue, newStringError(expr, "invalid type conversion")
				}
			}
			newSlice, err := appendSlice(expr, reflect.MakeSlice(lhsT, 0, value.Len()), value)
			if err != nil {
				return nilValue, err
			}
			lhsV = reflect.Append(lhsV, newSlice)
		}
		return lhsV, nil
	}

	return nilValue, newStringError(expr, "invalid type conversion")
}

func makeType(runInfo *runInfoStruct, typeStruct *ast.TypeStruct) reflect.Type {
	switch typeStruct.Kind {
	case ast.TypeDefault:
		return getTypeFromEnv(runInfo, typeStruct)
	case ast.TypePtr:
		var t reflect.Type
		if typeStruct.SubType != nil {
			t = makeType(runInfo, typeStruct.SubType)
		} else {
			t = getTypeFromEnv(runInfo, typeStruct)
		}
		if runInfo.err != nil {
			return nil
		}
		if t == nil {
			return nil
		}
		return reflect.PtrTo(t)
	case ast.TypeSlice:
		var t reflect.Type
		if typeStruct.SubType != nil {
			t = makeType(runInfo, typeStruct.SubType)
		} else {
			t = getTypeFromEnv(runInfo, typeStruct)
		}
		if runInfo.err != nil {
			return nil
		}
		if t == nil {
			return nil
		}
		for i := 1; i < typeStruct.Dimensions; i++ {
			t = reflect.SliceOf(t)
		}
		return reflect.SliceOf(t)
	case ast.TypeMap:
		key := makeType(runInfo, typeStruct.Key)
		if runInfo.err != nil {
			return nil
		}
		if key == nil {
			return nil
		}
		t := makeType(runInfo, typeStruct.SubType)
		if runInfo.err != nil {
			return nil
		}
		if t == nil {
			return nil
		}
		if !runInfo.options.Debug {
			// captures panic
			defer recoverFunc(runInfo)
		}
		t = reflect.MapOf(key, t)
		return t
	case ast.TypeChan:
		var t reflect.Type
		if typeStruct.SubType != nil {
			t = makeType(runInfo, typeStruct.SubType)
		} else {
			t = getTypeFromEnv(runInfo, typeStruct)
		}
		if runInfo.err != nil {
			return nil
		}
		if t == nil {
			return nil
		}
		return reflect.ChanOf(reflect.BothDir, t)
	case ast.TypeStructType:
		var t reflect.Type
		fields := make([]reflect.StructField, 0, len(typeStruct.StructNames))
		for i := 0; i < len(typeStruct.StructNames); i++ {
			t = makeType(runInfo, typeStruct.StructTypes[i])
			if runInfo.err != nil {
				return nil
			}
			if t == nil {
				return nil
			}

			fields = append(fields, reflect.StructField{Name: typeStruct.StructNames[i], Type: t})
		}
		if !runInfo.options.Debug {
			// captures panic
			defer recoverFunc(runInfo)
		}
		t = reflect.StructOf(fields)
		return t
	default:
		runInfo.err = fmt.Errorf("unknown kind")
		return nil
	}
}

func getTypeFromEnv(runInfo *runInfoStruct, typeStruct *ast.TypeStruct) reflect.Type {
	var e *env.Env
	e, runInfo.err = runInfo.env.GetEnvFromPath(typeStruct.Env)
	if runInfo.err != nil {
		return nil
	}

	var t reflect.Type
	t, runInfo.err = e.Type(typeStruct.Name)
	return t
}

func (runInfo *runInfoStruct) isVMType(name string) (b bool) {
	if runInfo == nil {
		return
	}

	for _, T := range runInfo.vmTypes {
		if b = T == name; b {
			break
		}
	}
	return
}

func makeValue(runInfo *runInfoStruct, name string, t reflect.Type) (reflect.Value, error) {
	if t.Implements(over.MemReflectType) {
		tl := reflect.New(t).Elem().Interface().(over.Mem)
		v, err := tl.New()

		return reflect.ValueOf(v), err
	}

	isVM := runInfo.isVMType(name)
	switch t.Kind() {
	case reflect.Chan:
		return reflect.MakeChan(t, 0), nil
	case reflect.Func:
		if t.NumIn() < 2 || t.In(0) != contextType {
			return reflect.MakeFunc(t, nil), nil
		}
		v := t.In(1)
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		V, err := runInfo.env.Method(name)
		if err != nil {
			return reflect.MakeFunc(t, nil), nil
		}

		return V, nil
	case reflect.Map:
		// note creating slice as work around to create map
		// just doing MakeMap can give incorrect type for defined types
		value := reflect.MakeSlice(reflect.SliceOf(t), 0, 1)
		value = reflect.Append(value, reflect.MakeMap(reflect.MapOf(t.Key(), t.Elem())))
		return value.Index(0), nil
	case reflect.Ptr:
		ptrV := reflect.New(t.Elem())
		v, err := makeValue(runInfo, "", t.Elem())
		if err != nil {
			return nilValue, err
		}

		ptrV.Elem().Set(v)
		return ptrV, nil
	case reflect.Slice:
		return reflect.MakeSlice(t, 0, 0), nil
	case reflect.Struct:
		structV := reflect.New(t).Elem()
		for i := 0; i < structV.NumField(); i++ {
			if structV.Field(i).Kind() == reflect.Ptr {
				continue
			}
			v, err := makeValue(runInfo, name+"."+t.Field(i).Name, structV.Field(i).Type())
			if err != nil {
				return nilValue, err
			}
			if structV.Field(i).CanSet() {
				structV.Field(i).Set(v)
			}
		}
		if isVM {
			structV = reflect.ValueOf(&vmStruct{structV})
		}
		return structV, nil
	}
	v := reflect.New(t).Elem()
	if isVM {
		v = reflect.ValueOf(&vmStruct{v})
	}

	return v, nil
}

// precedenceOfKinds returns the greater of two kinds
// string > float > int
func precedenceOfKinds(kind1 reflect.Kind, kind2 reflect.Kind) reflect.Kind {
	if kind1 == kind2 {
		return kind1
	}
	switch kind1 {
	case reflect.String:
		return kind1
	case reflect.Float64, reflect.Float32:
		switch kind2 {
		case reflect.String:
			return kind2
		}
		return kind1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch kind2 {
		case reflect.String, reflect.Float64, reflect.Float32:
			return kind2
		}
	}
	return kind1
}
