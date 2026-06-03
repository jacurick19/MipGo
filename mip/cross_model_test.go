package mip_test

import (
	"testing"

	. "mipgo/v2/mip"
)

func TestCrossModelVariablePanic(t *testing.T) {
	m1 := New("model1")
	m2 := New("model2")

	x1 := m1.Binary("x1")
	x2 := m2.Binary("x2")

	// Test case 1: Constraint with external variable in SubjectTo
	assertPanic(t, "SubjectTo with external variable", func() {
		m1.SubjectTo(x2.Leq(1.0))
	})

	// Test case 2: External variable in Minimize
	assertPanic(t, "Minimize with external variable", func() {
		m1.Minimize(x2)
	})

	// Test case 3: External variable in Maximize
	assertPanic(t, "Maximize with external variable", func() {
		m1.Maximize(x2)
	})

	// Test case 4: External variable in AddSOS1
	assertPanic(t, "AddSOS1 with external variable", func() {
		m1.AddSOS1(x1, x2)
	})

	// Test case 5: External variable in AddSOS2
	assertPanic(t, "AddSOS2 with external variable", func() {
		m1.AddSOS2(x1, x2)
	})

	// Test case 6: External variable in Indicator
	assertPanic(t, "Indicator binary with external variable", func() {
		m1.Indicator(x2, x1.Leq(1.0))
	})
	assertPanic(t, "Indicator constraint with external variable", func() {
		m1.Indicator(x1, x2.Leq(1.0))
	})

	// Test case 7: External variable in SetMIPStart
	assertPanic(t, "SetMIPStart with external variable", func() {
		m1.SetMIPStart(map[*Variable]float64{
			x2: 1.0,
		})
	})
}

func assertPanic(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for '%s', but function completed without panic", name)
		}
	}()
	f()
}
