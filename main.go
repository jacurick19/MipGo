package main

import (
	"fmt"
	"log"

	"mipgo/v2/mip"
	_ "mipgo/v2/mip/cbc"
	_ "mipgo/v2/mip/highs"
	_ "mipgo/v2/mip/scip"
)

func main() {
	solvers := mip.RegisteredSolvers()
	if len(solvers) == 0 {
		log.Fatal("No solver backends compiled in. Please build with -tags highs or -tags cbc.")
	}

	fmt.Printf("Detected %d registered solver backend(s):\n", len(solvers))
	for _, s := range solvers {
		var name string
		switch s {
		case mip.HiGHS:
			name = "HiGHS"
		case mip.CBC:
			name = "CBC"
		case mip.Gurobi:
			name = "Gurobi"
		case mip.SCIP:
			name = "SCIP"
		}
		fmt.Printf(" - %s\n", name)
	}
	fmt.Println()

	n := 8

	// Run N-Queens for each registered solver
	for _, s := range solvers {
		var sName string
		switch s {
		case mip.HiGHS:
			sName = "HiGHS"
		case mip.CBC:
			sName = "CBC"
		case mip.SCIP:
			sName = "SCIP"
		default:
			sName = "Unknown"
		}

		fmt.Printf("========================================\n")
		fmt.Printf("Solving %d-Queens with %s...\n", n, sName)
		fmt.Printf("========================================\n")

		// Rebuild model per solver to ensure clean variables references
		m := mip.New("N-Queens")
		m.SetSolver(s)

		status, sol := solveNQueens(n, m)

		if status == mip.Optimal {
			fmt.Printf("Success! %s found optimal solution.\n", sName)
			fmt.Printf("Objective value: %g\n", sol.Objective())
			fmt.Printf("Solver Runtime: %v\n", sol.Stats().Runtime)
			fmt.Println("Board Layout:")
			vars := m.Variables()
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					v := vars[i*n+j]
					if sol.Value(v) > 0.5 {
						fmt.Print("Q ")
					} else {
						fmt.Print(". ")
					}
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("Failed to find optimal solution using %s. Status: %v\n", sName, status)
		}
	}
}

func buildNQueens(n int, model *mip.Model) *mip.VarMatrix {
	// Create Binary Matrix
	x := model.BinaryMatrix("x", n, n)

	// 1. Row constraints: exactly one queen per row
	for i := 0; i < n; i++ {
		model.SubjectTo(
			mip.SumOver(n, func(j int) mip.Expr {
				return x.At(i, j)
			}).Eq(1.0).Named(fmt.Sprintf("row_%d", i)),
		)
	}

	// 2. Column constraints: exactly one queen per column
	for j := 0; j < n; j++ {
		model.SubjectTo(
			mip.SumOver(n, func(i int) mip.Expr {
				return x.At(i, j)
			}).Eq(1.0).Named(fmt.Sprintf("col_%d", j)),
		)
	}

	// 3. Diagonal constraints (i - j = k): at most one queen per diagonal
	for k := -n + 1; k < n; k++ {
		var diagVars []*mip.Variable
		for i := 0; i < n; i++ {
			j := i - k
			if j >= 0 && j < n {
				diagVars = append(diagVars, x.At(i, j))
			}
		}
		if len(diagVars) > 1 {
			model.SubjectTo(
				mip.SumVars(diagVars...).Leq(1.0).Named(fmt.Sprintf("diag_%d", k)),
			)
		}
	}

	// 4. Anti-diagonal constraints (i + j = k): at most one queen per anti-diagonal
	for k := 0; k < 2*n-1; k++ {
		var antiDiagVars []*mip.Variable
		for i := 0; i < n; i++ {
			j := k - i
			if j >= 0 && j < n {
				antiDiagVars = append(antiDiagVars, x.At(i, j))
			}
		}
		if len(antiDiagVars) > 1 {
			model.SubjectTo(
				mip.SumVars(antiDiagVars...).Leq(1.0).Named(fmt.Sprintf("antidiag_%d", k)),
			)
		}
	}

	return x
}

func solveNQueens(n int, model *mip.Model) (mip.Status, *mip.Solution) {
	buildNQueens(n, model)
	sol, err := model.Solve()
	if err != nil {
		log.Fatalf("Error solving model: %v", err)
	}
	return sol.Status(), sol
}
