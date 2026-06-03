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

	n := 10 // exactly 12 jobs (n=10 plus 2 dummies)

	p := []int{0, 3, 2, 5, 4, 2, 3, 4, 2, 4, 6, 0}

	u := [][]int{
		{0, 0}, {5, 1}, {0, 4}, {1, 4}, {1, 3}, {3, 2}, {3, 1}, {2, 4},
		{4, 0}, {5, 2}, {2, 5}, {0, 0},
	}

	c := []int{6, 8}

	S := [][]int{
		{0, 1}, {0, 2}, {0, 3}, {1, 4}, {1, 5}, {2, 9}, {2, 10}, {3, 8}, {4, 6},
		{4, 7}, {5, 9}, {5, 10}, {6, 8}, {6, 9}, {7, 8}, {8, 11}, {9, 11}, {10, 11},
	}

	sumP := 0
	for _, val := range p {
		sumP += val
	}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving RCPSP with %s...\n", sName)
		fmt.Printf("========================================\n")

		m := mip.New("RCPSP")
		m.SetSolver(s)

		// Binary variables x[j][t]
		x := m.BinaryMatrix("x", len(p), sumP)

		// Objective: Minimize sum(t * x[n+1][t])
		var objTerms []mip.Term
		for t := 0; t < sumP; t++ {
			objTerms = append(objTerms, mip.Term{Coeff: float64(t), Var: x.At(n+1, t)})
		}
		m.Minimize(mip.NewExpressionFromReader(objTerms, 0.0))

		// Constraint: Each job starts exactly once
		for j := 0; j < len(p); j++ {
			var rowVars []*mip.Variable
			for t := 0; t < sumP; t++ {
				rowVars = append(rowVars, x.At(j, t))
			}
			m.SubjectTo(mip.SumVars(rowVars...).Eq(1.0))
		}

		// Constraint: Resource capacity
		for r := 0; r < len(c); r++ {
			for t := 0; t < sumP; t++ {
				var terms []mip.Term
				for j := 0; j < len(p); j++ {
					startT := t - p[j] + 1
					if startT < 0 {
						startT = 0
					}
					for t2 := startT; t2 <= t; t2++ {
						if u[j][r] > 0 {
							terms = append(terms, mip.Term{Coeff: float64(u[j][r]), Var: x.At(j, t2)})
						}
					}
				}
				m.SubjectTo(mip.NewExpressionFromReader(terms, 0.0).Leq(float64(c[r])))
			}
		}

		// Constraint: Precedence relations
		for _, prec := range S {
			jVal := prec[0]
			sVal := prec[1]
			var terms []mip.Term
			for t := 0; t < sumP; t++ {
				terms = append(terms, mip.Term{Coeff: float64(t), Var: x.At(sVal, t)})
				terms = append(terms, mip.Term{Coeff: -float64(t), Var: x.At(jVal, t)})
			}
			m.SubjectTo(mip.NewExpressionFromReader(terms, 0.0).Geq(float64(p[jVal])))
		}

		sol, err := m.Solve()
		if err != nil {
			log.Printf("Failed to solve RCPSP with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value: %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			fmt.Println("Schedule:")
			for j := 0; j < len(p); j++ {
				for t := 0; t < sumP; t++ {
					if sol.Value(x.At(j, t)) >= 0.99 {
						fmt.Printf("Job %d: begins at t=%d and finishes at t=%d\n", j, t, t+p[j])
					}
				}
			}
			fmt.Printf("Makespan = %g\n\n", sol.Objective())
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
