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

	// Number of channels per node
	r := []int{3, 5, 8, 3, 6, 5, 7, 3}

	// Distance between channels in the same node (i, i) and in adjacent nodes
	d := [][]int{
		{3, 2, 0, 0, 2, 2, 0, 0}, // 0
		{2, 3, 2, 0, 0, 2, 2, 0}, // 1
		{0, 2, 3, 0, 0, 0, 3, 0}, // 2
		{0, 0, 0, 3, 2, 0, 0, 2}, // 3
		{2, 0, 0, 2, 3, 2, 0, 0}, // 4
		{2, 2, 0, 0, 2, 3, 2, 0}, // 5
		{0, 2, 2, 0, 0, 2, 3, 0}, // 6
		{0, 0, 0, 2, 0, 0, 0, 3}, // 7
	}

	n := len(r)

	// U range sum(d[i][j]) + sum(r[i])
	sumD := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sumD += d[i][j]
		}
	}
	sumR := 0
	for i := 0; i < n; i++ {
		sumR += r[i]
	}
	uSize := sumD + sumR

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving BMCP with %s...\n", sName)
		fmt.Printf("========================================\n")

		m := mip.New("BMCP")
		m.SetSolver(s)

		// Binary variables x[i][c]
		x := m.BinaryMatrix("x", n, uSize)

		// Continuous variable z
		z := m.Continuous("z", 0, mip.Inf)

		// Minimize z
		m.Minimize(z)

		// Constraint: sum(x[i][c] for c in U) == r[i]
		for i := 0; i < n; i++ {
			var rowVars []*mip.Variable
			for c := 0; c < uSize; c++ {
				rowVars = append(rowVars, x.At(i, c))
			}
			m.SubjectTo(mip.SumVars(rowVars...).Eq(float64(r[i])))
		}

		// Constraint: x[i][c1] + x[j][c2] <= 1
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j && d[i][j] > 0 {
					for c1 := 0; c1 < uSize; c1++ {
						for c2 := c1; c2 < uSize && c2 < c1+d[i][j]; c2++ {
							m.SubjectTo(mip.Sum(x.At(i, c1), x.At(j, c2)).Leq(1.0))
						}
					}
				}
			}
		}

		// Constraint: x[i][c1] + x[i][c2] <= 1
		for i := 0; i < n; i++ {
			if d[i][i] > 1 {
				for c1 := 0; c1 < uSize; c1++ {
					for c2 := c1 + 1; c2 < uSize && c2 < c1+d[i][i]; c2++ {
						m.SubjectTo(mip.Sum(x.At(i, c1), x.At(i, c2)).Leq(1.0))
					}
				}
			}
		}

		// Constraint: z >= (c+1)*x[i][c]
		for i := 0; i < n; i++ {
			for c := 0; c < uSize; c++ {
				expr := mip.Lin(1.0, z, -float64(c+1), x.At(i, c))
				m.SubjectTo(expr.Geq(0.0))
			}
		}

		sol, err := m.Solve()
		if err != nil {
			log.Printf("Failed to solve BMCP with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value (z): %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			for i := 0; i < n; i++ {
				var channels []int
				for c := 0; c < uSize; c++ {
					if sol.Value(x.At(i, c)) >= 0.99 {
						channels = append(channels, c)
					}
				}
				fmt.Printf("Channels of node %d: %v\n", i, channels)
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
