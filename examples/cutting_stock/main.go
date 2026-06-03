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

	n := 10  // maximum number of bars
	L := 250 // bar length
	m := 4   // number of requests
	w := []int{187, 119, 74, 90}
	b := []int{1, 2, 2, 1}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving Cutting Stock with %s...\n", sName)
		fmt.Printf("========================================\n")

		model := mip.New("CuttingStock")
		model.SetSolver(s)

		// Binary variables y[j]
		y := model.BinaryVars("y", n)

		// Integer variables x[i][j] (how many of request i in bar j)
		x := model.IntegerMatrix("x", m, n, 0, mip.Inf)

		// Minimize sum of y[j]
		model.Minimize(mip.SumVars(y...))

		// Constraint: Satisfy all demands
		for i := 0; i < m; i++ {
			var rowVars []*mip.Variable
			for j := 0; j < n; j++ {
				rowVars = append(rowVars, x.At(i, j))
			}
			model.SubjectTo(mip.SumVars(rowVars...).Geq(float64(b[i])))
		}

		// Constraint: Capacity of each bar
		for j := 0; j < n; j++ {
			var terms []mip.Term
			for i := 0; i < m; i++ {
				terms = append(terms, mip.Term{Coeff: float64(w[i]), Var: x.At(i, j)})
			}
			terms = append(terms, mip.Term{Coeff: -float64(L), Var: y[j]})
			model.SubjectTo(mip.NewExpressionFromReader(terms, 0.0).Leq(0.0))
		}

		// Symmetry reduction: y[j-1] >= y[j]
		for j := 1; j < n; j++ {
			model.SubjectTo(mip.Sum(y[j-1], mip.Neg(y[j])).Geq(0.0))
		}

		sol, err := model.Solve()
		if err != nil {
			log.Printf("Failed to solve Cutting Stock with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value: %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			for j := 0; j < n; j++ {
				if sol.Value(y[j]) >= 0.99 {
					fmt.Printf("Bar %d used. Contents: ", j)
					for i := 0; i < m; i++ {
						val := sol.Value(x.At(i, j))
						if val >= 0.99 {
							fmt.Printf("%dx size %d, ", int(val), w[i])
						}
					}
					fmt.Println()
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
