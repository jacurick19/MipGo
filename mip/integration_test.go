package mip_test

import (
	"fmt"
	"math"
	"testing"

	. "mipgo/v2/mip"
	_ "mipgo/v2/mip/cbc"
	_ "mipgo/v2/mip/highs"
	_ "mipgo/v2/mip/scip"
)

func TestNQueensSolversIntegration(t *testing.T) {
	solvers := RegisteredSolvers()
	if len(solvers) == 0 {
		t.Skip("No solver backends compiled in. Skipping integration test.")
	}

	// Test cases for different board sizes (N)
	tests := []struct {
		n              int
		expectedStatus Status
	}{
		{n: 1, expectedStatus: Optimal},
		{n: 2, expectedStatus: Infeasible},
		{n: 3, expectedStatus: Infeasible},
		{n: 4, expectedStatus: Optimal},
		{n: 8, expectedStatus: Optimal},
	}

	for _, s := range solvers {
		sName := "Unknown"
		switch s {
		case HiGHS:
			sName = "HiGHS"
		case CBC:
			sName = "CBC"
		case Gurobi:
			sName = "Gurobi"
		case SCIP:
			sName = "SCIP"
		}

		t.Run(sName, func(t *testing.T) {
			for _, tc := range tests {
				t.Run(fmt.Sprintf("N=%d", tc.n), func(t *testing.T) {
					model := New(fmt.Sprintf("N-Queens-Integration-%d", tc.n))
					model.SetSolver(s)

					// 1. Build N-Queens formulation
					x := model.BinaryMatrix("x", tc.n, tc.n)

					// Row constraints: exactly one queen per row
					for i := 0; i < tc.n; i++ {
						model.SubjectTo(
							SumOver(tc.n, func(j int) Expr {
								return x.At(i, j)
							}).Eq(1.0).Named(fmt.Sprintf("row_%d", i)),
						)
					}

					// Column constraints: exactly one queen per column
					for j := 0; j < tc.n; j++ {
						model.SubjectTo(
							SumOver(tc.n, func(i int) Expr {
								return x.At(i, j)
							}).Eq(1.0).Named(fmt.Sprintf("col_%d", j)),
						)
					}

					// Diagonal constraints (i - j = k): at most one queen per diagonal
					for k := -tc.n + 1; k < tc.n; k++ {
						var diagVars []*Variable
						for i := 0; i < tc.n; i++ {
							j := i - k
							if j >= 0 && j < tc.n {
								diagVars = append(diagVars, x.At(i, j))
							}
						}
						if len(diagVars) > 1 {
							model.SubjectTo(
								SumVars(diagVars...).Leq(1.0).Named(fmt.Sprintf("diag_%d", k)),
							)
						}
					}

					// Anti-diagonal constraints (i + j = k): at most one queen per anti-diagonal
					for k := 0; k < 2*tc.n-1; k++ {
						var antiDiagVars []*Variable
						for i := 0; i < tc.n; i++ {
							j := k - i
							if j >= 0 && j < tc.n {
								antiDiagVars = append(antiDiagVars, x.At(i, j))
							}
						}
						if len(antiDiagVars) > 1 {
							model.SubjectTo(
								SumVars(antiDiagVars...).Leq(1.0).Named(fmt.Sprintf("antidiag_%d", k)),
							)
						}
					}

					// 2. Solve the model
					sol, err := model.Solve()
					if err != nil {
						t.Fatalf("Solver %s failed with error: %v", sName, err)
					}

					if sol.Status() != tc.expectedStatus {
						t.Fatalf("Solver %s on N=%d expected status %v, got %v", sName, tc.n, tc.expectedStatus, sol.Status())
					}

					if tc.expectedStatus == Optimal {
						// 3. Verify Solution Validity for feasible boards
						type pos struct{ r, c int }
						var queens []pos

						for i := 0; i < tc.n; i++ {
							for j := 0; j < tc.n; j++ {
								val := sol.Value(x.At(i, j))
								if val > 0.5 {
									queens = append(queens, pos{i, j})
								}
							}
						}

						// Check total queens count
						if len(queens) != tc.n {
							t.Errorf("Solver %s placed %d queens, expected %d", sName, len(queens), tc.n)
						}

						// Verify that no two queens attack each other
						for a := 0; a < len(queens); a++ {
							for b := a + 1; b < len(queens); b++ {
								qa := queens[a]
								qb := queens[b]

								// Check same row
								if qa.r == qb.r {
									t.Errorf("Solver %s conflict: Queens at (%d,%d) and (%d,%d) share row %d", sName, qa.r, qa.c, qb.r, qb.c, qa.r)
								}

								// Check same column
								if qa.c == qb.c {
									t.Errorf("Solver %s conflict: Queens at (%d,%d) and (%d,%d) share column %d", sName, qa.r, qa.c, qb.r, qb.c, qa.c)
								}

								// Check diagonal
								if math.Abs(float64(qa.r-qb.r)) == math.Abs(float64(qa.c-qb.c)) {
									t.Errorf("Solver %s conflict: Queens at (%d,%d) and (%d,%d) share diagonal", sName, qa.r, qa.c, qb.r, qb.c)
								}
							}
						}
					}
				})
			}
		})
	}
}

func TestLazyConstraintsOuterLoop(t *testing.T) {
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
			m := New("Lazy-Constraints-Outer-Loop-Test")
			m.SetSolver(s)

			// 3 binary variables: x0, x1, x2
			x := m.BinaryVars("x", 3)

			// Maximize x0 + x1 + x2
			m.Maximize(SumVars(x...))

			// Without lazy constraints, the optimal solution would be x0=1, x1=1, x2=1 (obj = 3)
			// We add a lazy constraint: x0 + x1 + x2 <= 2
			callbackInvocations := 0
			m.AddLazyConstraintCallback(func(sol *Solution) []*Constraint {
				callbackInvocations++
				valSum := 0.0
				for _, v := range x {
					valSum += sol.Value(v)
				}
				// If solution sum is greater than 2, return violated constraint
				if valSum > 2.01 {
					var terms []any
					for _, v := range x {
						terms = append(terms, 1.0, v)
					}
					// Return constraint: x0 + x1 + x2 <= 2.0
					return []*Constraint{
						Lin(terms...).Leq(2.0),
					}
				}
				return nil
			})

			sol, err := m.Solve()
			if err != nil {
				t.Fatalf("Solver %s failed with error: %v", sName, err)
			}

			if sol.Status() != Optimal {
				t.Fatalf("Expected status Optimal, got %v", sol.Status())
			}

			if sol.Objective() != 2.0 {
				t.Errorf("Expected objective value 2.0, got %g (Callback invocations: %d)", sol.Objective(), callbackInvocations)
			}

			// Verify values sum to 2
			valSum := 0.0
			for _, v := range x {
				valSum += sol.Value(v)
			}
			if math.Abs(valSum-2.0) > 0.01 {
				t.Errorf("Expected solution variables to sum to 2.0, got %g", valSum)
			}
		})
	}
}

