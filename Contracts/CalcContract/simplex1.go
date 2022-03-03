// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package lp implements routines for solving linear programs.
package main

import (
	"fmt"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// Simplex solves a linear program in standard form using Danzig's Simplex
// algorithm. The standard form of a linear program is:
//  minimize	cᵀ x
//  s.t. 		A*x = b
//  			x >= 0 .
// The input tol sets how close to the optimal solution is found (specifically,
// when the maximal reduced cost is below tol). An error will be returned if the
// problem is infeasible or unbounded. In rare cases, numeric errors can cause
// the Simplex to fail. In this case, an error will be returned along with the
// most recently found feasible solution.
//
// The Convert function can be used to transform a general LP into standard form.
//
// The input matrix A must have at least as many columns as rows, len(c) must
// equal the number of columns of A, and len(b) must equal the number of rows of
// A or Simplex will panic. A must also have full row rank and may not contain any
// columns with all zeros, or Simplex will return an error.
//
// initialBasic can be used to set the initial set of indices for a feasible
// solution to the LP. If an initial feasible solution is not known, initialBasic
// may be nil. If initialBasic is non-nil, len(initialBasic) must equal the number
// of rows of A and must be an actual feasible solution to the LP, otherwise
// Simplex will panic.
//
// A description of the Simplex algorithm can be found in Ch. 8 of
//  Strang, Gilbert. "Linear Algebra and Applications." Academic, New York (1976).
// For a detailed video introduction, see lectures 11-13 of UC Math 352
//  https://www.youtube.com/watch?v=ESzYPFkY3og&index=11&list=PLh464gFUoJWOmBYla3zbZbc4nv2AXez6X.
func SimplexMax1(c []float64, A mat.Matrix, b []float64, tol float64, initialBasic []int) (optF float64, optX []float64, err error) {
	ans, x, _, err := simplex1(initialBasic, c, A, b, tol)
	return ans, x, err
}

func simplex1(initialBasic []int, c []float64, A mat.Matrix, b []float64, tol float64) (float64, []float64, []int, error) {
	err := verifyInputs(initialBasic, c, A, b)
	if err != nil {
		if err == ErrUnbounded {
			return math.Inf(-1), nil, nil, ErrUnbounded
		}
		return math.NaN(), nil, nil, err
	}
	m, n := A.Dims()

	if m == n {
		// Problem is exactly constrained, perform a linear solve.
		bVec := mat.NewVecDense(len(b), b)
		x := make([]float64, n)
		xVec := mat.NewVecDense(n, x)
		err := xVec.SolveVec(A, bVec)
		if err != nil {
			return math.NaN(), nil, nil, ErrSingular
		}
		for _, v := range x {
			if v < 0 {
				return math.NaN(), nil, nil, ErrInfeasible
			}
		}
		f := floats.Dot(x, c)
		return f, x, nil, nil
	}

	// There is at least one optimal solution to the LP which is at the intersection
	// to a set of constraint boundaries. For a standard form LP with m variables
	// and n equality constraints, at least m-n elements of x must equal zero
	// at optimality. The Simplex algorithm solves the standard-form LP by starting
	// at an initial constraint vertex and successively moving to adjacent constraint
	// vertices. At every vertex, the set of non-zero x values is the "basic
	// feasible solution". The list of non-zero x's are maintained in basicIdxs,
	// the respective columns of A are in ab, and the actual non-zero values of
	// x are in xb.
	//
	// The LP is equality constrained such that A * x = b. This can be expanded
	// to
	//  ab * xb + an * xn = b
	// where ab are the columns of a in the basic set, and an are all of the
	// other columns. Since each element of xn is zero by definition, this means
	// that for all feasible solutions xb = ab^-1 * b.
	//
	// Before the simplex algorithm can start, an initial feasible solution must
	// be found. If initialBasic is non-nil a feasible solution has been supplied.
	// Otherwise the "Phase I" problem must be solved to find an initial feasible
	// solution.

	var basicIdxs []int // The indices of the non-zero x values.
	var ab *mat.Dense   // The subset of columns of A listed in basicIdxs.
	var xb []float64    // The non-zero elements of x. xb = ab^-1 b

	if initialBasic != nil {
		// InitialBasic supplied. Panic if incorrect length or infeasible.
		if len(initialBasic) != m {
			panic("lp: incorrect number of initial vectors")
		}
		ab = mat.NewDense(m, len(initialBasic), nil)
		extractColumns(ab, A, initialBasic)
		xb = make([]float64, m)
		err = initializeFromBasic(xb, ab, b)
		if err != nil {
			panic(err)
		}
		basicIdxs = make([]int, len(initialBasic))
		copy(basicIdxs, initialBasic)
	} else {
		// No initial basis supplied. Solve the PhaseI problem.
		basicIdxs, ab, xb, err = findInitialBasic(A, b)
		if err != nil {
			return math.NaN(), nil, nil, err
		}
	}

	// basicIdxs contains the indexes for an initial feasible solution,
	// ab contains the extracted columns of A, and xb contains the feasible
	// solution. All x not in the basic set are 0 by construction.

	// nonBasicIdx is the set of nonbasic variables.
	nonBasicIdx := make([]int, 0, n-m)
	inBasic := make(map[int]struct{})
	for _, v := range basicIdxs {
		inBasic[v] = struct{}{}
	}
	for i := 0; i < n; i++ {
		_, ok := inBasic[i]
		if !ok {
			nonBasicIdx = append(nonBasicIdx, i)
		}
	}

	// cb is the subset of c for the basic variables. an and cn
	// are the equivalents to ab and cb but for the nonbasic variables.
	cb := make([]float64, len(basicIdxs))
	for i, idx := range basicIdxs {
		cb[i] = c[idx]
	}
	cn := make([]float64, len(nonBasicIdx))
	for i, idx := range nonBasicIdx {
		cn[i] = c[idx]
	}
	an := mat.NewDense(m, len(nonBasicIdx), nil)
	extractColumns(an, A, nonBasicIdx)

	bVec := mat.NewVecDense(len(b), b)
	cbVec := mat.NewVecDense(len(cb), cb)

	// Temporary data needed each iteration. (Described later)
	r := make([]float64, n-m)
	move := make([]float64, m)

	// Solve the linear program starting from the initial feasible set. This is
	// the "Phase 2" problem.
	//
	// Algorithm:
	// 1) Compute the "reduced costs" for the non-basic variables. The reduced
	// costs are the lagrange multipliers of the constraints.
	// 	 r = cn - anᵀ * ab¯ᵀ * cb
	// 2) If all of the reduced costs are positive, no improvement is possible,
	// and the solution is optimal (xn can only increase because of
	// non-negativity constraints). Otherwise, the solution can be improved and
	// one element will be exchanged in the basic set.
	// 3) Choose the x_n with the most negative value of r. Call this value xe.
	// This variable will be swapped into the basic set.
	// 4) Increase xe until the next constraint boundary is met. This will happen
	// when the first element in xb becomes 0. The distance xe can increase before
	// a given element in xb becomes negative can be found from
	//	xb = Ab^-1 b - Ab^-1 An xn
	//     = Ab^-1 b - Ab^-1 Ae xe
	//     = bhat + d x_e
	//  xe = bhat_i / - d_i
	// where Ae is the column of A corresponding to xe.
	// The constraining basic index is the first index for which this is true,
	// so remove the element which is min_i (bhat_i / -d_i), assuming d_i is negative.
	// If no d_i is less than 0, then the problem is unbounded.
	// 5) If the new xe is 0 (that is, bhat_i == 0), then this location is at
	// the intersection of several constraints. Use the Bland rule instead
	// of the rule in step 4 to avoid cycling.
	for {
		// Compute reduced costs -- r = cn - anᵀ ab¯ᵀ cb
		var tmp mat.VecDense
		err = tmp.SolveVec(ab.T(), cbVec)
		if err != nil {
			break
		}
		data := make([]float64, n-m)
		tmp2 := mat.NewVecDense(n-m, data)
		tmp2.MulVec(an.T(), &tmp)
		floats.SubTo(r, cn, data)

		// Replace the most negative element in the simplex. If there are no
		// negative entries then the optimal solution has been found.
		maxIdx := floats.MaxIdx(r)
		if r[maxIdx] <= -tol {
			break
		}

		for i, v := range r {
			if math.Abs(v) < rRoundTol {
				r[i] = 0
			}
		}

		// Compute the moving distance.
		err = computeMove(move, maxIdx, A, ab, xb, nonBasicIdx)
		if err != nil {
			if err == ErrUnbounded {
				return math.Inf(-1), nil, nil, ErrUnbounded
			}
			break
		}

		// Replace the basic index along the tightest constraint.
		replace := floats.MinIdx(move)
		if move[replace] <= 0 {
			replace, maxIdx, err = replaceBland(A, ab, xb, basicIdxs, nonBasicIdx, r, move)
			if err != nil {
				if err == ErrUnbounded {
					return math.Inf(-1), nil, nil, ErrUnbounded
				}
				break
			}
		}

		// Replace the constrained basicIdx with the newIdx.
		basicIdxs[replace], nonBasicIdx[maxIdx] = nonBasicIdx[maxIdx], basicIdxs[replace]
		cb[replace], cn[maxIdx] = cn[maxIdx], cb[replace]
		tmpCol1 := mat.Col(nil, replace, ab)
		tmpCol2 := mat.Col(nil, maxIdx, an)
		ab.SetCol(replace, tmpCol2)
		an.SetCol(maxIdx, tmpCol1)

		// Compute the new xb.
		xbVec := mat.NewVecDense(len(xb), xb)
		err = xbVec.SolveVec(ab, bVec)
		if err != nil {
			break
		}
	}
	// Found the optimum successfully or died trying. The basic variables get
	// their values, and the non-basic variables are all zero.
	opt := floats.Dot(cb, xb)
	xopt := make([]float64, n)
	for i, v := range basicIdxs {
		xopt[v] = xb[i]
	}
	return opt, xopt, basicIdxs, err
}

func SolveDualProblem(A mat.Matrix,c,b []float64){
	A = A.T()
	//_,n := A.Dims()
	//zerob := make([]float64,n)
	//for i := range b {
	//	b[i] = 0
	//}
	//a := make([]float64,n * n)
	//for i :=0;i<n*n;i+=1{
	//	a[i] = 0
	//}
	//aMatrix := mat.NewDense(n,n,a)
	//fmt.Println(aMatrix,zerob)
	cNew,aNew,bNew := lp.Convert(b , A, c, nil,nil)
	fmt.Println(aNew,bNew,cNew)
	// for i := range cNew{
	//	 cNew[i] = -cNew[i]
	// }
	opt, x, err := SimplexMax1(cNew, aNew, bNew, 0, nil)
	if err != nil {
		return
		fmt.Println(err)

	}
	fmt.Printf("optDual: %v\n", -opt)
	fmt.Printf("z: %v\n", x)
}

func TransferToStandardMin(c []float64,A mat.Dense) ([]float64,mat.Dense){
	// 形如 A * x >= b
	//     x 》=0
	m,n := A.Dims()
	for i := 0;i<m;i++{
		c = append(c,0)
	}
	srcData := make([]float64, m+n)
	for i := 0; i < m+n; i++ {
		srcData[i] = c[i]
	}
	number := DenseToSlice(A)
	newNumber := make([]float64,m*(m+n))
	counter :=0
	for i :=0;i<m;i++{
		var j int
		for j= 0; j<n;j++{
			newNumber[i*(m+n)+j] = number[counter]
			counter++
		}

		newNumber[i*(m+n)+j+i] = -1

	}
	newA := mat.NewDense(m,m+n,newNumber)
	return srcData,*newA



}

func TransferToStandardMax(c []float64,A mat.Dense) ([]float64,mat.Dense){
	// 形如 A * x >= b
	//     x 》=0
	m,n := A.Dims()
	for i := 0;i<m;i++{
		c = append(c,0)
	}
	srcData := make([]float64, m+n)
	for i := 0; i < m+n; i++ {
		srcData[i] = c[i]
	}
	number := DenseToSlice(A)
	newNumber := make([]float64,m*(m+n))
	counter :=0
	for i :=0;i<m;i++{
		var j int
		for j= 0; j<n;j++{
			newNumber[i*(m+n)+j] = number[counter]
			counter++
		}

		newNumber[i*(m+n)+j+i] = 1

	}
	newA := mat.NewDense(m,m+n,newNumber)
	return srcData,*newA
}

func MatrixToDense(A mat.Matrix) mat.Dense{
	m,n := A.Dims()
	num := MatrixToSlice(A)
	return *mat.NewDense(m,n,num)
}

func CalcAuxProblem(){


}