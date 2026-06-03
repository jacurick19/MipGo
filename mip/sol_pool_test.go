package mip_test

import (
	"testing"

	. "mipgo/v2/mip"
	_ "mipgo/v2/mip/cbc"
	_ "mipgo/v2/mip/highs"
	_ "mipgo/v2/mip/scip"
)

func TestSolutionPoolIntegration(t *testing.T) {
	solvers := RegisteredSolvers()
	if len(solvers) == 0 {
		t.Skip("No solver backends compiled in. Skipping integration test.")
	}

	for _, s := range solvers {
		sName := "Unknown"
		switch s {
		case HiGHS:
			sName = "HiGHS"
		case CBC:
			sName = "CBC"
		case SCIP:
			sName = "SCIP"
		}

		t.Run(sName, func(t *testing.T) {
			m := New("Solution-Pool-Test")
			m.SetSolver(s)

			// 3 binary variables: x0, x1, x2
			x := m.BinaryVars("x", 3)

			// Maximize 2*x0 + 1.5*x1 + 0.5*x2
			m.Maximize(Sum(Prod(2.0, x[0]), Prod(1.5, x[1]), Prod(0.5, x[2])))

			// Constraint: 1.5*x0 + 1.5*x1 + 1.0*x2 <= 2.0
			m.SubjectTo(Sum(Prod(1.5, x[0]), Prod(1.5, x[1]), Prod(1.0, x[2])).Leq(2.0))

			sol, err := m.Solve()
			if err != nil {
				t.Fatalf("Solver %s failed with error: %v", sName, err)
			}

			if sol.Status() != Optimal {
				t.Fatalf("Expected status Optimal, got %v", sol.Status())
			}

			pool := m.SolutionPool()
			if len(pool) == 0 {
				t.Fatalf("Expected solution pool to be populated, got 0 solutions")
			}

			t.Logf("Solver %s found %d solution(s) in the pool:", sName, len(pool))
			for i, pSol := range pool {
				t.Logf(" Solution %d: objective = %g, values = [x0=%g, x1=%g, x2=%g]",
					i, pSol.Objective(), pSol.Value(x[0]), pSol.Value(x[1]), pSol.Value(x[2]))
			}

			// The best solution (index 0) must match the returned solution
			if pool[0].Objective() != sol.Objective() {
				t.Errorf("Expected best solution in pool to have objective %g, got %g", sol.Objective(), pool[0].Objective())
			}

			// If the solver supports alternative solutions and found them, verify uniqueness
			seen := make(map[string]bool)
			for _, pSol := range pool {
				key := ""
				for _, v := range x {
					val := 0.0
					if pSol.Value(v) > 0.5 {
						val = 1.0
					}
					key += string(rune('0' + int(val)))
				}
				seen[key] = true
			}

			// HiGHS and CBC may only return 1 solution if they solve the problem quickly or don't collect alternatives,
			// but SCIP should find at least 2 distinct feasible/optimal solutions.
			if s == SCIP && len(seen) < 2 {
				t.Errorf("Expected SCIP solver pool to contain multiple distinct solutions, got only %d", len(seen))
			}
		})
	}
}
