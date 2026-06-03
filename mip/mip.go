package mip

import (
	"math"
	"time"
)

// --- Constants and Types ---

type Status int

const (
	Optimal Status = iota
	Feasible
	Infeasible
	Unbounded
	TimeLimit
	Interrupted
	Error
)

type VarType byte

const (
	Continuous VarType = 'C'
	Integer    VarType = 'I'
	Binary     VarType = 'B'
)

type SolverType int

const (
	HiGHS SolverType = iota
	CBC
	Gurobi
	SCIP
)

// Inf represents infinity in bounds
var Inf = math.Inf(1)

type Expr interface {
	Terms() []Term
	Constant() float64
	Eq(rhs any) *Constraint
	Leq(rhs any) *Constraint
	Geq(rhs any) *Constraint
}

type LazyConstraintCallback func(sol *Solution) []*Constraint
type CutCallback func(sol *Solution) []*Constraint

// --- Core Structs ---

type Model struct {
	name          string
	variables     []*Variable
	constraints   []*Constraint
	objective     Expr
	isMaximize    bool
	solverType    SolverType
	timeLimit     time.Duration
	mipGap        float64
	threads       int
	sos1          []SOS
	sos2          []SOS
	indicators    []IndicatorConstraint
	lazyCallback  LazyConstraintCallback
	cutCallback   CutCallback
	solutionPool  []*Solution
	mipStart      map[*Variable]float64
}

func (m *Model) AddLazyConstraintCallback(cb LazyConstraintCallback) {
	m.lazyCallback = cb
}

func (m *Model) AddCutCallback(cb CutCallback) {
	m.cutCallback = cb
}

func (m *Model) validateVar(v *Variable) {
	if v != nil && v.model != m {
		panic("variable belongs to a different model")
	}
}

func (m *Model) validateExpr(expr Expr) {
	if expr == nil {
		return
	}
	for _, t := range expr.Terms() {
		m.validateVar(t.Var)
	}
}

func (m *Model) SetMIPStart(start map[*Variable]float64) {
	for v := range start {
		m.validateVar(v)
	}
	m.mipStart = start
}

func (m *Model) SolutionPool() []*Solution {
	return m.solutionPool
}

type Variable struct {
	model     *Model
	id        int
	name      string
	vType     VarType
	lb        float64
	ub        float64
	termArray [1]Term
}

type Term struct {
	Coeff float64
	Var   *Variable
}

type Expression struct {
	termsList []Term
	constVal  float64
}

type Constraint struct {
	name   string
	termsList []Term
	rhs    float64
	sense  byte // 'L' for <=, 'G' for >=, 'E' for ==
}

type SOS struct {
	vars []*Variable
}

type IndicatorConstraint struct {
	binaryVar  *Variable
	constraint *Constraint
}

type Stats struct {
	Nodes             int64
	SimplexIterations int64
	BarrierIterations int64
	Runtime           time.Duration
	BestBound         float64
	Gap               float64
}

type Solution struct {
	status    Status
	objective float64
	values    map[*Variable]float64
	redCosts  map[*Variable]float64
	duals     map[*Constraint]float64
	slacks    map[*Constraint]float64
	stats     Stats
}

// --- Constructors ---

// New creates a new Model with the given name.
func New(name string) *Model {
	return &Model{
		name:       name,
		solverType: HiGHS, // default solver
		mipGap:     -1.0,  // unset indicator
		threads:    -1,    // unset indicator
	}
}

// NewModel is an alias for New.
func NewModel(name string) *Model {
	return New(name)
}

// --- Model Configuration ---

func (m *Model) SetSolver(s SolverType) {
	m.solverType = s
}

func (m *Model) SetTimeLimit(d time.Duration) {
	m.timeLimit = d
}

func (m *Model) SetMIPGap(gap float64) {
	m.mipGap = gap
}

func (m *Model) SetThreads(t int) {
	m.threads = t
}

// --- Single Variables ---

func (m *Model) newVar(name string, vType VarType, lb, ub float64) *Variable {
	v := &Variable{
		model: m,
		id:    len(m.variables),
		name:  name,
		vType: vType,
		lb:    lb,
		ub:    ub,
	}
	v.termArray[0] = Term{Coeff: 1.0, Var: v}
	m.variables = append(m.variables, v)
	return v
}

func (m *Model) Binary(name string) *Variable {
	return m.newVar(name, Binary, 0, 1)
}

func (m *Model) Integer(name string, lb, ub float64) *Variable {
	return m.newVar(name, Integer, lb, ub)
}

func (m *Model) Continuous(name string, lb, ub float64) *Variable {
	return m.newVar(name, Continuous, lb, ub)
}

// --- Bulk Variables (1D) ---

func (m *Model) BinaryVars(name string, n int) []*Variable {
	vars := make([]*Variable, n)
	for i := 0; i < n; i++ {
		vars[i] = m.Binary(name)
	}
	return vars
}

func (m *Model) IntegerVars(name string, n int, lb, ub float64) []*Variable {
	vars := make([]*Variable, n)
	for i := 0; i < n; i++ {
		vars[i] = m.Integer(name, lb, ub)
	}
	return vars
}

