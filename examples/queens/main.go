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
		log.Fatal("No solver backends registered.")
	}

	n := 40

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving %d-Queens with %s...\n", n, sName)
		fmt.Printf("========================================\n")

		queens := mip.New("n-Queens")
		queens.SetSolver(s)

		// Binary matrix x[i][j]
		x := queens.BinaryMatrix("x", n, n)

		// One per row
		for i := 0; i < n; i++ {
			var rowVars []*mip.Variable
			for j := 0; j < n; j++ {
				rowVars = append(rowVars, x.At(i, j))
			}
			queens.SubjectTo(mip.SumVars(rowVars...).Eq(1.0).Named(fmt.Sprintf("row(%d)", i)))
		}

		// One per column
		for j := 0; j < n; j++ {
			var colVars []*mip.Variable
			for i := 0; i < n; i++ {
				colVars = append(colVars, x.At(i, j))
			}
			queens.SubjectTo(mip.SumVars(colVars...).Eq(1.0).Named(fmt.Sprintf("col(%d)", j)))
		}

		// Diagonal \ (i - j = k)
		p := 0
		for k := 2 - n; k <= n-2; k++ {
			var diagVars []*mip.Variable
			for i := 0; i < n; i++ {
				j := i - k
				if j >= 0 && j < n {
					diagVars = append(diagVars, x.At(i, j))
				}
			}
			if len(diagVars) > 0 {
				queens.SubjectTo(mip.SumVars(diagVars...).Leq(1.0).Named(fmt.Sprintf("diag1(%d)", p)))
				p++
			}
		}

		// Diagonal / (i + j = k)
		p = 0
		for k := 3; k < n+n; k++ {
			var antiDiagVars []*mip.Variable
			for i := 0; i < n; i++ {
				j := k - i
				if j >= 0 && j < n {
					antiDiagVars = append(antiDiagVars, x.At(i, j))
				}
			}
			if len(antiDiagVars) > 0 {
				queens.SubjectTo(mip.SumVars(antiDiagVars...).Leq(1.0).Named(fmt.Sprintf("diag2(%d)", p)))
				p++
			}
		}

		sol, err := queens.Solve()
		if err != nil {
			log.Printf("Failed to solve n-Queens with %s: %v\n", sName, err)
			continue
		}

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			fmt.Println("Board layout:")
			vars := queens.Variables()
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					v := vars[i*n+j]
					if sol.Value(v) >= 0.99 {
						fmt.Print("Q ")
					} else {
						fmt.Print(". ")
					}
				}
				fmt.Println()
			}
			fmt.Println()
		} else {
			fmt.Println("No feasible solution found.")
			fmt.Println()
		}
	}
}

func solverName(s mip.SolverType) string {
	switch s {
	case mip.HiGHS:
		return "HiGHS"
	case mip.CBC:
		return "CBC"
	case mip.SCIP:
		return "SCIP"
	default:
		return "Unknown"
	}
}
