package scip

/*
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -lscip
#cgo darwin CFLAGS: -I/opt/homebrew/include
#cgo linux LDFLAGS: -lscip
#cgo linux CFLAGS: -I/usr/include
#include <scip/scip.h>
#include <scip/scipdefplugins.h>
#include <scip/cons_linear.h>
#include <scip/cons_sos1.h>
#include <scip/cons_sos2.h>
#include <scip/cons_indicator.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"

	"mipgo/v2/mip"
)

type SCIPBackend struct{}

func init() {
	mip.RegisterBackend(mip.SCIP, &SCIPBackend{})
}

func (b *SCIPBackend) Solve(model *mip.Model) (*mip.Solution, error) {
	var scip *C.SCIP
	if C.SCIPcreate(&scip) != C.SCIP_OKAY {
		return nil, errors.New("failed to create SCIP instance")
	}
	defer C.SCIPfree(&scip)

	if C.SCIPincludeDefaultPlugins(scip) != C.SCIP_OKAY {
		return nil, errors.New("failed to include default SCIP plugins")
	}

	cProbName := C.CString("model") // Model name is unexported now, we can use default name or export name. Wait, model has no Name() getter, we can export Name() on Model or just use "model" as default. Yes, let's just use "model".
	defer C.free(unsafe.Pointer(cProbName))
	if C.SCIPcreateProbBasic(scip, cProbName) != C.SCIP_OKAY {
		return nil, errors.New("failed to create basic SCIP problem")
	}

	// Set objective sense
	if model.IsMaximize() {
		C.SCIPsetObjsense(scip, C.SCIP_OBJSENSE_MAXIMIZE)
	} else {
		C.SCIPsetObjsense(scip, C.SCIP_OBJSENSE_MINIMIZE)
	}

	// Set verbosity to quiet by default for clean tests
	displayParam := C.CString("display/verblevel")
	C.SCIPsetIntParam(scip, displayParam, C.int(0))
	C.free(unsafe.Pointer(displayParam))

	// Set optional settings
	if model.TimeLimit() > 0 {
		timeParam := C.CString("limits/time")
		C.SCIPsetRealParam(scip, timeParam, C.double(model.TimeLimit().Seconds()))
		C.free(unsafe.Pointer(timeParam))
	}
	if model.MIPGap() >= 0 {
		gapParam := C.CString("limits/gap")
		C.SCIPsetRealParam(scip, gapParam, C.double(model.MIPGap()))
		C.free(unsafe.Pointer(gapParam))
	}
	if model.Threads() > 0 {
		threadParam := C.CString("parallel/maxnthreads")
		C.SCIPsetIntParam(scip, threadParam, C.int(model.Threads()))
		C.free(unsafe.Pointer(threadParam))
	}

	variables := model.Variables()
	numCols := len(variables)
	objCoeffs := make([]float64, numCols)
	if model.Objective() != nil {
		for _, t := range model.Objective().Terms() {
			if t.Var.ID() >= 0 && t.Var.ID() < numCols {
				objCoeffs[t.Var.ID()] += t.Coeff
			}
		}
	}

	// Create Variables
	varsMap := make(map[int]*C.SCIP_VAR)
	cVarsList := make([]*C.SCIP_VAR, numCols)

	for _, v := range variables {
		lb := v.LB()
		ub := v.UB()
		if math.IsInf(lb, -1) {
			lb = -float64(C.SCIPinfinity(scip))
		}
		if math.IsInf(ub, 1) {
			ub = float64(C.SCIPinfinity(scip))
		}

		var vartype C.SCIP_VARTYPE
		switch v.Type() {
		case mip.Binary:
			vartype = C.SCIP_VARTYPE_BINARY
		case mip.Integer:
			vartype = C.SCIP_VARTYPE_INTEGER
		case mip.Continuous:
			vartype = C.SCIP_VARTYPE_CONTINUOUS
		}

		var cVar *C.SCIP_VAR
		cName := C.CString(v.Name())
		if C.SCIPcreateVarBasic(
			scip,
			&cVar,
			cName,
			C.double(lb),
			C.double(ub),
			C.double(objCoeffs[v.ID()]),
			vartype,
		) != C.SCIP_OKAY {
			C.free(unsafe.Pointer(cName))
			return nil, fmt.Errorf("failed to create SCIP variable for %s", v.Name())
		}
		C.free(unsafe.Pointer(cName))

		if C.SCIPaddVar(scip, cVar) != C.SCIP_OKAY {
			return nil, fmt.Errorf("failed to add SCIP variable for %s", v.Name())
		}

		varsMap[v.ID()] = cVar
		cVarsList[v.ID()] = cVar
	}

	// Release variables reference on function exit
	defer func() {
		for i := range cVarsList {
			if cVarsList[i] != nil {
				C.SCIPreleaseVar(scip, &cVarsList[i])
			}
		}
	}()

	var cConsList []*C.SCIP_CONS
	defer func() {
		for i := range cConsList {
			if cConsList[i] != nil {
				C.SCIPreleaseCons(scip, &cConsList[i])
			}
		}
	}()

	constraints := model.Constraints()
	// 1. Process Standard Linear Constraints
	for _, c := range constraints {
		var cVars []*C.SCIP_VAR
		var cVals []C.double
		for _, t := range c.Terms() {
			cVars = append(cVars, varsMap[t.Var.ID()])
			cVals = append(cVals, C.double(t.Coeff))
		}

		var lhs, rhs float64
		switch c.Sense() {
		case 'E':
			lhs = c.RHS()
			rhs = c.RHS()
		case 'L':
			lhs = -float64(C.SCIPinfinity(scip))
			rhs = c.RHS()
		case 'G':
			lhs = c.RHS()
			rhs = float64(C.SCIPinfinity(scip))
		}

		var varsPtr **C.SCIP_VAR
		var valsPtr *C.double
		if len(cVars) > 0 {
			varsPtr = &cVars[0]
			valsPtr = &cVals[0]
		}

		var cCons *C.SCIP_CONS
		cName := C.CString(c.Name())
		if C.SCIPcreateConsBasicLinear(
			scip,
			&cCons,
			cName,
			C.int(len(cVars)),
			varsPtr,
			valsPtr,
			C.double(lhs),
			C.double(rhs),
		) != C.SCIP_OKAY {
			C.free(unsafe.Pointer(cName))
			return nil, fmt.Errorf("failed to create linear constraint %s", c.Name())
		}
		C.free(unsafe.Pointer(cName))

		if C.SCIPaddCons(scip, cCons) != C.SCIP_OKAY {
			return nil, fmt.Errorf("failed to add linear constraint %s", c.Name())
		}
		cConsList = append(cConsList, cCons)
	}

	// 2. Process Native SOS1 Constraints
	for i, s := range model.SOS1() {
		var cVars []*C.SCIP_VAR
		var weights []C.double
		for j, v := range s.Vars() {
			cVars = append(cVars, varsMap[v.ID()])
			weights = append(weights, C.double(j+1))
		}

		var varsPtr **C.SCIP_VAR
		var weightsPtr *C.double
		if len(cVars) > 0 {
			varsPtr = &cVars[0]
			weightsPtr = &weights[0]
		}

		var cCons *C.SCIP_CONS
		cName := C.CString(fmt.Sprintf("sos1_%d", i))
		if C.SCIPcreateConsBasicSOS1(
			scip,
			&cCons,
			cName,
			C.int(len(cVars)),
			varsPtr,
			weightsPtr,
		) != C.SCIP_OKAY {
			C.free(unsafe.Pointer(cName))
			return nil, fmt.Errorf("failed to create SOS1 constraint %d", i)
		}
		C.free(unsafe.Pointer(cName))

		if C.SCIPaddCons(scip, cCons) != C.SCIP_OKAY {
			return nil, fmt.Errorf("failed to add SOS1 constraint %d", i)
		}
		cConsList = append(cConsList, cCons)
	}

	// 3. Process Native SOS2 Constraints
	for i, s := range model.SOS2() {
		var cVars []*C.SCIP_VAR
		var weights []C.double
		for j, v := range s.Vars() {
			cVars = append(cVars, varsMap[v.ID()])
			weights = append(weights, C.double(j+1))
		}

		var varsPtr **C.SCIP_VAR
		var weightsPtr *C.double
		if len(cVars) > 0 {
			varsPtr = &cVars[0]
			weightsPtr = &weights[0]
		}

		var cCons *C.SCIP_CONS
		cName := C.CString(fmt.Sprintf("sos2_%d", i))
		if C.SCIPcreateConsBasicSOS2(
			scip,
			&cCons,
			cName,
			C.int(len(cVars)),
			varsPtr,
			weightsPtr,
		) != C.SCIP_OKAY {
			C.free(unsafe.Pointer(cName))
			return nil, fmt.Errorf("failed to create SOS2 constraint %d", i)
		}
		C.free(unsafe.Pointer(cName))

		if C.SCIPaddCons(scip, cCons) != C.SCIP_OKAY {
			return nil, fmt.Errorf("failed to add SOS2 constraint %d", i)
		}
		cConsList = append(cConsList, cCons)
	}

	// 4. Process Native Indicator Constraints
	for i, ind := range model.Indicators() {
		var cVars []*C.SCIP_VAR
		var cVals []C.double
		for _, t := range ind.Constraint().Terms() {
			cVars = append(cVars, varsMap[t.Var.ID()])
			cVals = append(cVals, C.double(t.Coeff))
		}

		var varsPtr **C.SCIP_VAR
		var valsPtr *C.double
		if len(cVars) > 0 {
			varsPtr = &cVars[0]
			valsPtr = &cVals[0]
		}

		binVar := varsMap[ind.BinaryVar().ID()]

		if ind.Constraint().Sense() == 'L' {
			var indCons *C.SCIP_CONS
			indName := C.CString(fmt.Sprintf("ind_%d", i))
			if C.SCIPcreateConsBasicIndicator(
				scip,
				&indCons,
				indName,
				binVar,
				C.int(len(cVars)),
				varsPtr,
				valsPtr,
				C.double(ind.Constraint().RHS()),
			) != C.SCIP_OKAY {
				C.free(unsafe.Pointer(indName))
				return nil, fmt.Errorf("failed to create indicator constraint %d", i)
			}
			C.free(unsafe.Pointer(indName))
			if C.SCIPaddCons(scip, indCons) != C.SCIP_OKAY {
				return nil, fmt.Errorf("failed to add indicator constraint %d", i)
			}
			cConsList = append(cConsList, indCons)
		} else if ind.Constraint().Sense() == 'G' {
			// Convert >= to <= by negating coefficients and rhs
			negVals := make([]C.double, len(cVals))
			for k, val := range cVals {
				negVals[k] = -val
			}
			var negValsPtr *C.double
			if len(negVals) > 0 {
				negValsPtr = &negVals[0]
			}
			var indCons *C.SCIP_CONS
			indName := C.CString(fmt.Sprintf("ind_%d", i))
			if C.SCIPcreateConsBasicIndicator(
				scip,
				&indCons,
				indName,
				binVar,
				C.int(len(cVars)),
				varsPtr,
				negValsPtr,
				C.double(-ind.Constraint().RHS()),
			) != C.SCIP_OKAY {
				C.free(unsafe.Pointer(indName))
				return nil, fmt.Errorf("failed to create indicator constraint %d", i)
			}
			C.free(unsafe.Pointer(indName))
			if C.SCIPaddCons(scip, indCons) != C.SCIP_OKAY {
				return nil, fmt.Errorf("failed to add indicator constraint %d", i)
			}
			cConsList = append(cConsList, indCons)
		} else if ind.Constraint().Sense() == 'E' {
			// Split = into <= and >= (converted to <=)
			var indCons1 *C.SCIP_CONS
			indName1 := C.CString(fmt.Sprintf("ind_%d_eq1", i))
			if C.SCIPcreateConsBasicIndicator(
				scip,
				&indCons1,
				indName1,
				binVar,
				C.int(len(cVars)),
				varsPtr,
				valsPtr,
				C.double(ind.Constraint().RHS()),
			) != C.SCIP_OKAY {
				C.free(unsafe.Pointer(indName1))
				return nil, fmt.Errorf("failed to create indicator constraint %d (eq1)", i)
			}
			C.free(unsafe.Pointer(indName1))
			if C.SCIPaddCons(scip, indCons1) != C.SCIP_OKAY {
				return nil, fmt.Errorf("failed to add indicator constraint %d (eq1)", i)
			}
			cConsList = append(cConsList, indCons1)

			negVals := make([]C.double, len(cVals))
			for k, val := range cVals {
				negVals[k] = -val
			}
			var negValsPtr *C.double
			if len(negVals) > 0 {
				negValsPtr = &negVals[0]
			}
			var indCons2 *C.SCIP_CONS
			indName2 := C.CString(fmt.Sprintf("ind_%d_eq2", i))
			if C.SCIPcreateConsBasicIndicator(
				scip,
				&indCons2,
				indName2,
				binVar,
				C.int(len(cVars)),
				varsPtr,
				negValsPtr,
				C.double(-ind.Constraint().RHS()),
			) != C.SCIP_OKAY {
				C.free(unsafe.Pointer(indName2))
				return nil, fmt.Errorf("failed to create indicator constraint %d (eq2)", i)
			}
			C.free(unsafe.Pointer(indName2))
			if C.SCIPaddCons(scip, indCons2) != C.SCIP_OKAY {
				return nil, fmt.Errorf("failed to add indicator constraint %d (eq2)", i)
			}
			cConsList = append(cConsList, indCons2)
		}
	}

	// Set MIPStart (initial solution) if provided
	if len(model.MIPStart()) > 0 {
		var sol *C.SCIP_SOL
		if C.SCIPcreateSol(scip, &sol, nil) == C.SCIP_OKAY {
			for v, val := range model.MIPStart() {
				C.SCIPsetSolVal(scip, sol, varsMap[v.ID()], C.double(val))
			}
			var stored C.SCIP_Bool
			C.SCIPaddSolFree(scip, &sol, &stored)
		}
	}

	// Solve the model
	if C.SCIPsolve(scip) != C.SCIP_OKAY {
		return nil, errors.New("failed to solve problem with SCIP")
	}

	scipStatus := C.SCIPgetStatus(scip)
	var status mip.Status
	switch scipStatus {
	case C.SCIP_STATUS_OPTIMAL:
		status = mip.Optimal
	case C.SCIP_STATUS_BESTSOLLIMIT, C.SCIP_STATUS_SOLLIMIT:
		status = mip.Feasible
	case C.SCIP_STATUS_INFEASIBLE, C.SCIP_STATUS_INFORUNBD:
		status = mip.Infeasible
	case C.SCIP_STATUS_UNBOUNDED:
		status = mip.Unbounded
	case C.SCIP_STATUS_TIMELIMIT:
		status = mip.TimeLimit
	case C.SCIP_STATUS_USERINTERRUPT, C.SCIP_STATUS_TERMINATE:
		status = mip.Interrupted
	default:
		status = mip.Error
	}

	solValues := make(map[*mip.Variable]float64)
	solRedCosts := make(map[*mip.Variable]float64)
	solDuals := make(map[*mip.Constraint]float64)
	solSlacks := make(map[*mip.Constraint]float64)
	var stats mip.Stats

	model.ClearSolutionPool()
	if status == mip.Optimal || status == mip.Feasible {
		nSols := int(C.SCIPgetNSols(scip))
		sols := C.SCIPgetSols(scip)
		if nSols > 0 && sols != nil {
			solsSlice := unsafe.Slice(sols, nSols)
			for i := 0; i < nSols; i++ {
				cSol := solsSlice[i]
				if cSol == nil {
					continue
				}
				poolSolValues := make(map[*mip.Variable]float64)
				poolSolStatus := mip.Feasible
				if i == 0 {
					poolSolStatus = status
				}
				for _, v := range variables {
					poolSolValues[v] = float64(C.SCIPgetSolVal(scip, cSol, varsMap[v.ID()]))
				}
				poolSol := mip.NewSolution(
					poolSolStatus,
					float64(C.SCIPgetSolOrigObj(scip, cSol)),
					poolSolValues,
					nil,
					nil,
					nil,
					mip.Stats{},
				)
				model.AddToSolutionPool(poolSol)
			}
		}
		if len(model.SolutionPool()) > 0 {
			solValues = model.SolutionPool()[0].VariablesValues()
		}
	}

	stats = mip.Stats{
		Nodes:             int64(C.SCIPgetNNodes(scip)),
		SimplexIterations: int64(C.SCIPgetNLPIterations(scip)),
		Runtime:           time.Duration(float64(C.SCIPgetSolvingTime(scip)) * float64(time.Second)),
		BestBound:         float64(C.SCIPgetDualbound(scip)),
		Gap:               float64(C.SCIPgetGap(scip)),
	}

	objVal := 0.0
	if len(model.SolutionPool()) > 0 {
		objVal = model.SolutionPool()[0].Objective()
	}

	sol := mip.NewSolution(status, objVal, solValues, solRedCosts, solDuals, solSlacks, stats)
	return sol, nil
}
