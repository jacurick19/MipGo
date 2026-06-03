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

	p := []float64{10, 13, 18, 31, 7, 15}
	w := []float64{11, 15, 20, 35, 10, 33}
	c := 47.0
	n := len(w)

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving Knapsack with %s...\n", sName)
		fmt.Printf("========================================\n")

		m := mip.New("knapsack")
		m.SetSolver(s)

		x := m.BinaryVars("x", n)

		// Set objective: Maximize sum(p[i] * x[i])
		m.Maximize(mip.Dot(p, x))

		// Add constraint: sum(w[i] * x[i]) <= c
		m.SubjectTo(mip.Dot(w, x).Leq(c))

		sol, err := m.Solve()
		if err != nil {
			log.Printf("Failed to solve knapsack with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Status: %v\n", sol.Status())
		fmt.Printf("Objective value: %g\n", sol.Objective())

		var selected []int
		for i := 0; i < n; i++ {
			if sol.Value(x[i]) >= 0.99 {
				selected = append(selected, i)
			}
		}
		fmt.Printf("Selected items: %v\n\n", selected)
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
