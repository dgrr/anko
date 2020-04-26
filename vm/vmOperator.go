package vm

import (
	"reflect"
	"strings"

	"github.com/dgrr/pako/ast"
	"github.com/dgrr/pako/v0.2.1/over"
)

// invokeOperator evaluates one Operator.
func (runInfo *runInfoStruct) invokeOperator() {
	switch operator := runInfo.operator.(type) {

	// BinaryOperator
	case *ast.BinaryOperator:
		runInfo.expr = operator.LHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}

		switch operator.Operator {
		case "||":
			if toBool(runInfo.rv) {
				runInfo.rv = trueValue
				return
			}
		case "&&":
			if !toBool(runInfo.rv) {
				runInfo.rv = falseValue
				return
			}
		default:
			runInfo.err = newStringError(operator, "unknown operator")
			runInfo.rv = nilValue
			return
		}

		runInfo.expr = operator.RHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}

		if toBool(runInfo.rv) {
			runInfo.rv = trueValue
		} else {
			runInfo.rv = falseValue
		}

	// ComparisonOperator
	case *ast.ComparisonOperator:
		runInfo.expr = operator.LHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}
		lhsV := runInfo.rv

		runInfo.expr = operator.RHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}

		switch operator.Operator {
		case "==":
			runInfo.rv = reflect.ValueOf(equal(lhsV, runInfo.rv))
		case "!=":
			runInfo.rv = reflect.ValueOf(!equal(lhsV, runInfo.rv))
		case "<":
			if lhsV.Type().Implements(over.ComparisonReflectType) {
				lhv := lhsV.Interface().(over.Comparison)
				v := getUnderlayedType(runInfo.rv)
				runInfo.err = lhv.Less(v)
				if runInfo.err == nil {
					runInfo.rv = trueValue
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toFloat64(lhsV) < toFloat64(runInfo.rv))
		case "<=":
			if lhsV.Type().Implements(over.ComparisonReflectType) {
				lhv := lhsV.Interface().(over.Comparison)
				v := getUnderlayedType(runInfo.rv)
				runInfo.err = lhv.LessEquals(v)
				if runInfo.err == nil {
					runInfo.rv = trueValue
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toFloat64(lhsV) <= toFloat64(runInfo.rv))
		case ">":
			if lhsV.Type().Implements(over.ComparisonReflectType) {
				lhv := lhsV.Interface().(over.Comparison)
				v := getUnderlayedType(runInfo.rv)
				runInfo.err = lhv.Greater(v)
				if runInfo.err == nil {
					runInfo.rv = trueValue
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toFloat64(lhsV) > toFloat64(runInfo.rv))
		case ">=":
			if lhsV.Type().Implements(over.ComparisonReflectType) {
				lhv := lhsV.Interface().(over.Comparison)
				v := getUnderlayedType(runInfo.rv)
				runInfo.err = lhv.GreaterEquals(v)
				if runInfo.err == nil {
					runInfo.rv = trueValue
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toFloat64(lhsV) >= toFloat64(runInfo.rv))
		default:
			runInfo.err = newStringError(operator, "unknown operator")
			runInfo.rv = nilValue
		}

	// AddOperator
	case *ast.AddOperator:
		runInfo.expr = operator.LHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}
		lhsV := runInfo.rv

		runInfo.expr = operator.RHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}

		switch operator.Operator {
		case "+":
			if lhsV.Type().Implements(over.AddReflectType) {
				adder := lhsV.Interface().(over.Add)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = adder.Add(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			lhsKind := lhsV.Kind()
			rhsKind := runInfo.rv.Kind()

			if lhsKind == reflect.Slice || lhsKind == reflect.Array {
				if rhsKind == reflect.Slice || rhsKind == reflect.Array {
					// append slice to slice
					runInfo.rv, runInfo.err = appendSlice(operator, lhsV, runInfo.rv)
					return
				}
				// try to append rhs non-slice to lhs slice
				runInfo.rv, runInfo.err = convertReflectValueToType(runInfo.rv, lhsV.Type().Elem())
				if runInfo.err != nil {
					runInfo.err = newStringError(operator, "invalid type conversion")
					runInfo.rv = nilValue
					return
				}
				runInfo.rv = reflect.Append(lhsV, runInfo.rv)
				return
			}
			if rhsKind == reflect.Slice || rhsKind == reflect.Array {
				// can not append rhs slice to lhs non-slice
				runInfo.err = newStringError(operator, "invalid type conversion")
				runInfo.rv = nilValue
				return
			}

			kind := precedenceOfKinds(lhsKind, rhsKind)
			switch kind {
			case reflect.String:
				runInfo.rv = reflect.ValueOf(toString(lhsV) + toString(runInfo.rv))
			case reflect.Float64, reflect.Float32:
				runInfo.rv = reflect.ValueOf(toFloat64(lhsV) + toFloat64(runInfo.rv))
			default:
				runInfo.rv = reflect.ValueOf(toInt64(lhsV) + toInt64(runInfo.rv))
			}

		case "-":
			if lhsV.Type().Implements(over.AddReflectType) {
				adder := lhsV.Interface().(over.Add)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = adder.Sub(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			switch lhsV.Kind() {
			case reflect.Float64, reflect.Float32:
				runInfo.rv = reflect.ValueOf(toFloat64(lhsV) - toFloat64(runInfo.rv))
				return
			}
			switch runInfo.rv.Kind() {
			case reflect.Float64, reflect.Float32:
				runInfo.rv = reflect.ValueOf(toFloat64(lhsV) - toFloat64(runInfo.rv))
			default:
				runInfo.rv = reflect.ValueOf(toInt64(lhsV) - toInt64(runInfo.rv))
			}

		case "|":
			if lhsV.Type().Implements(over.AddReflectType) {
				adder := lhsV.Interface().(over.Add)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = adder.Or(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}
			runInfo.rv = reflect.ValueOf(toInt64(lhsV) | toInt64(runInfo.rv))
		default:
			runInfo.err = newStringError(operator, "unknown operator")
			runInfo.rv = nilValue
		}

	// MultiplyOperator
	case *ast.MultiplyOperator:
		runInfo.expr = operator.LHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}
		lhsV := runInfo.rv

		runInfo.expr = operator.RHS
		runInfo.invokeExpr()
		if runInfo.err != nil {
			return
		}
		if runInfo.rv.Kind() == reflect.Interface && !runInfo.rv.IsNil() {
			runInfo.rv = runInfo.rv.Elem()
		}

		switch operator.Operator {
		case "*":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.Mul(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			if lhsV.Kind() == reflect.String && (runInfo.rv.Kind() == reflect.Int || runInfo.rv.Kind() == reflect.Int32 || runInfo.rv.Kind() == reflect.Int64) {
				runInfo.rv = reflect.ValueOf(strings.Repeat(toString(lhsV), int(toInt64(runInfo.rv))))
				return
			}
			if lhsV.Kind() == reflect.Float64 || runInfo.rv.Kind() == reflect.Float64 {
				runInfo.rv = reflect.ValueOf(toFloat64(lhsV) * toFloat64(runInfo.rv))
				return
			}
			runInfo.rv = reflect.ValueOf(toInt64(lhsV) * toInt64(runInfo.rv))
		case "/":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.Div(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toFloat64(lhsV) / toFloat64(runInfo.rv))
		case "%":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.Mod(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toInt64(lhsV) % toInt64(runInfo.rv))
		case ">>":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.Right(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toInt64(lhsV) >> uint64(toInt64(runInfo.rv)))
		case "<<":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.Left(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toInt64(lhsV) << uint64(toInt64(runInfo.rv)))
		case "&":
			if lhsV.Type().Implements(over.MultiplyReflectType) {
				mul := lhsV.Interface().(over.Multiply)
				v := getUnderlayedType(runInfo.rv)
				v, runInfo.err = mul.And(v)
				if runInfo.err == nil {
					runInfo.rv = reflect.ValueOf(v)
				}
				return
			}

			runInfo.rv = reflect.ValueOf(toInt64(lhsV) & toInt64(runInfo.rv))

		default:
			runInfo.err = newStringError(operator, "unknown operator")
			runInfo.rv = nilValue
		}

	default:
		runInfo.err = newStringError(operator, "unknown operator")
		runInfo.rv = nilValue

	}
}
