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

	w := []int{4, 3, 5, 2, 1, 4, 7, 3} // widths
	h := []int{2, 4, 1, 5, 6, 3, 5, 4} // heights
	n := len(w)

	W := 10 // raw material width

	// Build sets S and G
	S := make([][]int, n)
	G := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if h[j] <= h[i] {
				S[i] = append(S[i], j)
			}
			if h[j] >= h[i] {
				G[i] = append(G[i], j)
			}
		}
	}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving Level Packing with %s...\n", sName)
		fmt.Printf("========================================\n")

		model := mip.New("LevelPacking")
		model.SetSolver(s)

		// x[i][j] is defined only for j in S[i]
		x := make([]map[int]*mip.Variable, n)
		for i := 0; i < n; i++ {
			x[i] = make(map[int]*mip.Variable)
			for _, j := range S[i] {
				x[i][j] = model.Binary(fmt.Sprintf("x_%d_%d", i, j))
			}
		}

		// Minimize sum of h[i] * x[i][i]
		var objTerms []mip.Term
		for i := 0; i < n; i++ {
			objTerms = append(objTerms, mip.Term{Coeff: float64(h[i]), Var: x[i][i]})
		}
		model.Minimize(mip.NewExpressionFromReader(objTerms, 0.0))

		// Constraint: each item i belongs to exactly one level
		for i := 0; i < n; i++ {
			var list []*mip.Variable
			for _, j := range G[i] {
				list = append(list, x[j][i])
			}
			model.SubjectTo(mip.SumVars(list...).Eq(1.0))
		}

		// Constraint: remaining width capacity for level i
		for i := 0; i < n; i++ {
			var terms []mip.Term
			for _, j := range S[i] {
				if j != i {
					terms = append(terms, mip.Term{Coeff: float64(w[j]), Var: x[i][j]})
				}
			}
			coeff := -float64(W - w[i])
			terms = append(terms, mip.Term{Coeff: coeff, Var: x[i][i]})
			model.SubjectTo(mip.NewExpressionFromReader(terms, 0.0).Leq(0.0))
		}

		sol, err := model.Solve()
		if err != nil {
			log.Printf("Failed to solve Level Packing with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value: %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			for i := 0; i < n; i++ {
				if sol.Value(x[i][i]) >= 0.99 {
					var grouped []int
					for _, j := range S[i] {
						if j != i && sol.Value(x[i][j]) >= 0.99 {
							grouped = append(grouped, j)
						}
					}
					fmt.Printf("Items grouped with %d: %v\n", i, grouped)
				}
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
