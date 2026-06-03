package main

import (
	"fmt"
	"log"
	"math"

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

	// Possible plants
	F := []int{1, 2, 3, 4, 5, 6}

	// Possible plant installation positions
	pf := map[int][2]float64{
		1: {1, 38}, 2: {31, 40}, 3: {23, 59},
		4: {76, 51}, 5: {93, 51}, 6: {63, 74},
	}

	// Maximum plant capacity
	c := map[int]float64{
		1: 1955, 2: 1932, 3: 1987,
		4: 1823, 5: 1718, 6: 1742,
	}

	// Clients
	C := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Position of clients
	pc := map[int][2]float64{
		1: {94, 10}, 2: {57, 26}, 3: {74, 44}, 4: {27, 51}, 5: {78, 30},
		6: {23, 30}, 7: {20, 72}, 8: {3, 27}, 9: {5, 39}, 10: {51, 1},
	}

	// Demands
	d := map[int]float64{
		1: 302, 2: 273, 3: 275, 4: 266, 5: 287,
		6: 296, 7: 297, 8: 310, 9: 302, 10: 309,
	}

	// Distance matrix rounded to 1 decimal place
	dist := make(map[[2]int]float64)
	for _, fVal := range F {
		for _, cVal := range C {
			dx := pf[fVal][0] - pc[cVal][0]
			dy := pf[fVal][1] - pc[cVal][1]
			val := math.Sqrt(dx*dx + dy*dy)
			dist[[2]int{fVal, cVal}] = math.Round(val*10) / 10
		}
	}

	for _, s := range solvers {
		sName := solverName(s)
		fmt.Printf("========================================\n")
		fmt.Printf("Solving Plant Location with %s...\n", sName)
		fmt.Printf("========================================\n")

		if s == mip.HiGHS {
			fmt.Println("Skipping HiGHS for Plant Location because HiGHS does not support SOS2 constraints.")
			fmt.Println()
			continue
		}

		model := mip.New("PlantLocation")
		model.SetSolver(s)

		// Continuous variables z[i] for plant capacity
		z := make(map[int]*mip.Variable)
		for _, i := range F {
			z[i] = model.Continuous(fmt.Sprintf("z_%d", i), 0, c[i])
		}

		// Type 1 SOS: only one plant per region
		for r := 0; r <= 1; r++ {
			var Fr []*mip.Variable
			for _, i := range F {
				xCoord := pf[i][0]
				if float64(r)*50.0 <= xCoord && xCoord <= 50.0+float64(r)*50.0 {
					Fr = append(Fr, z[i])
				}
			}
			model.AddSOS1(Fr...)
		}

		// Continuous variables x[i][j]: amount plant i supplies to client j
		x := make(map[[2]int]*mip.Variable)
		for _, i := range F {
			for _, j := range C {
				x[[2]int{i, j}] = model.Continuous(fmt.Sprintf("x_%d_%d", i, j), 0, mip.Inf)
			}
		}

		// Constraint: Satisfy client demand
		for _, j := range C {
			var list []*mip.Variable
			for _, i := range F {
				list = append(list, x[[2]int{i, j}])
			}
			model.SubjectTo(mip.SumVars(list...).Eq(d[j]))
		}

		// SOS type 2 to model non-linear installation costs
		y := make(map[int]*mip.Variable)
		for _, i := range F {
			y[i] = model.Continuous(fmt.Sprintf("y_%d", i), 0, mip.Inf)
		}

		for _, f := range F {
			D := 6
			v := make([]float64, D)
			vn := make([]float64, D)
			for k := 0; k < D; k++ {
				v[k] = c[f] * float64(k) / float64(D-1)
				if k == 0 {
					vn[k] = 0.0
				} else {
					vn[k] = 1520.0 * math.Log(v[k])
				}
			}

			// Discretization weight variables w
			w := model.ContinuousVars(fmt.Sprintf("w_%d", f), D, 0, mip.Inf)

			// Convexification: sum(w_k) == 1
			model.SubjectTo(mip.SumVars(w...).Eq(1.0))

			// Link to z var: z[f] = sum(v[k] * w[k])
			var zTerms []mip.Term
			zTerms = append(zTerms, mip.Term{Coeff: 1.0, Var: z[f]})
			for k := 0; k < D; k++ {
				zTerms = append(zTerms, mip.Term{Coeff: -v[k], Var: w[k]})
			}
			model.SubjectTo(mip.NewExpressionFromReader(zTerms, 0.0).Eq(0.0))

			// Link to y var: y[f] = sum(vn[k] * w[k])
			var yTerms []mip.Term
			yTerms = append(yTerms, mip.Term{Coeff: 1.0, Var: y[f]})
			for k := 0; k < D; k++ {
				yTerms = append(yTerms, mip.Term{Coeff: -vn[k], Var: w[k]})
			}
			model.SubjectTo(mip.NewExpressionFromReader(yTerms, 0.0).Eq(0.0))

			// Add SOS2
			model.AddSOS2(w...)
		}

		// Constraint: Capacity of each plant
		for _, i := range F {
			var terms []mip.Term
			terms = append(terms, mip.Term{Coeff: 1.0, Var: z[i]})
			for _, j := range C {
				terms = append(terms, mip.Term{Coeff: -1.0, Var: x[[2]int{i, j}]})
			}
			model.SubjectTo(mip.NewExpressionFromReader(terms, 0.0).Geq(0.0))
		}

		// Objective: Minimize shipping cost + installation cost
		var objTerms []mip.Term
		for _, i := range F {
			for _, j := range C {
				key := [2]int{i, j}
				objTerms = append(objTerms, mip.Term{Coeff: dist[key], Var: x[key]})
			}
			objTerms = append(objTerms, mip.Term{Coeff: 1.0, Var: y[i]})
		}
		model.Minimize(mip.NewExpressionFromReader(objTerms, 0.0))

		sol, err := model.Solve()
		if err != nil {
			log.Printf("Failed to solve Plant Location with %s: %v\n", sName, err)
			continue
		}

		fmt.Printf("Objective value: %g\n", sol.Objective())

		if sol.Status() == mip.Optimal || sol.Status() == mip.Feasible {
			for _, f := range F {
				fmt.Printf("Plant %d: capacity = %g, cost = %g\n", f, sol.Value(z[f]), sol.Value(y[f]))
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
