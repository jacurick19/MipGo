package mip

import (
	"reflect"
)

// --- Expr Interface Implementation ---

// Implement Expr for Variable
func (v *Variable) Terms() []Term {
	return v.termArray[:]
}

func (v *Variable) Constant() float64 {
	return 0.0
}

func (v *Variable) Eq(rhs any) *Constraint {
	return newConstraint(v, rhs, 'E')
}

func (v *Variable) Leq(rhs any) *Constraint {
	return newConstraint(v, rhs, 'L')
}

func (v *Variable) Geq(rhs any) *Constraint {
	return newConstraint(v, rhs, 'G')
}

// Implement Expr for Expression
var _ Expr = Expression{}

func (e Expression) Terms() []Term {
	return e.termsList
}

func (e Expression) Constant() float64 {
	return e.constVal
}

func (e Expression) Eq(rhs any) *Constraint {
	return newConstraint(e, rhs, 'E')
}

func (e Expression) Leq(rhs any) *Constraint {
	return newConstraint(e, rhs, 'L')
}

func (e Expression) Geq(rhs any) *Constraint {
	return newConstraint(e, rhs, 'G')
}

func (e Expression) AddConst(val float64) Expression {
	newTerms := make([]Term, len(e.termsList))
	copy(newTerms, e.termsList)
	return Expression{
		termsList: newTerms,
		constVal:  e.constVal + val,
	}
}

// --- Expression Helpers ---

// Const creates a constant expression.
func Const(val float64) Expression {
	return Expression{constVal: val}
}



// Lin constructs a linear expression from coefficient-variable pairs.
// Example: Lin(2, x, -3, y) represents 2x - 3y
func Lin(args ...any) Expression {
	if len(args)%2 != 0 {
		panic("Lin requires an even number of arguments (coefficient/variable pairs)")
	}
	var terms []Term
	var constant float64
	for i := 0; i < len(args); i += 2 {
		coeffVal := toFloatOrPanic(args[i])
		varParam := args[i+1]
		switch v := varParam.(type) {
		case *Variable:
			terms = append(terms, Term{Coeff: coeffVal, Var: v})
		case Expr:
			for _, t := range v.Terms() {
				terms = append(terms, Term{Coeff: t.Coeff * coeffVal, Var: t.Var})
			}
			constant += v.Constant() * coeffVal
		default:
			panic("Lin second argument of a pair must be a Variable or an Expression")
		}
	}
	return Expression{termsList: terms, constVal: constant}
}

// Sum aggregates multiple variables, expressions, or constants.
func Sum(args ...any) Expression {
	var terms []Term
	var constant float64
	for _, arg := range args {
		switch v := arg.(type) {
		case *Variable:
			terms = append(terms, Term{Coeff: 1.0, Var: v})
		case Expr:
			terms = append(terms, v.Terms()...)
			constant += v.Constant()
		default:
			if f, ok := toFloat(arg); ok {
				constant += f
			} else {
				panic("Sum argument must be numeric, a Variable, or an Expression")
			}
		}
	}
	return Expression{termsList: terms, constVal: constant}
}

