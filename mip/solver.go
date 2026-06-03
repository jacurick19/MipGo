package mip

import (
	"errors"
	"math"
)

// SolverBackend defines the interface that solver drivers must implement.
type SolverBackend interface {
	Solve(model *Model) (*Solution, error)
}

var backends = make(map[SolverType]SolverBackend)

// RegisterBackend registers a solver backend for a specific SolverType.
// This is typically called from a driver's init() function.
func RegisterBackend(t SolverType, backend SolverBackend) {
	backends[t] = backend
}

// NativeCallbackSupport is an optional interface that solver backends can implement
// to declare native support for branch-and-cut callbacks.
type NativeCallbackSupport interface {
	SupportsNativeCallbacks() bool
}

// Solve solves the model using the selected solver backend.
func (m *Model) Solve() (*Solution, error) {
	backend, ok := backends[m.solverType]
	if !ok {
		return nil, errors.New("requested solver backend is not registered")
	}

	nativeSupport, isNative := backend.(NativeCallbackSupport)
	hasNative := isNative && nativeSupport.SupportsNativeCallbacks()

	if m.lazyCallback != nil && !hasNative {
		for {
			sol, err := backend.Solve(m)
			if err != nil {
				return nil, err
			}
			if sol.Status() != Optimal && sol.Status() != Feasible {
				return sol, nil
			}

			violated := m.lazyCallback(sol)
			if len(violated) == 0 {
				return sol, nil
			}

			// Add violated constraints back into the model and resolve
			m.SubjectTo(violated...)
		}
	}

	return backend.Solve(m)
}

// RegisteredSolvers returns a slice of all registered solver types in the current build.
func RegisteredSolvers() []SolverType {
	var list []SolverType
	for t := range backends {
		list = append(list, t)
	}
	return list
}

// GetBigM calculates a safe Big-M value for indicator constraint linearization based on variable bounds.
func GetBigM(c *Constraint) float64 {
	var boundVal float64 = 0.0
	if c.sense == 'L' {
		for _, t := range c.termsList {
			if t.Coeff > 0 {
				if !math.IsInf(t.Var.ub, 1) {
					boundVal += t.Coeff * t.Var.ub
				} else {
					return 1e5
				}
			} else {
				if !math.IsInf(t.Var.lb, -1) {
					boundVal += t.Coeff * t.Var.lb
				} else {
					return 1e5
				}
			}
		}
		mVal := boundVal - c.rhs
		if mVal < 1e5 {
			return 1e5
		}
		return mVal
	} else {
		for _, t := range c.termsList {
			if t.Coeff > 0 {
				if !math.IsInf(t.Var.lb, -1) {
					boundVal += t.Coeff * t.Var.lb
				} else {
					return 1e5
				}
			} else {
				if !math.IsInf(t.Var.ub, 1) {
					boundVal += t.Coeff * t.Var.ub
				} else {
					return 1e5
				}
			}
		}
		mVal := c.rhs - boundVal
		if mVal < 1e5 {
			return 1e5
		}
		return mVal
	}
}

// ModelWriter defines an optional interface for solver backends to write models.
type ModelWriter interface {
	WriteLP(model *Model, filename string) error
	WriteMPS(model *Model, filename string) error
}

// ModelReader defines an optional interface for solver backends to read models.
type ModelReader interface {
	ReadLP(filename string) (*Model, error)
	ReadMPS(filename string) (*Model, error)
}

func (m *Model) WriteLP(filename string) error {
	for _, backend := range backends {
		if writer, ok := backend.(ModelWriter); ok {
			return writer.WriteLP(m, filename)
		}
	}
	return errors.New("no registered solver backend supports writing LP files")
}

func (m *Model) WriteMPS(filename string) error {
	for _, backend := range backends {
		if writer, ok := backend.(ModelWriter); ok {
			return writer.WriteMPS(m, filename)
		}
	}
	return errors.New("no registered solver backend supports writing MPS files")
}

func ReadLP(filename string) (*Model, error) {
	for _, backend := range backends {
		if reader, ok := backend.(ModelReader); ok {
			return reader.ReadLP(filename)
		}
	}
	return nil, errors.New("no registered solver backend supports reading LP files")
}

func ReadMPS(filename string) (*Model, error) {
	for _, backend := range backends {
		if reader, ok := backend.(ModelReader); ok {
			return reader.ReadMPS(filename)
		}
	}
	return nil, errors.New("no registered solver backend supports reading MPS files")
}
