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

	places := []string{
		"Antwerp", "Bruges", "C-Mine", "Dinant", "Ghent",
		"Grand-Place de Bruxelles", "Hasselt", "Leuven",
		"Mechelen", "Mons", "Montagne de Bueren", "Namur",
		"Remouchamps", "Waterloo",
	}

	dists := [][]float64{
		{83, 81, 113, 52, 42, 73, 44, 23, 91, 105, 90, 124, 57},
		{161, 160, 39, 89, 151, 110, 90, 99, 177, 143, 193, 100},
		{90, 125, 82, 13, 57, 71, 123, 38, 72, 59, 82},
		{123, 77, 81, 71, 91, 72, 64, 24, 62, 63},
		{51, 114, 72, 54, 69, 139, 105, 155, 62},
		{70, 25, 22, 52, 90, 56, 105, 16},
		{45, 61, 111, 36, 61, 57, 70},
		{23, 71, 67, 48, 85, 29},
		{74, 89, 69, 107, 36},
		{117, 65, 125, 43},
		{54, 22, 84},
		{60, 44},
		{97},
		{},
	}

	n := len(places)

	// Build distances matrix
	c := make([][]float64, n)
	for i := 0; i < n; i++ {
		c[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				c[i][j] = 0
			} else if j > i {
				c[i][j] = dists[i][j-i-1]
			} else {
				c[i][j] = dists[j][i-j-1]
			}
		}
	}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving TSP with %s...\n", sName)
		fmt.Printf("========================================\n")

		model := mip.New("TSP")
		model.SetSolver(s)

		// Binary variables x[i][j]
		x := model.BinaryMatrix("x", n, n)

		// Continuous variables y[i] for subtour elimination
		y := model.ContinuousVars("y", n, 0, mip.Inf)

		// Objective function: minimize sum(c[i][j] * x[i][j])
		var objTerms []mip.Term
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if c[i][j] > 0 {
					objTerms = append(objTerms, mip.Term{Coeff: c[i][j], Var: x.At(i, j)})
				}
			}
		}
		model.Minimize(mip.NewExpressionFromReader(objTerms, 0.0))

		// Constraint: leave each city only once
		for i := 0; i < n; i++ {
			var rowVars []*mip.Variable
			for j := 0; j < n; j++ {
				if j != i {
					rowVars = append(rowVars, x.At(i, j))
				}
			}
			model.SubjectTo(mip.SumVars(rowVars...).Eq(1.0))
		}

		// Constraint: enter each city only once
		for i := 0; i < n; i++ {
			var colVars []*mip.Variable
			for j := 0; j < n; j++ {
				if j != i {
					colVars = append(colVars, x.At(j, i))
				}
			}
			model.SubjectTo(mip.SumVars(colVars...).Eq(1.0))
		}

		// Subtour elimination: y[i] - y[j] - (n+1)*x[i][j] >= -n for all i,j in V - {0}, i != j
		for i := 1; i < n; i++ {
			for j := 1; j < n; j++ {
				if i != j {
					expr := mip.Lin(1.0, y[i], -1.0, y[j], -float64(n+1), x.At(i, j))
					model.SubjectTo(expr.Geq(-float64(n)))
				}
			}
		}

		sol, err := model.Solve()
		if err != nil {
			log.Printf("Failed to solve TSP with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value: %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			fmt.Printf("Route with total distance %g found: %s", sol.Objective(), places[0])
			nc := 0
			for {
				nextCity := -1
				for i := 0; i < n; i++ {
					if sol.Value(x.At(nc, i)) >= 0.99 {
						nextCity = i
						break
					}
				}
				if nextCity == -1 {
					break
				}
				fmt.Printf(" -> %s", places[nextCity])
				nc = nextCity
				if nc == 0 {
					break
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
