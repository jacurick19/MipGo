package mip

import (
	"testing"
)

func TestCallbackRegistration(t *testing.T) {
	m := New("test_callbacks")

	lazyCalled := false
	lazyCb := func(sol *Solution) []*Constraint {
		lazyCalled = true
		return nil
	}

	cutCalled := false
	cutCb := func(sol *Solution) []*Constraint {
		cutCalled = true
		return nil
	}

	m.AddLazyConstraintCallback(lazyCb)
	m.AddCutCallback(cutCb)

	if m.lazyCallback == nil {
		t.Errorf("Expected lazyCallback to be registered")
	}
	if m.cutCallback == nil {
		t.Errorf("Expected cutCallback to be registered")
	}

	// Invoke to verify
	m.lazyCallback(nil)
	m.cutCallback(nil)

	if !lazyCalled || !cutCalled {
		t.Errorf("Registered callback closures were not executed correctly")
	}
}