func (m *Model) ContinuousVars(name string, n int, lb, ub float64) []*Variable {
	vars := make([]*Variable, n)
	for i := 0; i < n; i++ {
		vars[i] = m.Continuous(name, lb, ub)
	}
	return vars
}

// --- Objectives ---

func (m *Model) Minimize(expr Expr) {
	m.validateExpr(expr)
	m.objective = expr
	m.isMaximize = false
}

func (m *Model) Maximize(expr Expr) {
	m.validateExpr(expr)
	m.objective = expr
	m.isMaximize = true
}

// --- Constraints ---

func (m *Model) SubjectTo(constrs ...*Constraint) {
	for _, c := range constrs {
		if c != nil {
			for _, t := range c.termsList {
				m.validateVar(t.Var)
			}
		}
	}
	m.constraints = append(m.constraints, constrs...)
}

// --- Special Constraints ---

func (m *Model) AddSOS1(vars ...*Variable) {
	for _, v := range vars {
		m.validateVar(v)
	}
	m.sos1 = append(m.sos1, SOS{vars: vars})
}

func (m *Model) AddSOS2(vars ...*Variable) {
	for _, v := range vars {
		m.validateVar(v)
	}
	m.sos2 = append(m.sos2, SOS{vars: vars})
}

func (m *Model) Indicator(b *Variable, c *Constraint) {
	m.validateVar(b)
	if c != nil {
		for _, t := range c.termsList {
			m.validateVar(t.Var)
		}
	}
	m.indicators = append(m.indicators, IndicatorConstraint{
		binaryVar:  b,
		constraint: c,
	})
}

// --- Solution API Accessors ---

func (s *Solution) Objective() float64 {
	return s.objective
}

func (s *Solution) Value(v *Variable) float64 {
	if s.values == nil {
		return 0
	}
	return s.values[v]
}

func (s *Solution) ReducedCost(v *Variable) float64 {
	if s.redCosts == nil {
		return 0
	}
	return s.redCosts[v]
}

func (s *Solution) Dual(c *Constraint) float64 {
	if s.duals == nil {
		return 0
	}
	return s.duals[c]
}

func (s *Solution) Slack(c *Constraint) float64 {
	if s.slacks == nil {
		return 0
	}
	return s.slacks[c]
}

func (s *Solution) Status() Status {
	return s.status
}

func (s *Solution) Stats() Stats {
	return s.stats
}

func (s *Solution) VariablesValues() map[*Variable]float64 {
	return s.values
}

// --- Introspection ---

func (m *Model) Variables() []*Variable {
	return m.variables
}

func (m *Model) Constraints() []*Constraint {
	return m.constraints
}

func (m *Model) Objective() Expr {
	return m.objective
}

func (m *Model) IsMaximize() bool {
	return m.isMaximize
}

func (m *Model) MIPStart() map[*Variable]float64 {
	return m.mipStart
}

func (m *Model) TimeLimit() time.Duration {
	return m.timeLimit
}

func (m *Model) MIPGap() float64 {
	return m.mipGap
}

func (m *Model) Threads() int {
	return m.threads
}

func (m *Model) SOS1() []SOS {
	return m.sos1
}

func (m *Model) SOS2() []SOS {
	return m.sos2
}

func (m *Model) Indicators() []IndicatorConstraint {
	return m.indicators
}

func (m *Model) ClearSolutionPool() {
	m.solutionPool = nil
}

func (m *Model) AddToSolutionPool(sol *Solution) {
	m.solutionPool = append(m.solutionPool, sol)
}

func (v *Variable) ID() int {
	return v.id
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Type() VarType {
	return v.vType
}

func (v *Variable) LB() float64 {
	return v.lb
}

func (v *Variable) UB() float64 {
	return v.ub
}

func (c *Constraint) Name() string {
	return c.name
}

func (c *Constraint) Terms() []Term {
	return c.termsList
}

func (c *Constraint) RHS() float64 {
	return c.rhs
}

func (c *Constraint) Sense() byte {
	return c.sense
}

func (s *SOS) Vars() []*Variable {
	return s.vars
}

func (ind *IndicatorConstraint) BinaryVar() *Variable {
	return ind.binaryVar
}

func (ind *IndicatorConstraint) Constraint() *Constraint {
	return ind.constraint
}

func NewSolution(status Status, objective float64, values map[*Variable]float64, redCosts map[*Variable]float64, duals map[*Constraint]float64, slacks map[*Constraint]float64, stats Stats) *Solution {
	return &Solution{
		status:    status,
		objective: objective,
		values:    values,
		redCosts:  redCosts,
		duals:     duals,
		slacks:    slacks,
		stats:     stats,
	}
}

func (m *Model) NewVarFromReader(name string, vType VarType, lb, ub float64) *Variable {
	return m.newVar(name, vType, lb, ub)
}

func NewExpressionFromReader(terms []Term, constant float64) Expr {
	return Expression{termsList: terms, constVal: constant}
}

func (m *Model) SetObjectiveFromReader(expr Expr, isMaximize bool) {
	m.objective = expr
	m.isMaximize = isMaximize
}

func NewConstraintFromTerms(terms []Term, rhs float64, sense byte) *Constraint {
	return &Constraint{
		termsList: terms,
		rhs:       rhs,
		sense:     sense,
	}
}
