package mip

import (
	"reflect"
)

// --- VarMatrix (2D Variable Container) ---

type VarMatrix struct {
	name string
	rows int
	cols int
	vars []*Variable
}

func (vm *VarMatrix) At(i, j int) *Variable {
	if i < 0 || i >= vm.rows || j < 0 || j >= vm.cols {
		panic("index out of matrix bounds")
	}
	return vm.vars[i*vm.cols + j]
}

func (vm *VarMatrix) Variables() []*Variable {
	return vm.vars
}

// --- VarTensor (N-Dimensional Variable Container) ---

type VarTensor struct {
	name    string
	dims    []int
	strides []int
	vars    []*Variable
}

func (vt *VarTensor) At(indices ...int) *Variable {
	if len(indices) != len(vt.dims) {
		panic("incorrect number of dimensions for tensor access")
	}
	flatIdx := 0
	for i, idx := range indices {
		if idx < 0 || idx >= vt.dims[i] {
			panic("index out of tensor bounds")
		}
		flatIdx += idx * vt.strides[i]
	}
	return vt.vars[flatIdx]
}

func (vt *VarTensor) Variables() []*Variable {
	return vt.vars
}

// --- Model Integration for Matrix and Tensor creation ---

func resolveDim(d any) int {
	switch val := d.(type) {
	case int:
		return val
	default:
		rv := reflect.ValueOf(d)
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			return rv.Len()
		}
		panic("dimension must be an int or a slice")
	}
}

// BinaryMatrix creates a 2D matrix of binary variables.
func (m *Model) BinaryMatrix(name string, rows, cols any) *VarMatrix {
	r := resolveDim(rows)
	c := resolveDim(cols)
	vars := make([]*Variable, r*c)
	for i := 0; i < r*c; i++ {
		vars[i] = m.Binary(name)
	}
	return &VarMatrix{
		name: name,
		rows: r,
		cols: c,
		vars: vars,
	}
}

// IntegerMatrix creates a 2D matrix of integer variables.
func (m *Model) IntegerMatrix(name string, rows, cols any, lb, ub float64) *VarMatrix {
	r := resolveDim(rows)
	c := resolveDim(cols)
	vars := make([]*Variable, r*c)
	for i := 0; i < r*c; i++ {
		vars[i] = m.Integer(name, lb, ub)
	}
	return &VarMatrix{
		name: name,
		rows: r,
		cols: c,
		vars: vars,
	}
}

// ContinuousMatrix creates a 2D matrix of continuous variables.
func (m *Model) ContinuousMatrix(name string, rows, cols any, lb, ub float64) *VarMatrix {
	r := resolveDim(rows)
	c := resolveDim(cols)
	vars := make([]*Variable, r*c)
	for i := 0; i < r*c; i++ {
		vars[i] = m.Continuous(name, lb, ub)
	}
	return &VarMatrix{
		name: name,
		rows: r,
		cols: c,
		vars: vars,
	}
}

// BinaryTensor creates an N-dimensional tensor of binary variables.
func (m *Model) BinaryTensor(name string, dims ...any) *VarTensor {
	resolvedDims := make([]int, len(dims))
	totalSize := 1
	for i, d := range dims {
		resolvedDims[i] = resolveDim(d)
		totalSize *= resolvedDims[i]
	}

	strides := make([]int, len(resolvedDims))
	currentStride := 1
	for i := len(resolvedDims) - 1; i >= 0; i-- {
		strides[i] = currentStride
		currentStride *= resolvedDims[i]
	}

	vars := make([]*Variable, totalSize)
	for i := 0; i < totalSize; i++ {
		vars[i] = m.Binary(name)
	}

	return &VarTensor{
		name:    name,
		dims:    resolvedDims,
		strides: strides,
		vars:    vars,
	}
}

// IntegerTensor creates an N-dimensional tensor of integer variables.
func (m *Model) IntegerTensor(name string, lb, ub float64, dims ...any) *VarTensor {
	resolvedDims := make([]int, len(dims))
	totalSize := 1
	for i, d := range dims {
		resolvedDims[i] = resolveDim(d)
		totalSize *= resolvedDims[i]
	}

	strides := make([]int, len(resolvedDims))
	currentStride := 1
	for i := len(resolvedDims) - 1; i >= 0; i-- {
		strides[i] = currentStride
		currentStride *= resolvedDims[i]
	}

	vars := make([]*Variable, totalSize)
	for i := 0; i < totalSize; i++ {
		vars[i] = m.Integer(name, lb, ub)
	}

	return &VarTensor{
		name:    name,
		dims:    resolvedDims,
		strides: strides,
		vars:    vars,
	}
}

// ContinuousTensor creates an N-dimensional tensor of continuous variables.
func (m *Model) ContinuousTensor(name string, lb, ub float64, dims ...any) *VarTensor {
	resolvedDims := make([]int, len(dims))
	totalSize := 1
	for i, d := range dims {
		resolvedDims[i] = resolveDim(d)
		totalSize *= resolvedDims[i]
	}

	strides := make([]int, len(resolvedDims))
	currentStride := 1
	for i := len(resolvedDims) - 1; i >= 0; i-- {
		strides[i] = currentStride
		currentStride *= resolvedDims[i]
	}

	vars := make([]*Variable, totalSize)
	for i := 0; i < totalSize; i++ {
		vars[i] = m.Continuous(name, lb, ub)
	}

	return &VarTensor{
		name:    name,
		dims:    resolvedDims,
		strides: strides,
		vars:    vars,
	}
}
