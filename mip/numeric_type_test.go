package mip_test

import (
	"testing"

	. "mipgo/v2/mip"
)

func TestNumericTypesInLin(t *testing.T) {
	m := New("model")
	x := m.Binary("x")

	// Helper to check different numeric types as coefficients in Lin
	testCases := []struct {
		val      any
		expected float64
	}{
		{int(5), 5.0},
		{int64(5), 5.0},
		{int32(5), 5.0},
		{int16(5), 5.0},
		{int8(5), 5.0},
		{uint(5), 5.0},
		{uint64(5), 5.0},
		{uint32(5), 5.0},
		{uint16(5), 5.0},
		{uint8(5), 5.0},
		{float64(5.5), 5.5},
		{float32(5.5), 5.5},
	}

	for _, tc := range testCases {
		expr := Lin(tc.val, x)
		if len(expr.Terms()) != 1 {
			t.Errorf("Expected 1 term for type %T, got %d", tc.val, len(expr.Terms()))
			continue
		}
		coeff := expr.Terms()[0].Coeff
		if coeff != tc.expected {
			t.Errorf("For type %T: expected coefficient %g, got %g (silent conversion bug)", tc.val, tc.expected, coeff)
		}
	}
}

func TestNumericTypesInSum(t *testing.T) {
	testCases := []any{
		int(10), int64(10), int32(10), int16(10), int8(10),
		uint(10), uint64(10), uint32(10), uint16(10), uint8(10),
		float64(10.0), float32(10.0),
	}

	for _, tc := range testCases {
		expr := Sum(tc)
		if expr.Constant() != 10.0 {
			t.Errorf("For type %T in Sum: expected constant 10.0, got %g", tc, expr.Constant())
		}
	}
}

func TestNumericTypesInConstraintRHS(t *testing.T) {
	m := New("model")
	x := m.Binary("x")

	testCases := []any{
		int(3), int64(3), int32(3), int16(3), int8(3),
		uint(3), uint64(3), uint32(3), uint16(3), uint8(3),
		float64(3.0), float32(3.0),
	}

	for _, tc := range testCases {
		// This should not panic and should construct a constraint with RHS = 3.0
		var c *Constraint
		err := capturePanic(func() {
			c = x.Eq(tc)
		})
		if err != nil {
			t.Errorf("For type %T in constraint RHS: Eq() panicked: %v", tc, err)
			continue
		}
		if c.RHS() != 3.0 {
			t.Errorf("For type %T in constraint RHS: expected RHS 3.0, got %g", tc, c.RHS())
		}
	}
}

func TestInvalidTypesPanic(t *testing.T) {
	m := New("model")
	x := m.Binary("x")

	// Invalid coefficient in Lin
	if err := capturePanic(func() { Lin("invalid", x) }); err == nil {
		t.Error("Expected panic for invalid string coefficient in Lin")
	}

	// Invalid argument in Sum
	if err := capturePanic(func() { Sum("invalid") }); err == nil {
		t.Error("Expected panic for invalid string argument in Sum")
	}

	// Invalid RHS in Eq
	if err := capturePanic(func() { x.Eq("invalid") }); err == nil {
		t.Error("Expected panic for invalid string RHS in Eq")
	}
}

func capturePanic(f func()) (err any) {
	defer func() {
		err = recover()
	}()
	f()
	return nil
}
