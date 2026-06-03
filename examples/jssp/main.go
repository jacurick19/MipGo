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

	n := 3
	m := 3

	times := [][]int{
		{2, 1, 2},
		{1, 2, 2},
		{1, 2, 1},
	}

	M := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			M += times[i][j]
		}
	}

	machines := [][]int{
		{2, 0, 1},
		{1, 2, 0},
		{2, 1, 0},
	}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving JSSP with %s...\n", sName)
		fmt.Printf("========================================\n")

		model := mip.New("JSSP")
		model.SetSolver(s)

		// Continuous variable C (makespan)
		C := model.Continuous("C", 0, mip.Inf)

		// Continuous variables x[j][i] (starting times)
		x := model.ContinuousMatrix("x", n, m, 0, mip.Inf)

		// Binary tensor y[j][k][i]
		y := model.BinaryTensor("y", n, n, m)

		model.Minimize(C)

		// Precedence constraints
		for j := 0; j < n; j++ {
			for i := 1; i < m; i++ {
				mCurr := machines[j][i]
				mPrev := machines[j][i-1]
				expr := mip.Lin(1.0, x.At(j, mCurr), -1.0, x.At(j, mPrev))
				model.SubjectTo(expr.Geq(float64(times[j][mPrev])))
			}
		}

		// Disjunctive machine constraints
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				if k != j {
					for i := 0; i < m; i++ {
						expr1 := mip.Lin(1.0, x.At(j, i), -1.0, x.At(k, i), float64(M), y.At(j, k, i))
						model.SubjectTo(expr1.Geq(float64(times[k][i])))

						expr2 := mip.Lin(-1.0, x.At(j, i), 1.0, x.At(k, i), -float64(M), y.At(j, k, i))
						model.SubjectTo(expr2.Geq(float64(times[j][i] - M)))
					}
				}
			}
		}

		// Makespan constraints
		for j := 0; j < n; j++ {
			mLast := machines[j][m-1]
			expr := mip.Lin(1.0, C, -1.0, x.At(j, mLast))
			model.SubjectTo(expr.Geq(float64(times[j][mLast])))
		}

		sol, err := model.Solve()
		if err != nil {
			log.Printf("Failed to solve JSSP with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Completion time: %g\n", sol.Value(C))

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			for j := 0; j < n; j++ {
				for i := 0; i < m; i++ {
					fmt.Printf("task %d starts on machine %d at time %g\n", j+1, i+1, sol.Value(x.At(j, i)))
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
