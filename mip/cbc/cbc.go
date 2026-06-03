package cbc

/*
#cgo LDFLAGS: -L/opt/homebrew/lib
#cgo LDFLAGS: -lCbcSolver -lcbc -lCgl -lOsiClp -lOsi -lClp -lCoinUtils -lstdc++ -lm
#include <stdlib.h>
typedef int(*cbc_progress_callback)(void *model,
                                        int phase,
                                        int step,
                                        const char *phaseName,
                                        double seconds,
                                        double lb,
                                        double ub,
                                        int nint,
                                        int *vecint,
                                        void *cbData
                                        );

    typedef void(*cbc_callback)(void *model, int msgno, int ndouble,
        const double *dvec, int nint, const int *ivec,
        int nchar, char **cvec);

    typedef void(*cbc_cut_callback)(void *osiSolver, void *osiCuts, void *appdata, int level, int npass);

    typedef int (*cbc_incumbent_callback)(void *cbcModel,
        double obj, int nz,
        char **vnames, double *x, void *appData);

    typedef void Cbc_Model;

    void *Cbc_newModel();

    void Cbc_readLp(Cbc_Model *model, const char *file);

    int Cbc_readBasis(Cbc_Model *model, const char *filename);

    int Cbc_writeBasis(Cbc_Model *model, const char *filename, char
        writeValues, int formatType);

    void Cbc_readMps(Cbc_Model *model, const char *file);

    char Cbc_supportsGzip();

    char Cbc_supportsBzip2();

    void Cbc_writeLp(Cbc_Model *model, const char *file);

    void Cbc_writeMps(Cbc_Model *model, const char *file);

    int Cbc_getNumCols(Cbc_Model *model);

    int Cbc_getNumRows(Cbc_Model *model);

    int Cbc_getNumIntegers(Cbc_Model *model);

    int Cbc_getNumElements(Cbc_Model *model);

    int Cbc_getRowNz(Cbc_Model *model, int row);

    int *Cbc_getRowIndices(Cbc_Model *model, int row);

    double *Cbc_getRowCoeffs(Cbc_Model *model, int row);

    double Cbc_getRowRHS(Cbc_Model *model, int row);

    void Cbc_setRowRHS(Cbc_Model *model, int row, double rhs);

    char Cbc_getRowSense(Cbc_Model *model, int row);

    const double *Cbc_getRowActivity(Cbc_Model *model);

    const double *Cbc_getRowSlack(Cbc_Model *model);

    int Cbc_getColNz(Cbc_Model *model, int col);

    int *Cbc_getColIndices(Cbc_Model *model, int col);

    double *Cbc_getColCoeffs(Cbc_Model *model, int col);

    void Cbc_addCol(Cbc_Model *model, const char *name,
        double lb, double ub, double obj, char isInteger,
        int nz, int *rows, double *coefs);

    void Cbc_addRow(Cbc_Model *model, const char *name, int nz,
        const int *cols, const double *coefs, char sense, double rhs);

    void Cbc_addLazyConstraint(Cbc_Model *model, int nz,
        int *cols, double *coefs, char sense, double rhs);

    void Cbc_addSOS(Cbc_Model *model, int numRows, const int *rowStarts,
        const int *colIndices, const double *weights, const int type);

    int Cbc_numberSOS(Cbc_Model *model);

    void Cbc_setObjCoeff(Cbc_Model *model, int index, double value);

    double Cbc_getObjSense(Cbc_Model *model);

    const double *Cbc_getObjCoefficients(Cbc_Model *model);

    const double *Cbc_getColSolution(Cbc_Model *model);

    const double *Cbc_getReducedCost(Cbc_Model *model);

    double *Cbc_bestSolution(Cbc_Model *model);

    int Cbc_numberSavedSolutions(Cbc_Model *model);

    const double *Cbc_savedSolution(Cbc_Model *model, int whichSol);

    double Cbc_savedSolutionObj(Cbc_Model *model, int whichSol);

    double Cbc_getObjValue(Cbc_Model *model);

    void Cbc_setObjSense(Cbc_Model *model, double sense);

    int Cbc_isProvenOptimal(Cbc_Model *model);

    int Cbc_isProvenInfeasible(Cbc_Model *model);

    int Cbc_isContinuousUnbounded(Cbc_Model *model);

    int Cbc_isAbandoned(Cbc_Model *model);

    const double *Cbc_getColLower(Cbc_Model *model);

    const double *Cbc_getColUpper(Cbc_Model *model);

    double Cbc_getColObj(Cbc_Model *model, int colIdx);

    double Cbc_getColLB(Cbc_Model *model, int colIdx);

    double Cbc_getColUB(Cbc_Model *model, int colIdx);

    void Cbc_setColLower(Cbc_Model *model, int index, double value);

    void Cbc_setColUpper(Cbc_Model *model, int index, double value);

    int Cbc_isInteger(Cbc_Model *model, int i);

    void Cbc_getColName(Cbc_Model *model,
        int iColumn, char *name, size_t maxLength);

    void Cbc_getRowName(Cbc_Model *model,
        int iRow, char *name, size_t maxLength);

    void Cbc_setContinuous(Cbc_Model *model, int iColumn);

    void Cbc_setInteger(Cbc_Model *model, int iColumn);

    enum IntParam {
        CbcMaxNumNode = 0,
        CbcMaxNumSol,
        CbcLogLevel,
        CbcMaxNumNode2
    };

    void Cbc_setIntParam(Cbc_Model *model, enum IntParam which, const int val);

    enum DblParam {
        CbcIntegerTolerance = 0,
        CbcAllowableGap,
        CbcAllowableFractionGap,
        CbcMaximumSeconds,
        CbcCutoff
    };

    void Cbc_setDblParam(Cbc_Model *model, enum DblParam which, const double val);

    void Cbc_setParameter(Cbc_Model *model, const char *name,
        const char *value);

    double Cbc_getCutoff(Cbc_Model *model);

    void Cbc_setCutoff(Cbc_Model *model, double cutoff);

    double Cbc_getAllowableGap(Cbc_Model *model);

    void Cbc_setAllowableGap(Cbc_Model *model, double allowedGap);

    double Cbc_getAllowableFractionGap(Cbc_Model *model);

    void Cbc_setAllowableFractionGap(Cbc_Model *model,
        double allowedFracionGap);

    double Cbc_getAllowablePercentageGap(Cbc_Model *model);

    void Cbc_setAllowablePercentageGap(Cbc_Model *model,
        double allowedPercentageGap);

    double Cbc_getMaximumSeconds(Cbc_Model *model);

    void Cbc_setMaximumSeconds(Cbc_Model *model, double maxSeconds);

    int Cbc_getMaximumNodes(Cbc_Model *model);

    void Cbc_setMaximumNodes(Cbc_Model *model, int maxNodes);

    int Cbc_getMaximumSolutions(Cbc_Model *model);

    void Cbc_setMaximumSolutions(Cbc_Model *model, int maxSolutions);

    int Cbc_getLogLevel(Cbc_Model *model);

    void Cbc_setLogLevel(Cbc_Model *model, int logLevel);

    double Cbc_getBestPossibleObjValue(Cbc_Model *model);

    void Cbc_setMIPStart(Cbc_Model *model, int count,
        const char **colNames, const double colValues[]);

    void Cbc_setMIPStartI(Cbc_Model *model, int count, const int colIdxs[],
        const double colValues[]);

    enum LPMethod {
        CbcLPMethodDual = 0,
        CbcLPMethodPrimal
    };

    void
    Cbc_setLPmethod(Cbc_Model *model, enum LPMethod lpm );

    void Cbc_updateConflictGraph( Cbc_Model *model );

    const void *Cbc_conflictGraph( Cbc_Model *model );

    int Cbc_solve(Cbc_Model *model);

    int Cbc_solveLinearProgram(Cbc_Model *model);

    enum CutType {
        CbcCutTypeCgl = 0
    };

    void Cbc_generateCuts( Cbc_Model *cbcModel, enum CutType ct, void *oc, int depth, int pass );

    void Cbc_strengthenPacking(Cbc_Model *model);

    void Cbc_strengthenPackingRows(Cbc_Model *model, size_t n, const size_t rows[]);

    void *Cbc_getSolverPtr(Cbc_Model *model);

    void *Cbc_deleteModel(Cbc_Model *model);

    int Osi_getNumIntegers( void *osi );

    const double *Osi_getReducedCost( void *osi );

    const double *Osi_getObjCoefficients();

    double Osi_getObjSense();

    void *Osi_newSolver();

    void Osi_deleteSolver( void *osi );

    void Osi_initialSolve(void *osi);

    void Osi_resolve(void *osi);

    void Osi_branchAndBound(void *osi);

    char Osi_isAbandoned(void *osi);

    char Osi_isProvenOptimal(void *osi);

    char Osi_isProvenPrimalInfeasible(void *osi);

    char Osi_isProvenDualInfeasible(void *osi);

    char Osi_isPrimalObjectiveLimitReached(void *osi);

    char Osi_isDualObjectiveLimitReached(void *osi);

    char Osi_isIterationLimitReached(void *osi);

    double Osi_getObjValue( void *osi );

    void Osi_setColUpper (void *osi, int elementIndex, double ub);

    void Osi_setColLower(void *osi, int elementIndex, double lb);

    int Osi_getNumCols( void *osi );

    void Osi_getColName( void *osi, int i, char *name, int maxLen );

    const double *Osi_getColLower( void *osi );

    const double *Osi_getColUpper( void *osi );

    int Osi_isInteger( void *osi, int col );

    int Osi_getNumRows( void *osi );

    int Osi_getRowNz(void *osi, int row);

    const int *Osi_getRowIndices(void *osi, int row);

    const double *Osi_getRowCoeffs(void *osi, int row);

    double Osi_getRowRHS(void *osi, int row);

    char Osi_getRowSense(void *osi, int row);

    void Osi_setObjCoef(void *osi, int index, double obj);

    void Osi_setObjSense(void *osi, double sense);

    const double *Osi_getColSolution(void *osi);

    void *OsiCuts_new();

    void OsiCuts_addRowCut( void *osiCuts, int nz, const int *idx,
        const double *coef, char sense, double rhs );

    void OsiCuts_addGlobalRowCut( void *osiCuts, int nz, const int *idx,
        const double *coef, char sense, double rhs );

    int OsiCuts_sizeRowCuts( void *osiCuts );

    int OsiCuts_nzRowCut( void *osiCuts, int iRowCut );

    const int * OsiCuts_idxRowCut( void *osiCuts, int iRowCut );

    const double *OsiCuts_coefRowCut( void *osiCuts, int iRowCut );

    double OsiCuts_rhsRowCut( void *osiCuts, int iRowCut );

    char OsiCuts_senseRowCut( void *osiCuts, int iRowCut );

    void OsiCuts_delete( void *osiCuts );

    void Osi_addCol(void *osi, const char *name, double lb, double ub,
       double obj, char isInteger, int nz, int *rows, double *coefs);

    void Osi_addRow(void *osi, const char *name, int nz,
        const int *cols, const double *coefs, char sense, double rhs);

    void Cbc_deleteRows(Cbc_Model *model, int numRows, const int rows[]);

    void Cbc_deleteCols(Cbc_Model *model, int numCols, const int cols[]);

    void Cbc_storeNameIndexes(Cbc_Model *model, char _store);

    int Cbc_getColNameIndex(Cbc_Model *model, const char *name);

    int Cbc_getRowNameIndex(Cbc_Model *model, const char *name);

    void Cbc_problemName(Cbc_Model *model, int maxNumberCharacters,
                         char *array);

    int Cbc_setProblemName(Cbc_Model *model, const char *array);

    double Cbc_getPrimalTolerance(Cbc_Model *model);

    void Cbc_setPrimalTolerance(Cbc_Model *model, double tol);

    double Cbc_getDualTolerance(Cbc_Model *model);

    void Cbc_setDualTolerance(Cbc_Model *model, double tol);

    void Cbc_addCutCallback(Cbc_Model *model, cbc_cut_callback cutcb,
        const char *name, void *appData, int howOften, char atSolution );

    void Cbc_addIncCallback(
        void *model, cbc_incumbent_callback inccb,
        void *appData );

    void Cbc_registerCallBack(Cbc_Model *model,
        cbc_callback userCallBack);

    void Cbc_addProgrCallback(void *model,
        cbc_progress_callback prgcbc, void *appData);

    void Cbc_clearCallBack(Cbc_Model *model);

    const double *Cbc_getRowPrice(Cbc_Model *model);

    const double *Osi_getRowPrice(void *osi);

    double Osi_getIntegerTolerance(void *osi);

    void Osi_checkCGraph( void *osi );

    const void * Osi_CGraph( void *osi );

    size_t CG_nodes( void *cgraph );

    char CG_conflicting( void *cgraph, int n1, int n2 );

    double CG_density( void *cgraph );

    typedef struct {
      size_t n;
      const size_t *neigh;
    } CGNeighbors;

    CGNeighbors CG_conflictingNodes(Cbc_Model *model, void *cgraph, size_t node);

    void Cbc_computeFeatures(Cbc_Model *model, double *features);

    int Cbc_nFeatures();

    const char *Cbc_featureName(int i);

    void Cbc_reset(Cbc_Model *model);
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"unsafe"

	"mipgo/v2/mip"
)

type CBCBackend struct{}

func init() {
	mip.RegisterBackend(mip.CBC, &CBCBackend{})
}

func (b *CBCBackend) Solve(model *mip.Model) (*mip.Solution, error) {
	model.ClearSolutionPool()
	cModel := C.Cbc_newModel()
	if cModel == nil {
		return nil, errors.New("failed to create CBC model instance")
	}
	defer C.Cbc_deleteModel(cModel)

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

	for _, v := range variables {
		var isInt C.char = 0
		if v.Type() == mip.Integer || v.Type() == mip.Binary {
			isInt = 1
		}

		var cName *C.char = nil
		if v.Name() != "" {
			cName = C.CString(v.Name())
		}

		C.Cbc_addCol(
			cModel,
			cName,
			C.double(v.LB()),
			C.double(v.UB()),
			C.double(objCoeffs[v.ID()]),
			isInt,
			0,
			nil,
			nil,
		)

		if cName != nil {
			C.free(unsafe.Pointer(cName))
		}
	}

	sense := 1.0
	if model.IsMaximize() {
		sense = -1.0
	}
	C.Cbc_setObjSense(cModel, C.double(sense))

	constraints := model.Constraints()
	for _, c := range constraints {
		nnz := len(c.Terms())
		if nnz == 0 {
			continue
		}

		cIndices := make([]C.int, nnz)
		cValues := make([]C.double, nnz)
		for i, t := range c.Terms() {
			cIndices[i] = C.int(t.Var.ID())
			cValues[i] = C.double(t.Coeff)
		}

		var cName *C.char = nil
		if c.Name() != "" {
			cName = C.CString(c.Name())
		} else {
			cName = C.CString("")
		}
		C.Cbc_addRow(
			cModel,
			cName,
			C.int(nnz),
			&cIndices[0],
			&cValues[0],
			C.char(c.Sense()),
			C.double(c.RHS()),
		)
		C.free(unsafe.Pointer(cName))
	}

	for _, ind := range model.Indicators() {
		addLinearizedIndicatorCBC(cModel, ind.BinaryVar(), ind.Constraint())
	}

	if len(model.SOS1()) > 0 {
		addSOSBCBC(cModel, model.SOS1(), 1)
	}
	if len(model.SOS2()) > 0 {
		addSOSBCBC(cModel, model.SOS2(), 2)
	}

	if model.TimeLimit() > 0 {
		C.Cbc_setMaximumSeconds(cModel, C.double(model.TimeLimit().Seconds()))
	}
	if model.MIPGap() >= 0 {
		C.Cbc_setAllowableFractionGap(cModel, C.double(model.MIPGap()))
	}
	if model.Threads() > 0 {
		cOptName := C.CString("threads")
		cOptVal := C.CString(fmt.Sprintf("%d", model.Threads()))
		C.Cbc_setParameter(cModel, cOptName, cOptVal)
		C.free(unsafe.Pointer(cOptName))
		C.free(unsafe.Pointer(cOptVal))
	}

	cOptNameMaxSaved := C.CString("maxSaved")
	cOptValMaxSaved := C.CString("10")
	C.Cbc_setParameter(cModel, cOptNameMaxSaved, cOptValMaxSaved)
	C.free(unsafe.Pointer(cOptNameMaxSaved))
	C.free(unsafe.Pointer(cOptValMaxSaved))

	if len(model.MIPStart()) > 0 {
		count := len(model.MIPStart())
		colIdxs := make([]C.int, count)
		colValues := make([]C.double, count)
		i := 0
		for v, val := range model.MIPStart() {
			colIdxs[i] = C.int(v.ID())
			colValues[i] = C.double(val)
			i++
		}
		C.Cbc_setMIPStartI(cModel, C.int(count), &colIdxs[0], &colValues[0])
	}

	C.Cbc_solve(cModel)

	var status mip.Status
	if C.Cbc_isProvenOptimal(cModel) != 0 {
		status = mip.Optimal
	} else if C.Cbc_isProvenInfeasible(cModel) != 0 {
		status = mip.Infeasible
	} else if C.Cbc_isContinuousUnbounded(cModel) != 0 {
		status = mip.Unbounded
	} else if C.Cbc_isAbandoned(cModel) != 0 {
		status = mip.Error
	} else {
		status = mip.Error
	}

	solValues := make(map[*mip.Variable]float64)
	solRedCosts := make(map[*mip.Variable]float64)
	solDuals := make(map[*mip.Constraint]float64)
	solSlacks := make(map[*mip.Constraint]float64)
	var stats mip.Stats

	if status == mip.Optimal || status == mip.Feasible {


		colSol := C.Cbc_getColSolution(cModel)
		colDuals := C.Cbc_getReducedCost(cModel)
		var rowDuals *C.double = nil
		var rowSlacks *C.double = nil

		if colSol != nil {
			colSolArr := unsafe.Slice(colSol, numCols)
			for _, v := range variables {
				solValues[v] = float64(colSolArr[v.ID()])
			}
		}

		if colDuals != nil {
			colDualsArr := unsafe.Slice(colDuals, numCols)
			for _, v := range variables {
				solRedCosts[v] = float64(colDualsArr[v.ID()])
			}
		}

		numRows := len(constraints)
		if rowDuals != nil && numRows > 0 {
			rowDualsArr := unsafe.Slice(rowDuals, numRows)
			for i, c := range constraints {
				solDuals[c] = float64(rowDualsArr[i])
			}
		}

		if rowSlacks != nil && numRows > 0 {
			rowSlacksArr := unsafe.Slice(rowSlacks, numRows)
			for i, c := range constraints {
				solSlacks[c] = float64(rowSlacksArr[i])
			}
		}

		nSols := int(C.Cbc_numberSavedSolutions(cModel))
		for i := 0; i < nSols; i++ {
			cVals := C.Cbc_savedSolution(cModel, C.int(i))
			obj := float64(C.Cbc_savedSolutionObj(cModel, C.int(i)))
			if model.IsMaximize() {
				obj = -obj
			}
			poolSolValues := make(map[*mip.Variable]float64)
			poolSolStatus := mip.Feasible
			if i == 0 {
				poolSolStatus = status
			}
			if cVals != nil {
				cValsArr := unsafe.Slice(cVals, numCols)
				for _, v := range variables {
					poolSolValues[v] = float64(cValsArr[v.ID()])
				}
			}
			poolSol := mip.NewSolution(
				poolSolStatus,
				obj,
				poolSolValues,
				nil,
				nil,
				nil,
				mip.Stats{},
			)
			model.AddToSolutionPool(poolSol)
		}
	}

	bestBound := float64(C.Cbc_getBestPossibleObjValue(cModel))
	gap := 0.0
	objVal := 0.0
	if status == mip.Optimal || status == mip.Feasible {
		objVal = float64(C.Cbc_getObjValue(cModel))
		if objVal != 0.0 {
			gap = math.Abs(bestBound-objVal) / math.Abs(objVal)
		}
	}

	stats = mip.Stats{
		BestBound: bestBound,
		Gap:       gap,
	}

	sol := mip.NewSolution(status, objVal, solValues, solRedCosts, solDuals, solSlacks, stats)
	return sol, nil
}

func addLinearizedIndicatorCBC(cModel unsafe.Pointer, b *mip.Variable, c *mip.Constraint) {
	if c.Sense() == 'E' {
		c1 := mip.NewConstraintFromTerms(c.Terms(), c.RHS(), 'L')
		c2 := mip.NewConstraintFromTerms(c.Terms(), c.RHS(), 'G')
		addLinearizedIndicatorCBC(cModel, b, c1)
		addLinearizedIndicatorCBC(cModel, b, c2)
		return
	}

	mVal := mip.GetBigM(c)
	nnz := len(c.Terms()) + 1
	cIndices := make([]C.int, nnz)
	cValues := make([]C.double, nnz)

	for i, t := range c.Terms() {
		cIndices[i] = C.int(t.Var.ID())
		cValues[i] = C.double(t.Coeff)
	}
	cIndices[nnz-1] = C.int(b.ID())

	var rhsVal float64
	var senseChar byte

	if c.Sense() == 'L' {
		cValues[nnz-1] = C.double(mVal)
		rhsVal = c.RHS() + mVal
		senseChar = 'L'
	} else { // 'G'
		cValues[nnz-1] = C.double(-mVal)
		rhsVal = c.RHS() - mVal
		senseChar = 'G'
	}

	cName := C.CString("")
	C.Cbc_addRow(
		cModel,
		cName,
		C.int(nnz),
		&cIndices[0],
		&cValues[0],
		C.char(senseChar),
		C.double(rhsVal),
	)
	C.free(unsafe.Pointer(cName))
}

func addSOSBCBC(cModel unsafe.Pointer, sosSets []mip.SOS, sosType int) {
	if len(sosSets) == 0 {
		return
	}
	numRows := len(sosSets)
	rowStarts := make([]C.int, numRows+1)
	var colIndices []C.int
	var weights []C.double

	var currentStart C.int = 0
	for i, s := range sosSets {
		rowStarts[i] = currentStart
		for j, v := range s.Vars() {
			colIndices = append(colIndices, C.int(v.ID()))
			weights = append(weights, C.double(j+1))
		}
		currentStart += C.int(len(s.Vars()))
	}
	rowStarts[numRows] = currentStart

	C.Cbc_addSOS(
		cModel,
		C.int(numRows),
		&rowStarts[0],
		&colIndices[0],
		&weights[0],
		C.int(sosType),
	)
}
