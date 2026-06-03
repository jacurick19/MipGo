package mip_test

import (
	"fmt"
	"testing"

	. "mipgo/v2/mip"
	_ "mipgo/v2/mip/cbc"
	_ "mipgo/v2/mip/highs"
	_ "mipgo/v2/mip/scip"
)

func TestMIPStartIntegration(t *testing.T) {
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
			m := New("MIPStart-Test")
			m.SetSolver(s)

			// 3 binary variables: x0, x1, x2
			x := m.BinaryVars("x", 3)

			// Maximize 2*x0 + 1.5*x1 + 0.5*x2
			m.Maximize(Sum(Prod(2.0, x[0]), Prod(1.5, x[1]), Prod(0.5, x[2])))

			// Constraint: 1.5*x0 + 1.5*x1 + 1.0*x2 <= 2.0
			m.SubjectTo(Sum(Prod(1.5, x[0]), Prod(1.5, x[1]), Prod(1.0, x[2])).Leq(2.0))

			// We set a MIPStart (initial feasible solution):
			// x0 = 1, x1 = 0, x2 = 0 (objective = 2) which is optimal and feasible.
			startVal := map[*Variable]float64{
				x[0]: 1.0,
				x[1]: 0.0,
				x[2]: 0.0,
			}
			m.SetMIPStart(startVal)

			sol, err := m.Solve()
			if err != nil {
				t.Fatalf("Solver %s failed with error: %v", sName, err)
			}

			if sol.Status() != Optimal {
				t.Fatalf("Expected status Optimal, got %v", sol.Status())
			}

			if sol.Objective() != 2.0 {
				t.Errorf("Expected objective 2.0, got %g", sol.Objective())
			}

			if sol.Value(x[0]) != 1.0 || sol.Value(x[1]) != 0.0 || sol.Value(x[2]) != 0.0 {
				t.Errorf("Expected solution x=[1, 0, 0], got x=[%g, %g, %g]",
					sol.Value(x[0]), sol.Value(x[1]), sol.Value(x[2]))
			}
		})
	}
}

func TestNRooksMIPStart(t *testing.T) {
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
			m := New("NRooks-MIPStart-Test")
			m.SetSolver(s)

			n := 6
			x := make([][]*Variable, n)
			for i := 0; i < n; i++ {
				x[i] = m.BinaryVars(fmt.Sprintf("x_%d", i), n)
			}

			// Constraints: Row sums = 1
			for i := 0; i < n; i++ {
				var rowVars []any
				for j := 0; j < n; j++ {
					rowVars = append(rowVars, x[i][j])
				}
				m.SubjectTo(Sum(rowVars...).Eq(1.0))
			}

			// Constraints: Column sums = 1
			for j := 0; j < n; j++ {
				var colVars []any
				for i := 0; i < n; i++ {
					colVars = append(colVars, x[i][j])
				}
				m.SubjectTo(Sum(colVars...).Eq(1.0))
			}

			// Objective: Maximize 10 * sum(top half diagonal) + sum(others)
			// This forces the unique optimal solution to include the top half diagonal variables.
			var terms []any
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if i < n/2 && i == j {
						terms = append(terms, Prod(10.0, x[i][j]))
					} else {
						terms = append(terms, Prod(1.0, x[i][j]))
					}
				}
			}
			m.Maximize(Sum(terms...))

			// MIPStart: top half diagonal filled (x[0][0], x[1][1], x[2][2])
			// and all other variables in their rows/columns set to 0.0,
			// leaving the rest of the chessboard (i,j >= 3) unspecified.
			startVal := make(map[*Variable]float64)
			for i := 0; i < n/2; i++ {
				for j := 0; j < n; j++ {
					if i == j {
						startVal[x[i][j]] = 1.0
					} else {
						startVal[x[i][j]] = 0.0
						startVal[x[j][i]] = 0.0
					}
				}
			}
			m.SetMIPStart(startVal)

			sol, err := m.Solve()
			if err != nil {
				t.Fatalf("Solver %s failed with error: %v", sName, err)
			}

			if sol.Status() != Optimal && sol.Status() != Feasible {
				t.Fatalf("Expected status Optimal or Feasible, got %v", sol.Status())
			}

			// Verify row sums
			for i := 0; i < n; i++ {
				sum := 0.0
				for j := 0; j < n; j++ {
					if sol.Value(x[i][j]) > 0.5 {
						sum += 1.0
					}
				}
				if sum != 1.0 {
					t.Errorf("Row %d sum is %g, expected 1.0", i, sum)
				}
			}

			// Verify column sums
			for j := 0; j < n; j++ {
				sum := 0.0
				for i := 0; i < n; i++ {
					if sol.Value(x[i][j]) > 0.5 {
						sum += 1.0
					}
				}
				if sum != 1.0 {
					t.Errorf("Column %d sum is %g, expected 1.0", j, sum)
				}
			}

			// Verify that the MIPStart choices are respected
			for i := 0; i < n/2; i++ {
				if sol.Value(x[i][i]) < 0.5 {
					t.Errorf("Expected MIPStart diagonal variable x[%d][%d] to be 1.0, got %g", i, i, sol.Value(x[i][i]))
				}
			}
		})
	}
}
