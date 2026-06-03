package mip

import "sync"

var (
	callbackMu        sync.RWMutex
	callbackNextID    int
	activeSolversVars = make(map[int][]*Variable)
	activeModels      = make(map[int]*Model)
)

func RegisterHiGHSModel(m *Model) int {
	callbackMu.Lock()
	defer callbackMu.Unlock()
	id := callbackNextID
	callbackNextID++
	activeSolversVars[id] = m.variables
	activeModels[id] = m
	return id
}

func UnregisterHiGHSModel(id int) {
	callbackMu.Lock()
	defer callbackMu.Unlock()
	delete(activeSolversVars, id)
	delete(activeModels, id)
}

func GetActiveVars(id int) ([]*Variable, bool) {
	callbackMu.RLock()
	defer callbackMu.RUnlock()
	vars, ok := activeSolversVars[id]
	return vars, ok
}

func AppendSolutionToPool(id int, sol *Solution) {
	callbackMu.Lock()
	defer callbackMu.Unlock()
	m, ok := activeModels[id]
	if !ok {
		return
	}
	// Verify if we already have this solution to avoid duplicates
	for _, existing := range m.solutionPool {
		if mathAbs(existing.objective-sol.objective) < 1e-9 {
			// Check if variable values are also identical
			match := true
			for _, v := range m.variables {
				if mathAbs(existing.values[v]-sol.values[v]) > 1e-6 {
					match = false
					break
				}
			}
			if match {
				return // Already captured
			}
		}
	}
	m.solutionPool = append(m.solutionPool, sol)
}

func mathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