type Numeric interface {
	~float64 | ~float32 | ~int | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

// Dot computes the dot product of a slice of coefficients and a slice of variables/expressions.
func Dot[T Numeric, V Expr](coeffs []T, variables []V) Expression {
	n := len(coeffs)
	if len(variables) < n {
		n = len(variables)
	}
	var terms []Term
	var constant float64
	for i := 0; i < n; i++ {
		c := float64(coeffs[i])
		vItem := variables[i]
		if vVar, ok := any(vItem).(*Variable); ok {
			terms = append(terms, Term{Coeff: c, Var: vVar})
		} else {
			for _, t := range vItem.Terms() {
				terms = append(terms, Term{Coeff: t.Coeff * c, Var: t.Var})
			}
			constant += vItem.Constant() * c
		}
	}
	return Expression{termsList: terms, constVal: constant}
}

// Prod multiplies a variable or expression by a coefficient.
func Prod(coeff float64, expr any) Expression {
	return Lin(coeff, expr)
}

// Neg negates a variable or expression.
func Neg(expr any) Expression {
	return Prod(-1.0, expr)
}

// --- Summation Utilities ---

// SumVars sums a set of variables.
func SumVars(vars ...*Variable) Expression {
	terms := make([]Term, len(vars))
	for i, v := range vars {
		terms[i] = Term{Coeff: 1.0, Var: v}
	}
	return Expression{termsList: terms}
}

// SumOver sums an expression over indices.
// nOrSlice can be an int (loop from 0 to n-1) or a slice (loop from 0 to len-1).
func SumOver(nOrSlice any, f func(i int) Expr) Expression {
	var limit int
	switch v := nOrSlice.(type) {
	case int:
		limit = v
	default:
		val := reflect.ValueOf(nOrSlice)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			limit = val.Len()
		} else {
			panic("SumOver first argument must be an int or a slice")
		}
	}
	terms := make([]Term, 0, limit)
	var constant float64
	for i := 0; i < limit; i++ {
		expr := f(i)
		if expr != nil {
			terms = append(terms, expr.Terms()...)
			constant += expr.Constant()
		}
	}
	return Expression{termsList: terms, constVal: constant}
}

// Sum2D sums an expression over two sets of indices.
func Sum2D(rows any, cols any, f func(i, j int) Expr) Expression {
	var rLimit, cLimit int
	// rows
	switch v := rows.(type) {
	case int:
		rLimit = v
	default:
		val := reflect.ValueOf(rows)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			rLimit = val.Len()
		} else {
			panic("Sum2D rows argument must be an int or a slice")
		}
	}
	// cols
	switch v := cols.(type) {
	case int:
		cLimit = v
	default:
		val := reflect.ValueOf(cols)
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			cLimit = val.Len()
		} else {
			panic("Sum2D cols argument must be an int or a slice")
		}
	}
	terms := make([]Term, 0, rLimit*cLimit)
	var constant float64
	for i := 0; i < rLimit; i++ {
		for j := 0; j < cLimit; j++ {
			expr := f(i, j)
			if expr != nil {
				terms = append(terms, expr.Terms()...)
				constant += expr.Constant()
			}
		}
	}
	return Expression{termsList: terms, constVal: constant}
}

// --- Constraint Chaining ---

func (c *Constraint) Named(name string) *Constraint {
	c.name = name
	return c
}

// --- Internal Helpers ---

func toFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case int16:
		return float64(v), true
	case int8:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint8:
		return float64(v), true
	default:
		return 0.0, false
	}
}

func toFloatOrPanic(val any) float64 {
	f, ok := toFloat(val)
	if !ok {
		panic("expected numeric value")
	}
	return f
}

func newConstraint(lhs Expr, rhs any, sense byte) *Constraint {
	var rExpr Expr
	switch v := rhs.(type) {
	case *Variable:
		rExpr = v
	case Expression:
		rExpr = v
	default:
		if e, ok := rhs.(Expr); ok {
			rExpr = e
		} else if f, ok := toFloat(rhs); ok {
			rExpr = Const(f)
		} else {
			panic("rhs must be numeric, a Variable, or an Expression")
		}
	}

	diffConstant := lhs.Constant() - rExpr.Constant()
	rhsVal := -diffConstant

	if len(rExpr.Terms()) == 0 {
		return &Constraint{
			termsList: lhs.Terms(),
			rhs:       rhsVal,
			sense:     sense,
		}
	}

	coeffMap := make(map[*Variable]float64)
	for _, t := range lhs.Terms() {
		coeffMap[t.Var] += t.Coeff
	}
	for _, t := range rExpr.Terms() {
		coeffMap[t.Var] -= t.Coeff
	}

	var terms []Term
	for v, coeff := range coeffMap {
		if coeff != 0.0 {
			terms = append(terms, Term{Coeff: coeff, Var: v})
		}
	}

	return &Constraint{
		termsList: terms,
		rhs:       rhsVal,
		sense:     sense,
	}
}
