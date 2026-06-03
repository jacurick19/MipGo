package mip

import (
	"math"
	"os"
	"testing"
)

func TestModelCreation(t *testing.T) {
	m := New("test_model")
	if m.name != "test_model" {
		t.Errorf("Expected model name to be 'test_model', got %s", m.name)
	}

	x := m.Binary("x")
	if x.vType != Binary {
		t.Errorf("Expected variable x to be Binary")
	}

	y := m.Integer("y", -10, 10)
	if y.vType != Integer || y.lb != -10 || y.ub != 10 {
		t.Errorf("Expected variable y to be Integer in [-10, 10]")
	}

	z := m.Continuous("z", 0, Inf)
	if z.vType != Continuous || z.lb != 0 || !math.IsInf(z.ub, 1) {
		t.Errorf("Expected variable z to be Continuous in [0, Inf]")
	}
}

func TestExpressions(t *testing.T) {
	m := New("expr_model")
	x := m.Binary("x")
	y := m.Continuous("y", 0, 10)

	// Single expression
	expr1 := Expr(x)
	if len(expr1.Terms()) != 1 || expr1.Terms()[0].Var != x || expr1.Terms()[0].Coeff != 1.0 {
		t.Errorf("Expr(x) incorrect structure")
	}

	// Lin
	expr2 := Lin(2.5, x, -1.0, y)
	if len(expr2.Terms()) != 2 {
		t.Errorf("Lin expression terms count mismatch")
	}
	if expr2.Terms()[0].Coeff != 2.5 || expr2.Terms()[0].Var != x {
		t.Errorf("First term mismatch in Lin")
	}
	if expr2.Terms()[1].Coeff != -1.0 || expr2.Terms()[1].Var != y {
		t.Errorf("Second term mismatch in Lin")
	}

	// Sum
	expr3 := Sum(x, y, 5.0)
	if len(expr3.Terms()) != 2 || expr3.Constant() != 5.0 {
		t.Errorf("Sum structure mismatch")
	}

	// Dot
	coeffs := []float64{1.2, 3.4}
	vars := []*Variable{x, y}
	expr4 := Dot(coeffs, vars)
	if len(expr4.Terms()) != 2 || expr4.Terms()[0].Coeff != 1.2 || expr4.Terms()[1].Coeff != 3.4 {
		t.Errorf("Dot expression mismatch")
	}

	// Immutability
	expr5 := expr2.AddConst(10.0)
	if expr2.Constant() != 0.0 {
		t.Errorf("Expression was mutated (original constant changed)")
	}
	if expr5.Constant() != 10.0 {
		t.Errorf("New expression constant not set")
	}
}

func TestConstraintNormalization(t *testing.T) {
	m := New("constr_model")
	x := m.Continuous("x", 0, 10)
	y := m.Continuous("y", 0, 10)

	// x <= y  ->  x - y <= 0
	c := x.Leq(y)
	if c.sense != 'L' {
		t.Errorf("Expected sense 'L', got %c", c.sense)
	}
	if c.rhs != 0.0 {
		t.Errorf("Expected rhs 0, got %g", c.rhs)
	}

	// Map coefficients
	coeffs := make(map[*Variable]float64)
	for _, term := range c.termsList {
		coeffs[term.Var] = term.Coeff
	}
	if coeffs[x] != 1.0 || coeffs[y] != -1.0 {
		t.Errorf("Expected coeff of x to be 1.0 and y to be -1.0, got x=%g, y=%g", coeffs[x], coeffs[y])
	}
}

func TestContainers(t *testing.T) {
	m := New("container_model")

	// Test Matrix
	matrix := m.BinaryMatrix("x", 3, 4)
	if matrix.rows != 3 || matrix.cols != 4 {
		t.Errorf("Expected matrix shape 3x4")
	}
	v := matrix.At(2, 3)
	if v == nil || v.vType != Binary {
		t.Errorf("Failed to retrieve valid variable from Matrix")
	}

	// Test Tensor
	tensor := m.BinaryTensor("y", 2, 3, 4)
	if len(tensor.dims) != 3 || tensor.dims[0] != 2 || tensor.dims[1] != 3 || tensor.dims[2] != 4 {
		t.Errorf("Expected tensor shape 2x3x4")
	}
	vt := tensor.At(1, 2, 3)
	if vt == nil || vt.vType != Binary {
		t.Errorf("Failed to retrieve valid variable from Tensor")
	}
}

func TestIndicatorsAndSOS(t *testing.T) {
	m := New("special_model")
	b := m.Binary("b")
	x := m.Continuous("x", 0, 10)
	y := m.Continuous("y", 0, 10)

	// Test Indicator
	m.Indicator(b, x.Leq(y))
	if len(m.indicators) != 1 {
		t.Errorf("Expected 1 indicator constraint")
	}
	if m.indicators[0].binaryVar != b {
		t.Errorf("Indicator binary variable mismatch")
	}

	// Test SOS
	m.AddSOS1(x, y)
	if len(m.sos1) != 1 || len(m.sos1[0].vars) != 2 {
		t.Errorf("Expected SOS1 constraint of size 2")
	}

	m.AddSOS2(x, y)
	if len(m.sos2) != 1 || len(m.sos2[0].vars) != 2 {
		t.Errorf("Expected SOS2 constraint of size 2")
	}
}

func TestLPFileIO(t *testing.T) {
	m := New("io_test")
	x := m.Binary("x")
	y := m.Continuous("y", 0.0, 10.0)

	m.Maximize(Sum(Prod(2.0, x), Prod(1.5, y)))
	m.SubjectTo(Sum(Prod(1.0, x), Prod(1.0, y)).Leq(5.0).Named("c1"))

	lpFile := "test_model.lp"
	defer os.Remove(lpFile)

	err := m.WriteLP(lpFile)
	if err != nil {
		t.Fatalf("Failed to write LP: %v", err)
	}

	imported, err := ReadLP(lpFile)
	if err != nil {
		t.Fatalf("Failed to read LP: %v", err)
	}

	if len(imported.Variables()) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(imported.Variables()))
	}
	if len(imported.Constraints()) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(imported.Constraints()))
	}

	importedVars := imported.Variables()
	var importedX, importedY *Variable
	for _, v := range importedVars {
		if v.name == "x" {
			importedX = v
		} else if v.name == "y" {
			importedY = v
		}
	}

	if importedX == nil || importedY == nil {
		t.Fatalf("Failed to locate variables by name in imported model")
	}

	if importedX.vType != Binary {
		t.Errorf("Expected imported x to be Binary, got %v", importedX.vType)
	}
	if importedY.vType != Continuous || importedY.ub != 10.0 {
		t.Errorf("Expected imported y to be Continuous in [0, 10], got ub=%g", importedY.ub)
	}
}

func BenchmarkVariableTerms(b *testing.B) {
	m := New("bench")
	x := m.Continuous("x", 0, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = x.Terms()
	}
}


