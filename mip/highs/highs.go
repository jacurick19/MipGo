package highs

/*
#cgo LDFLAGS: -L/opt/homebrew/lib
#cgo LDFLAGS: -lhighs -lstdc++ -lm
#include <stdlib.h>

        typedef struct {
          void* cbdata;
          int log_type;
          double running_time;
          int simplex_iteration_count;
          int ipm_iteration_count;
          int pdlp_iteration_count;
          double objective_function_value;
          int64_t mip_node_count;
          int64_t mip_total_lp_iterations;
          double mip_primal_bound;
          double mip_dual_bound;
          double mip_gap;
          double* mip_solution;
          int mip_solution_size;
        } HighsCallbackDataOut;

        typedef struct {
          int user_interrupt;
        } HighsCallbackDataIn;

        typedef void (*HighsCCallbackType)(int, const char*,
                                           const HighsCallbackDataOut*,
                                           HighsCallbackDataIn*, void*);

        #define kHighsMaximumStringLength 512

        #define kHighsStatusError -1
        #define kHighsStatusOk 0
        #define kHighsStatusWarning 1

        #define kHighsVarTypeContinuous 0
        #define kHighsVarTypeInteger 1
        #define kHighsVarTypeSemiContinuous 2
        #define kHighsVarTypeSemiInteger 3
        #define kHighsVarTypeImplicitInteger 4

        #define kHighsOptionTypeBool 0
        #define kHighsOptionTypeInt 1
        #define kHighsOptionTypeDouble 2
        #define kHighsOptionTypeString 3

        #define kHighsInfoTypeInt64 -1
        #define kHighsInfoTypeInt 1
        #define kHighsInfoTypeDouble 2

        #define kHighsObjSenseMinimize 1
        #define kHighsObjSenseMaximize -1

        #define kHighsMatrixFormatColwise 1
        #define kHighsMatrixFormatRowwise 2

        #define kHighsHessianFormatTriangular 1
        #define kHighsHessianFormatSquare 2

        #define kHighsSolutionStatusNone 0
        #define kHighsSolutionStatusInfeasible 1
        #define kHighsSolutionStatusFeasible 2

        #define kHighsBasisValidityInvalid 0
        #define kHighsBasisValidityValid 1

        #define kHighsPresolveStatusNotPresolved -1
        #define kHighsPresolveStatusNotReduced 0
        #define kHighsPresolveStatusInfeasible 1
        #define kHighsPresolveStatusUnboundedOrInfeasible 2
        #define kHighsPresolveStatusReduced 3
        #define kHighsPresolveStatusReducedToEmpty 4
        #define kHighsPresolveStatusTimeout 5
        #define kHighsPresolveStatusNullError 6
        #define kHighsPresolveStatusOptionsError 7

        #define kHighsModelStatusNotset 0
        #define kHighsModelStatusLoadError 1
        #define kHighsModelStatusModelError 2
        #define kHighsModelStatusPresolveError 3
        #define kHighsModelStatusSolveError 4
        #define kHighsModelStatusPostsolveError 5
        #define kHighsModelStatusModelEmpty 6
        #define kHighsModelStatusOptimal 7
        #define kHighsModelStatusInfeasible 8
        #define kHighsModelStatusUnboundedOrInfeasible 9
        #define kHighsModelStatusUnbounded 10
        #define kHighsModelStatusObjectiveBound 11
        #define kHighsModelStatusObjectiveTarget 12
        #define kHighsModelStatusTimeLimit 13
        #define kHighsModelStatusIterationLimit 14
        #define kHighsModelStatusUnknown 15
        #define kHighsModelStatusSolutionLimit 16
        #define kHighsModelStatusInterrupt 17

        #define kHighsBasisStatusLower 0
        #define kHighsBasisStatusBasic 1
        #define kHighsBasisStatusUpper 2
        #define kHighsBasisStatusZero 3
        #define kHighsBasisStatusNonbasic 4

        #define kHighsCallbackLogging 0
        #define kHighsCallbackSimplexInterrupt 1
        #define kHighsCallbackIpmInterrupt 2
        #define kHighsCallbackMipSolution 3
        #define kHighsCallbackMipImprovingSolution 4
        #define kHighsCallbackMipLogging 5
        #define kHighsCallbackMipInterrupt 6

        int Highs_lpCall(const int num_col, const int num_row,
                              const int num_nz, const int a_format,
                              const int sense, const double offset,
                              const double* col_cost, const double* col_lower,
                              const double* col_upper, const double* row_lower,
                              const double* row_upper, const int* a_start,
                              const int* a_index, const double* a_value,
                              double* col_value, double* col_dual, double* row_value,
                              double* row_dual, int* col_basis_status,
                              int* row_basis_status, int* model_status);

        int Highs_mipCall(const int num_col, const int num_row,
                               const int num_nz, const int a_format,
                               const int sense, const double offset,
                               const double* col_cost, const double* col_lower,
                               const double* col_upper, const double* row_lower,
                               const double* row_upper, const int* a_start,
                               const int* a_index, const double* a_value,
                               const int* integrality, double* col_value,
                               double* row_value, int* model_status);

        int Highs_qpCall(
            const int num_col, const int num_row, const int num_nz,
            const int q_num_nz, const int a_format, const int q_format,
            const int sense, const double offset, const double* col_cost,
            const double* col_lower, const double* col_upper, const double* row_lower,
            const double* row_upper, const int* a_start, const int* a_index,
            const double* a_value, const int* q_start, const int* q_index,
            const double* q_value, double* col_value, double* col_dual,
            double* row_value, double* row_dual, int* col_basis_status,
            int* row_basis_status, int* model_status);

        void* Highs_create(void);

        void Highs_destroy(void* highs);

        const char* Highs_version(void);

        int Highs_versionMajor(void);

        int Highs_versionMinor(void);

        int Highs_versionPatch(void);

        const char* Highs_githash(void);

        const char* Highs_compilationDate(void);

        int Highs_readModel(void* highs, const char* filename);

        int Highs_writeModel(void* highs, const char* filename);

        int Highs_clear(void* highs);

        int Highs_clearModel(void* highs);

        int Highs_clearSolver(void* highs);

        int Highs_run(void* highs);

        int Highs_writeSolution(const void* highs, const char* filename);

        int Highs_writeSolutionPretty(const void* highs, const char* filename);

        int Highs_passLp(void* highs, const int num_col,
                              const int num_row, const int num_nz,
                              const int a_format, const int sense,
                              const double offset, const double* col_cost,
                              const double* col_lower, const double* col_upper,
                              const double* row_lower, const double* row_upper,
                              const int* a_start, const int* a_index,
                              const double* a_value);

        int Highs_passMip(void* highs, const int num_col,
                               const int num_row, const int num_nz,
                               const int a_format, const int sense,
                               const double offset, const double* col_cost,
                               const double* col_lower, const double* col_upper,
                               const double* row_lower, const double* row_upper,
                               const int* a_start, const int* a_index,
                               const double* a_value, const int* integrality);

        int Highs_passModel(void* highs, const int num_col,
                                 const int num_row, const int num_nz,
                                 const int q_num_nz, const int a_format,
                                 const int q_format, const int sense,
                                 const double offset, const double* col_cost,
                                 const double* col_lower, const double* col_upper,
                                 const double* row_lower, const double* row_upper,
                                 const int* a_start, const int* a_index,
                                 const double* a_value, const int* q_start,
                                 const int* q_index, const double* q_value,
                                 const int* integrality);

        int Highs_passHessian(void* highs, const int dim,
                                   const int num_nz, const int format,
                                   const int* start, const int* index,
                                   const double* value);

        int Highs_passRowName(const void* highs, const int row,
                                   const char* name);

        int Highs_passColName(const void* highs, const int col,
                                   const char* name);

        int Highs_readOptions(const void* highs, const char* filename);

        int Highs_setBoolOptionValue(void* highs, const char* option,
                                          const int value);

        int Highs_setIntOptionValue(void* highs, const char* option,
                                          const int value);

        int Highs_setDoubleOptionValue(void* highs, const char* option,
                                            const double value);

        int Highs_setStringOptionValue(void* highs, const char* option,
                                            const char* value);

        int Highs_getBoolOptionValue(const void* highs, const char* option,
                                          int* value);

        int Highs_getIntOptionValue(const void* highs, const char* option,
                                         int* value);

        int Highs_getDoubleOptionValue(const void* highs, const char* option,
                                            double* value);

        int Highs_getStringOptionValue(const void* highs, const char* option,
                                            char* value);

        int Highs_getOptionType(const void* highs, const char* option,
                                     int* type);

        int Highs_resetOptions(void* highs);

        int Highs_writeOptions(const void* highs, const char* filename);

        int Highs_writeOptionsDeviations(const void* highs, const char* filename);

        int Highs_getNumOptions(const void* highs);

        int Highs_getOptionName(const void* highs, const int index,
                                     char** name);

        int Highs_getBoolOptionValues(const void* highs, const char* option,
                                           int* current_value,
                                           int* default_value);
        int Highs_getIntOptionValues(const void* highs, const char* option,
                                          int* current_value, int* min_value,
                                          int* max_value, int* default_value);

        int Highs_getDoubleOptionValues(const void* highs, const char* option,
                                             double* current_value, double* min_value,
                                             double* max_value, double* default_value);

        int Highs_getStringOptionValues(const void* highs, const char* option,
                                             char* current_value, char* default_value);

        int Highs_getIntInfoValue(const void* highs, const char* info,
                                       int* value);

        int Highs_getDoubleInfoValue(const void* highs, const char* info,
                                           double* value);

        int Highs_getInt64InfoValue(const void* highs, const char* info,
                                         int64_t* value);

        int Highs_getInfoType(const void* highs, const char* info, int* type);

        int Highs_getSolution(const void* highs, double* col_value,
                                   double* col_dual, double* row_value,
                                   double* row_dual);

        int Highs_getBasis(const void* highs, int* col_status,
                                int* row_status);

        int Highs_getModelStatus(const void* highs);

        int Highs_getDualRay(const void* highs, int* has_dual_ray,
                                  double* dual_ray_value);

        int Highs_getPrimalRay(const void* highs, int* has_primal_ray,
                                    double* primal_ray_value);

        double Highs_getObjectiveValue(const void* highs);

        int Highs_getBasicVariables(const void* highs, int* basic_variables);

        int Highs_getBasisInverseRow(const void* highs, const int row,
                                          double* row_vector, int* row_num_nz,
                                          int* row_index);

        int Highs_getBasisInverseCol(const void* highs, const int col,
                                          double* col_vector, int* col_num_nz,
                                          int* col_index);

        int Highs_getBasisSolve(const void* highs, const double* rhs,
                                     double* solution_vector, int* solution_num_nz,
                                     int* solution_index);

        int Highs_getBasisTransposeSolve(const void* highs, const double* rhs,
                                              double* solution_vector,
                                              int* solution_nz,
                                              int* solution_index);

        int Highs_getReducedRow(const void* highs, const int row,
                                     double* row_vector, int* row_num_nz,
                                     int* row_index);

        int Highs_getReducedColumn(const void* highs, const int col,
                                        double* col_vector, int* col_num_nz,
                                        int* col_index);

        int Highs_setBasis(void* highs, const int* col_status,
                                const int* row_status);

        int Highs_setLogicalBasis(void* highs);

        int Highs_setSolution(void* highs, const double* col_value,
                                   const double* row_value, const double* col_dual,
                                   const double* row_dual);

        int Highs_setCallback(void* highs, HighsCCallbackType user_callback,
                                   void* user_callback_data);

        int Highs_startCallback(void* highs, const int callback_type);

        int Highs_stopCallback(void* highs, const int callback_type);

        double Highs_getRunTime(const void* highs);

        int Highs_zeroAllClocks(const void* highs);

        int Highs_addCol(void* highs, const double cost, const double lower,
                              const double upper, const int num_new_nz,
                              const int* index, const double* value);

        int Highs_addCols(void* highs, const int num_new_col,
                               const double* costs, const double* lower,
                               const double* upper, const int num_new_nz,
                               const int* starts, const int* index,
                               const double* value);

        int Highs_addVar(void* highs, const double lower, const double upper);

        int Highs_addVars(void* highs, const int num_new_var,
                               const double* lower, const double* upper);

        int Highs_addRow(void* highs, const double lower, const double upper,
                              const int num_new_nz, const int* index,
                              const double* value);

        int Highs_addRows(void* highs, const int num_new_row,
                               const double* lower, const double* upper,
                               const int num_new_nz, const int* starts,
                               const int* index, const double* value);

        int Highs_changeObjectiveSense(void* highs, const int sense);

        int Highs_changeObjectiveOffset(void* highs, const double offset);

        int Highs_changeColIntegrality(void* highs, const int col,
                                            const int integrality);

        int Highs_changeColsIntegralityByRange(void* highs,
                                                    const int from_col,
                                                    const int to_col,
                                                    const int* integrality);

        int Highs_changeColsIntegralityBySet(void* highs,
                                                  const int num_set_entries,
                                                  const int* set,
                                                  const int* integrality);

        int Highs_changeColsIntegralityByMask(void* highs, const int* mask,
                                                   const int* integrality);

        int Highs_changeColCost(void* highs, const int col,
                                     const double cost);

        int Highs_changeColsCostByRange(void* highs, const int from_col,
                                             const int to_col, const double* cost);

        int Highs_changeColsCostBySet(void* highs, const int num_set_entries,
                                           const int* set, const double* cost);

        int Highs_changeColsCostByMask(void* highs, const int* mask,
                                            const double* cost);

        int Highs_changeColBounds(void* highs, const int col,
                                       const double lower, const double upper);

        int Highs_changeColsBoundsByRange(void* highs, const int from_col,
                                               const int to_col,
                                               const double* lower,
                                               const double* upper);

        int Highs_changeColsBoundsBySet(void* highs,
                                             const int num_set_entries,
                                             const int* set, const double* lower,
                                             const double* upper);

        int Highs_changeColsBoundsByMask(void* highs, const int* mask,
                                              const double* lower, const double* upper);

        int Highs_changeRowBounds(void* highs, const int row,
                                       const double lower, const double upper);

        int Highs_changeRowsBoundsBySet(void* highs,
                                             const int num_set_entries,
                                             const int* set, const double* lower,
                                             const double* upper);

        int Highs_changeRowsBoundsByMask(void* highs, const int* mask,
                                              const double* lower, const double* upper);

        int Highs_changeCoeff(void* highs, const int row, const int col,
                                   const double value);

        int Highs_getObjectiveSense(const void* highs, int* sense);

        int Highs_getObjectiveOffset(const void* highs, double* offset);

        int Highs_getColsByRange(const void* highs, const int from_col,
                                      const int to_col, int* num_col,
                                      double* costs, double* lower, double* upper,
                                      int* num_nz, int* matrix_start,
                                      int* matrix_index, double* matrix_value);

        int Highs_getColsBySet(const void* highs, const int num_set_entries,
                                    const int* set, int* num_col,
                                    double* costs, double* lower, double* upper,
                                    int* num_nz, int* matrix_start,
                                    int* matrix_index, double* matrix_value);

        int Highs_getColsByMask(const void* highs, const int* mask,
                                     int* num_col, double* costs, double* lower,
                                     double* upper, int* num_nz,
                                     int* matrix_start, int* matrix_index,
                                     double* matrix_value);

        int Highs_getRowsByRange(const void* highs, const int from_row,
                                      const int to_row, int* num_row,
                                      double* lower, double* upper, int* num_nz,
                                      int* matrix_start, int* matrix_index,
                                      double* matrix_value);

        int Highs_getRowsBySet(const void* highs, const int num_set_entries,
                                    const int* set, int* num_row,
                                    double* lower, double* upper, int* num_nz,
                                    int* matrix_start, int* matrix_index,
                                    double* matrix_value);

        int Highs_getRowsByMask(const void* highs, const int* mask,
                                     int* num_row, double* lower, double* upper,
                                     int* num_nz, int* matrix_start,
                                     int* matrix_index, double* matrix_value);
        int Highs_getRowName(const void* highs, const int row, char* name);

        int Highs_getRowByName(const void* highs, const char* name, int* row);

        int Highs_getColName(const void* highs, const int col, char* name);

        int Highs_getColByName(const void* highs, const char* name, int* col);

        int Highs_getColIntegrality(const void* highs, const int col,
                                         int* integrality);

        int Highs_deleteColsByRange(void* highs, const int from_col,
                                         const int to_col);

        int Highs_deleteColsBySet(void* highs, const int num_set_entries,
                                       const int* set);

        int Highs_deleteColsByMask(void* highs, int* mask);

        int Highs_deleteRowsByRange(void* highs, const int from_row,
                                         const int to_row);

        int Highs_deleteRowsBySet(void* highs, const int num_set_entries,
                                       const int* set);

        int Highs_deleteRowsByMask(void* highs, int* mask);

        int Highs_scaleCol(void* highs, const int col, const double scaleval);

        int Highs_scaleRow(void* highs, const int row, const double scaleval);

        double Highs_getInfinity(const void* highs);

        int Highs_getSizeofint(const void* highs);

        int Highs_getNumCol(const void* highs);

        int Highs_getNumRow(const void* highs);

        int Highs_getNumNz(const void* highs);

        int Highs_getHessianNumNz(const void* highs);

        int Highs_getModel(const void* highs, const int a_format,
                                const int q_format, int* num_col,
                                int* num_row, int* num_nz,
                                int* hessian_num_nz, int* sense,
                                double* offset, double* col_cost, double* col_lower,
                                double* col_upper, double* row_lower, double* row_upper,
                                int* a_start, int* a_index, double* a_value,
                                int* q_start, int* q_index, double* q_value,
                                int* integrality);

        int Highs_crossover(void* highs, const int num_col, const int num_row,
                                 const double* col_value, const double* col_dual,
                                 const double* row_dual);

        int Highs_getRanging(void* highs,
            double* col_cost_up_value, double* col_cost_up_objective,
            int* col_cost_up_in_var, int* col_cost_up_ou_var,
            double* col_cost_dn_value, double* col_cost_dn_objective,
            int* col_cost_dn_in_var, int* col_cost_dn_ou_var,
            double* col_bound_up_value, double* col_bound_up_objective,
            int* col_bound_up_in_var, int* col_bound_up_ou_var,
            double* col_bound_dn_value, double* col_bound_dn_objective,
            int* col_bound_dn_in_var, int* col_bound_dn_ou_var,
            double* row_bound_up_value, double* row_bound_up_objective,
            int* row_bound_up_in_var, int* row_bound_up_ou_var,
            double* row_bound_dn_value, double* row_bound_dn_objective,
            int* row_bound_dn_in_var, int* row_bound_dn_ou_var);

        void Highs_resetGlobalScheduler(const int blocking);

        // *********************
        // * Deprecated methods*
        // *********************

        #define HighsStatuskError -1
        #define HighsStatuskOk 0
        #define HighsStatuskWarning 1

        int Highs_call(const int num_col, const int num_row,
                             const int num_nz, const double* col_cost,
                             const double* col_lower, const double* col_upper,
                             const double* row_lower, const double* row_upper,
                             const int* a_start, const int* a_index,
                             const double* a_value, double* col_value, double* col_dual,
                             double* row_value, double* row_dual,
                             int* col_basis_status, int* row_basis_status,
                             int* model_status);

        int Highs_runQuiet(void* highs);

        int Highs_setHighsLogfile(void* highs, const void* logfile);

        int Highs_setHighsOutput(void* highs, const void* outputfile);

        int Highs_getIterationCount(const void* highs);

        int Highs_getSimplexIterationCount(const void* highs);

        int Highs_setHighsBoolOptionValue(void* highs, const char* option,
                                                const int value);

        int Highs_setintOptionValue(void* highs, const char* option,
                                               const int value);

        int Highs_setHighsDoubleOptionValue(void* highs, const char* option,
                                                   const double value);

        int Highs_setHighsStringOptionValue(void* highs, const char* option,
                                                   const char* value);

        int Highs_setHighsOptionValue(void* highs, const char* option,
                                            const char* value);

        int Highs_getHighsBoolOptionValue(const void* highs, const char* option,
                                                int* value);

        int Highs_getintOptionValue(const void* highs, const char* option,
                                               int* value);

        int Highs_getHighsDoubleOptionValue(const void* highs, const char* option,
                                                  double* value);

        int Highs_getHighsStringOptionValue(const void* highs, const char* option,
                                                  char* value);

        int Highs_getHighsOptionType(const void* highs, const char* option,
                                           int* type);

        int Highs_resetHighsOptions(void* highs);

        int Highs_getintInfoValue(const void* highs, const char* info,
                                             int* value);

        int Highs_getHighsDoubleInfoValue(const void* highs, const char* info,
                                                double* value);

        int Highs_getNumCols(const void* highs);

        int Highs_getNumRows(const void* highs);

        double Highs_getHighsInfinity(const void* highs);

        double Highs_getHighsRunTime(const void* highs);

        int Highs_setOptionValue(void* highs, const char* option,
                                       const char* value);

        int Highs_getScaledModelStatus(const void* highs);

        void goHighsCallbackHandler(int callback_type, void* data_out, void* user_data);

        static void highsCallbackGateway(int callback_type, const char* message, const HighsCallbackDataOut* data_out, HighsCallbackDataIn* data_in, void* user_data) {
            if (callback_type == 3 || callback_type == 4) {
                goHighsCallbackHandler(callback_type, (void*)data_out, user_data);
            }
        }

        static HighsCCallbackType getHighsCallbackGateway(void) {
            return &highsCallbackGateway;
        }
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

type HiGHSBackend struct{}

func init() {
	mip.RegisterBackend(mip.HiGHS, &HiGHSBackend{})
}

func buildHighsModel(model *mip.Model) (unsafe.Pointer, error) {
	hModel := C.Highs_create()
	if hModel == nil {
		return nil, errors.New("failed to create HiGHS model instance")
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

	for _, v := range variables {
		C.Highs_addCol(
			hModel,
			C.double(objCoeffs[v.ID()]),
			C.double(v.LB()),
			C.double(v.UB()),
			0,
			nil,
			nil,
		)
		if v.Type() == mip.Integer || v.Type() == mip.Binary {
			C.Highs_changeColIntegrality(hModel, C.int(v.ID()), 1)
		}
		if v.Name() != "" {
			cName := C.CString(v.Name())
			C.Highs_passColName(hModel, C.int(v.ID()), cName)
			C.free(unsafe.Pointer(cName))
		}
	}

	sense := C.int(1)
	if model.IsMaximize() {
		sense = -1
	}
	C.Highs_changeObjectiveSense(hModel, sense)

	highsInf := 1e200
	constraints := model.Constraints()
	for _, c := range constraints {
		var lower, upper float64
		switch c.Sense() {
		case 'L':
			lower = -highsInf
			upper = c.RHS()
		case 'G':
			lower = c.RHS()
			upper = highsInf
		case 'E':
			lower = c.RHS()
			upper = c.RHS()
		default:
			lower = -highsInf
			upper = c.RHS()
		}

		nnz := len(c.Terms())
		if nnz == 0 {
			C.Highs_addRow(hModel, C.double(lower), C.double(upper), 0, nil, nil)
			continue
		}

		cIndices := make([]C.int, nnz)
		cValues := make([]C.double, nnz)
		for i, t := range c.Terms() {
			cIndices[i] = C.int(t.Var.ID())
			cValues[i] = C.double(t.Coeff)
		}

		C.Highs_addRow(
			hModel,
			C.double(lower),
			C.double(upper),
			C.int(nnz),
			&cIndices[0],
			&cValues[0],
		)
	}

	for i, c := range constraints {
		if c.Name() != "" {
			cName := C.CString(c.Name())
			C.Highs_passRowName(hModel, C.int(i), cName)
			C.free(unsafe.Pointer(cName))
		}
	}

	// Add Indicators (Linearized)
	for _, ind := range model.Indicators() {
		addLinearizedIndicator(hModel, ind.BinaryVar(), ind.Constraint())
	}

	// Add SOS1 (linearized for binary variables)
	for _, s := range model.SOS1() {
		allBinary := true
		for _, v := range s.Vars() {
			if v.Type() != mip.Binary {
				allBinary = false
				break
			}
		}
		if allBinary && len(s.Vars()) > 0 {
			nnz := len(s.Vars())
			cIndices := make([]C.int, nnz)
			cValues := make([]C.double, nnz)
			for i, v := range s.Vars() {
				cIndices[i] = C.int(v.ID())
				cValues[i] = 1.0
			}
			C.Highs_addRow(hModel, C.double(-1e200), C.double(1.0), C.int(nnz), &cIndices[0], &cValues[0])
		}
	}

	if model.TimeLimit() > 0 {
		cOpt := C.CString("time_limit")
		C.Highs_setDoubleOptionValue(hModel, cOpt, C.double(model.TimeLimit().Seconds()))
		C.free(unsafe.Pointer(cOpt))
	}
	if model.MIPGap() >= 0 {
		cOpt := C.CString("mip_rel_gap")
		C.Highs_setDoubleOptionValue(hModel, cOpt, C.double(model.MIPGap()))
		C.free(unsafe.Pointer(cOpt))
	}
	if model.Threads() > 0 {
		cOpt := C.CString("threads")
		C.Highs_setIntOptionValue(hModel, cOpt, C.int(model.Threads()))
		C.free(unsafe.Pointer(cOpt))
	}

	return hModel, nil
}

func (b *HiGHSBackend) Solve(model *mip.Model) (*mip.Solution, error) {
	model.ClearSolutionPool()

	hModel, err := buildHighsModel(model)
	if err != nil {
		return nil, err
	}
	defer C.Highs_destroy(hModel)

	// Register callbacks for solution pool collection
	modelID := mip.RegisterHiGHSModel(model)
	defer mip.UnregisterHiGHSModel(modelID)

	cbUserData := modelID
	C.Highs_setCallback(hModel, C.getHighsCallbackGateway(), unsafe.Pointer(uintptr(cbUserData)))

	C.Highs_startCallback(hModel, C.kHighsCallbackMipSolution)
	defer C.Highs_stopCallback(hModel, C.kHighsCallbackMipSolution)

	// Set MIPStart (initial solution) if provided
	if len(model.MIPStart()) > 0 {
		variables := model.Variables()
		numCols := len(variables)
		highsInf := float64(C.Highs_getInfinity(hModel))
		colValues := make([]C.double, numCols)
		for i := 0; i < numCols; i++ {
			colValues[i] = C.double(highsInf)
		}
		for v, val := range model.MIPStart() {
			if v.ID() >= 0 && v.ID() < numCols {
				colValues[v.ID()] = C.double(val)
			}
		}
		C.Highs_setSolution(hModel, &colValues[0], nil, nil, nil)
	}

	// Run the Solver
	C.Highs_run(hModel)

	// Get Solver Status
	hStatus := int(C.Highs_getModelStatus(hModel))
	var status mip.Status
	switch hStatus {
	case 7: // kHighsModelStatusOptimal
		status = mip.Optimal
	case 8: // kHighsModelStatusInfeasible
		status = mip.Infeasible
	case 9: // kHighsModelStatusUnboundedOrInfeasible
		status = mip.Infeasible
	case 10: // kHighsModelStatusUnbounded
		status = mip.Unbounded
	case 13: // kHighsModelStatusTimeLimit
		status = mip.TimeLimit
	case 17: // kHighsModelStatusInterrupt
		status = mip.Interrupted
	default:
		status = mip.Error
	}

	solValues := make(map[*mip.Variable]float64)
	solRedCosts := make(map[*mip.Variable]float64)
	solDuals := make(map[*mip.Constraint]float64)
	solSlacks := make(map[*mip.Constraint]float64)
	var stats mip.Stats

	variables := model.Variables()
	constraints := model.Constraints()
	numCols := len(variables)
	numRows := len(constraints)

	if status == mip.Optimal || status == mip.Feasible || status == mip.TimeLimit {
		objVal := float64(C.Highs_getObjectiveValue(hModel))

		colValues := make([]C.double, numCols)
		colDuals := make([]C.double, numCols)
		rowValues := make([]C.double, numRows)
		rowDuals := make([]C.double, numRows)

		var colValPtr, colDualPtr, rowValPtr, rowDualPtr *C.double
		if numCols > 0 {
			colValPtr = &colValues[0]
			colDualPtr = &colDuals[0]
		}
		if numRows > 0 {
			rowValPtr = &rowValues[0]
			rowDualPtr = &rowDuals[0]
		}

		C.Highs_getSolution(hModel, colValPtr, colDualPtr, rowValPtr, rowDualPtr)

		for _, v := range variables {
			solValues[v] = float64(colValues[v.ID()])
			solRedCosts[v] = float64(colDuals[v.ID()])
		}

		for i, c := range constraints {
			solDuals[c] = float64(rowDuals[i])
			solSlacks[c] = c.RHS() - float64(rowValues[i])
		}

		// Append final best solution to solution pool
		finalPoolSolValues := make(map[*mip.Variable]float64)
		for _, v := range variables {
			finalPoolSolValues[v] = solValues[v]
		}
		finalPoolSol := mip.NewSolution(
			status,
			objVal,
			finalPoolSolValues,
			nil,
			nil,
			nil,
			mip.Stats{},
		)
		mip.AppendSolutionToPool(modelID, finalPoolSol)
	}

	// Collect Statistics
	var simplexCount C.int
	C.Highs_getIntInfoValue(hModel, C.CString("simplex_iteration_count"), &simplexCount)

	var ipmCount C.int
	C.Highs_getIntInfoValue(hModel, C.CString("ipm_iteration_count"), &ipmCount)

	var nodeCount C.int64_t
	C.Highs_getInt64InfoValue(hModel, C.CString("mip_node_count"), &nodeCount)

	var bestBound C.double
	C.Highs_getDoubleInfoValue(hModel, C.CString("mip_dual_bound"), &bestBound)

	var gap C.double
	C.Highs_getDoubleInfoValue(hModel, C.CString("mip_gap"), &gap)

	stats = mip.Stats{
		Nodes:             int64(nodeCount),
		SimplexIterations: int64(simplexCount),
		BarrierIterations: int64(ipmCount),
		Runtime:           time.Duration(float64(C.Highs_getRunTime(hModel)) * float64(time.Second)),
		BestBound:         float64(bestBound),
		Gap:               float64(gap),
	}

	objVal := 0.0
	if status == mip.Optimal || status == mip.Feasible || status == mip.TimeLimit {
		objVal = float64(C.Highs_getObjectiveValue(hModel))
	}

	sol := mip.NewSolution(status, objVal, solValues, solRedCosts, solDuals, solSlacks, stats)
	return sol, nil
}

func addLinearizedIndicator(hModel unsafe.Pointer, b *mip.Variable, c *mip.Constraint) {
	if c.Sense() == 'E' {
		c1 := mip.NewConstraintFromTerms(c.Terms(), c.RHS(), 'L')
		c2 := mip.NewConstraintFromTerms(c.Terms(), c.RHS(), 'G')
		addLinearizedIndicator(hModel, b, c1)
		addLinearizedIndicator(hModel, b, c2)
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

	var lower, upper float64
	highsInf := 1e200

	if c.Sense() == 'L' {
		cValues[nnz-1] = C.double(mVal)
		lower = -highsInf
		upper = c.RHS() + mVal
	} else { // 'G'
		cValues[nnz-1] = C.double(-mVal)
		lower = c.RHS() - mVal
		upper = highsInf
	}

	C.Highs_addRow(
		hModel,
		C.double(lower),
		C.double(upper),
		C.int(nnz),
		&cIndices[0],
		&cValues[0],
	)
}

//export goHighsCallbackHandler
func goHighsCallbackHandler(callbackType C.int, dataOut unsafe.Pointer, userData unsafe.Pointer) {
	if userData == nil {
		return
	}
	id := int(uintptr(userData))
	vars, ok := mip.GetActiveVars(id)
	if !ok || dataOut == nil {
		return
	}

	hData := (*C.HighsCallbackDataOut)(dataOut)
	if hData.mip_solution == nil {
		return
	}

	objVal := float64(hData.objective_function_value)
	values := make(map[*mip.Variable]float64)

	solArr := unsafe.Slice(hData.mip_solution, len(vars))
	for idx, v := range vars {
		values[v] = float64(solArr[idx])
	}

	poolSol := mip.NewSolution(
		mip.Feasible,
		objVal,
		values,
		nil,
		nil,
		nil,
		mip.Stats{},
	)

	mip.AppendSolutionToPool(id, poolSol)
}

func (b *HiGHSBackend) WriteLP(model *mip.Model, filename string) error {
	return writeModelUsingHiGHS(model, filename)
}

func (b *HiGHSBackend) WriteMPS(model *mip.Model, filename string) error {
	return writeModelUsingHiGHS(model, filename)
}

func (b *HiGHSBackend) ReadLP(filename string) (*mip.Model, error) {
	return readModelUsingHiGHS(filename)
}

func (b *HiGHSBackend) ReadMPS(filename string) (*mip.Model, error) {
	return readModelUsingHiGHS(filename)
}

func writeModelUsingHiGHS(model *mip.Model, filename string) error {
	hModel, err := buildHighsModel(model)
	if err != nil {
		return err
	}
	defer C.Highs_destroy(hModel)

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	status := int(C.Highs_writeModel(hModel, cFilename))
	if status == -1 {
		return fmt.Errorf("failed to write model to %s (HiGHS error)", filename)
	}
	return nil
}

func readModelUsingHiGHS(filename string) (*mip.Model, error) {
	hModel := C.Highs_create()
	if hModel == nil {
		return nil, fmt.Errorf("failed to create HiGHS instance")
	}
	defer C.Highs_destroy(hModel)

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	status := int(C.Highs_readModel(hModel, cFilename))
	if status == -1 {
		return nil, fmt.Errorf("failed to read model from %s (HiGHS error)", filename)
	}

	numCols := int(C.Highs_getNumCol(hModel))
	numRows := int(C.Highs_getNumRow(hModel))
	numNz := int(C.Highs_getNumNz(hModel))

	m := mip.New("imported_model")

	if numCols == 0 {
		return m, nil
	}

	colCost := make([]C.double, numCols)
	colLower := make([]C.double, numCols)
	colUpper := make([]C.double, numCols)
	integrality := make([]C.int, numCols)

	rowLower := make([]C.double, numRows)
	rowUpper := make([]C.double, numRows)

	var aStart []C.int
	var aIndex []C.int
	var aValue []C.double

	if numNz > 0 {
		aStart = make([]C.int, numCols+1)
		aIndex = make([]C.int, numNz)
		aValue = make([]C.double, numNz)
	} else {
		aStart = make([]C.int, numCols+1)
	}

	var sense C.int
	var offset C.double
	var numColOut C.int
	var numRowOut C.int
	var numNzOut C.int
	var hessianNzOut C.int

	var aStartPtr, aIndexPtr *C.int
	var aValuePtr *C.double
	if numNz > 0 {
		aStartPtr = &aStart[0]
		aIndexPtr = &aIndex[0]
		aValuePtr = &aValue[0]
	} else {
		aStartPtr = &aStart[0]
	}

	var colCostPtr, colLowerPtr, colUpperPtr *C.double
	var integralityPtr *C.int
	if numCols > 0 {
		colCostPtr = &colCost[0]
		colLowerPtr = &colLower[0]
		colUpperPtr = &colUpper[0]
		integralityPtr = &integrality[0]
	}

	var rowLowerPtr, rowUpperPtr *C.double
	if numRows > 0 {
		rowLowerPtr = &rowLower[0]
		rowUpperPtr = &rowUpper[0]
	}

	C.Highs_getModel(
		hModel,
		1, // kHighsMatrixFormatColwise
		1, // q_format
		&numColOut,
		&numRowOut,
		&numNzOut,
		&hessianNzOut,
		&sense,
		&offset,
		colCostPtr,
		colLowerPtr,
		colUpperPtr,
		rowLowerPtr,
		rowUpperPtr,
		aStartPtr,
		aIndexPtr,
		aValuePtr,
		nil,
		nil,
		nil,
		integralityPtr,
	)

	vars := make([]*mip.Variable, numCols)
	highsInf := float64(C.Highs_getInfinity(hModel))

	for j := 0; j < numCols; j++ {
		cName := make([]C.char, 256)
		C.Highs_getColName(hModel, C.int(j), &cName[0])
		goName := C.GoString(&cName[0])

		lb := float64(colLower[j])
		ub := float64(colUpper[j])

		if lb <= -highsInf {
			lb = -math.Inf(1)
		}
		if ub >= highsInf {
			ub = math.Inf(1)
		}

		var vType mip.VarType
		if integrality[j] == 1 {
			if lb == 0.0 && ub == 1.0 {
				vType = mip.Binary
			} else {
				vType = mip.Integer
			}
		} else {
			vType = mip.Continuous
		}

		vars[j] = m.NewVarFromReader(goName, vType, lb, ub)
	}

	var objTerms []mip.Term
	for j, v := range vars {
		cost := float64(colCost[j])
		if cost != 0.0 {
			objTerms = append(objTerms, mip.Term{Coeff: cost, Var: v})
		}
	}
	m.SetObjectiveFromReader(mip.NewExpressionFromReader(objTerms, 0), sense == -1)

	if numRows > 0 {
		rowsCoeffs := make([][]mip.Term, numRows)
		for col := 0; col < numCols; col++ {
			start := int(aStart[col])
			end := numNz
			if col < numCols-1 {
				end = int(aStart[col+1])
			}
			for idx := start; idx < end; idx++ {
				row := int(aIndex[idx])
				val := float64(aValue[idx])
				rowsCoeffs[row] = append(rowsCoeffs[row], mip.Term{
					Coeff: val,
					Var:   vars[col],
				})
			}
		}

		for i := 0; i < numRows; i++ {
			cName := make([]C.char, 256)
			C.Highs_getRowName(hModel, C.int(i), &cName[0])
			goName := C.GoString(&cName[0])

			rLower := float64(rowLower[i])
			rUpper := float64(rowUpper[i])

			var senseChar byte
			var rhsVal float64

			if rLower == rUpper {
				senseChar = 'E'
				rhsVal = rLower
			} else if rLower <= -highsInf && rUpper < highsInf {
				senseChar = 'L'
				rhsVal = rUpper
			} else if rLower > -highsInf && rUpper >= highsInf {
				senseChar = 'G'
				rhsVal = rLower
			} else {
				if len(rowsCoeffs[i]) > 0 {
					m.SubjectTo(mip.NewConstraintFromTerms(rowsCoeffs[i], rUpper, 'L').Named(goName + "_ub"))
					m.SubjectTo(mip.NewConstraintFromTerms(rowsCoeffs[i], rLower, 'G').Named(goName + "_lb"))
				}
				continue
			}

			m.SubjectTo(mip.NewConstraintFromTerms(rowsCoeffs[i], rhsVal, senseChar).Named(goName))
		}
	}

	return m, nil
}
