package expr

import (
	"math"
	"math/big"
	"testing"
)

func assertEvalF64(t *testing.T, node ExprNode, n float64, expected float64, tol float64) {
	t.Helper()
	got, ok := node.EvalF64(n)
	if !ok {
		t.Fatalf("EvalF64 returned ok=false for n=%v", n)
	}
	if math.Abs(got-expected) > tol {
		t.Errorf("EvalF64(n=%v) = %v, want %v (tol=%v)", n, got, expected, tol)
	}
}

// TestEvalF64_MatchesBigFloat verifies float64 and big.Float paths produce the
// same results (within float64 epsilon) for a variety of expression trees.
func TestEvalF64_MatchesBigFloat(t *testing.T) {
	trees := []struct {
		name string
		node ExprNode
	}{
		{"var", &VarNode{}},
		{"const7", &ConstNode{Val: 7}},
		{"neg(n)", &UnaryNode{Op: OpNeg, Child: &VarNode{}}},
		{"5!", &UnaryNode{Op: OpFactorial, Child: &ConstNode{Val: 5}}},
		{"(-1)^n", &UnaryNode{Op: OpAltSign, Child: &VarNode{}}},
		{"5!!", &UnaryNode{Op: OpDoubleFactorial, Child: &ConstNode{Val: 5}}},
		{"fib(10)", &UnaryNode{Op: OpFibonacci, Child: &ConstNode{Val: 10}}},
		{"sin(n)", &UnaryNode{Op: OpSin, Child: &VarNode{}}},
		{"cos(n)", &UnaryNode{Op: OpCos, Child: &VarNode{}}},
		{"ln(n)", &UnaryNode{Op: OpLn, Child: &VarNode{}}},
		{"floor(n/3)", &UnaryNode{Op: OpFloor, Child: &BinaryNode{
			Op: OpDiv, Left: &VarNode{}, Right: &ConstNode{Val: 3}}}},
		{"ceil(n/3)", &UnaryNode{Op: OpCeil, Child: &BinaryNode{
			Op: OpDiv, Left: &VarNode{}, Right: &ConstNode{Val: 3}}}},
		{"abs(n)", &UnaryNode{Op: OpAbs, Child: &VarNode{}}},
		{"sqrt(n)", &UnaryNode{Op: OpSqrt, Child: &VarNode{}}},
		{"n+2", &BinaryNode{Op: OpAdd, Left: &VarNode{}, Right: &ConstNode{Val: 2}}},
		{"n-2", &BinaryNode{Op: OpSub, Left: &VarNode{}, Right: &ConstNode{Val: 2}}},
		{"n*3", &BinaryNode{Op: OpMul, Left: &VarNode{}, Right: &ConstNode{Val: 3}}},
		{"n/4", &BinaryNode{Op: OpDiv, Left: &VarNode{}, Right: &ConstNode{Val: 4}}},
		{"2^n", &BinaryNode{Op: OpPow, Left: &ConstNode{Val: 2}, Right: &VarNode{}}},
		{"C(10,n)", &BinaryNode{Op: OpBinomial, Left: &ConstNode{Val: 10}, Right: &VarNode{}}},
		{"1/n!", &BinaryNode{Op: OpDiv, Left: &ConstNode{Val: 1},
			Right: &UnaryNode{Op: OpFactorial, Child: &VarNode{}}}},
	}

	const prec = 512
	testNs := []float64{1, 2, 3, 4, 5, 7, 10}

	for _, tc := range trees {
		t.Run(tc.name, func(t *testing.T) {
			for _, nv := range testNs {
				f64val, f64ok := tc.node.EvalF64(nv)
				bfval, bfok := tc.node.Eval(
					new(big.Float).SetPrec(prec).SetFloat64(nv), prec)

				if f64ok != bfok {
					t.Errorf("n=%v: ok mismatch f64=%v bf=%v", nv, f64ok, bfok)
					continue
				}
				if !f64ok {
					continue
				}
				bfF64, _ := bfval.Float64()
				diff := math.Abs(f64val - bfF64)
				tol := math.Max(math.Abs(bfF64)*1e-12, 1e-12)
				if diff > tol {
					t.Errorf("n=%v: f64=%v bf=%v diff=%v", nv, f64val, bfF64, diff)
				}
			}
		})
	}
}

func TestEvalF64_Factorial(t *testing.T) {
	node := &UnaryNode{Op: OpFactorial, Child: &ConstNode{Val: 5}}
	assertEvalF64(t, node, 0, 120, 0)

	node = &UnaryNode{Op: OpFactorial, Child: &ConstNode{Val: 0}}
	assertEvalF64(t, node, 0, 1, 0)

	// Negative factorial should fail
	node = &UnaryNode{Op: OpFactorial, Child: &ConstNode{Val: -1}}
	_, ok := node.EvalF64(0)
	if ok {
		t.Error("Factorial(-1) should return ok=false")
	}

	// 170! should succeed (last finite float64 factorial)
	node = &UnaryNode{Op: OpFactorial, Child: &ConstNode{Val: 170}}
	_, ok = node.EvalF64(0)
	if !ok {
		t.Error("170! should succeed in float64")
	}
}

func TestEvalF64_DivisionByZero(t *testing.T) {
	node := &BinaryNode{Op: OpDiv, Left: &ConstNode{Val: 1}, Right: &ConstNode{Val: 0}}
	_, ok := node.EvalF64(0)
	if ok {
		t.Error("Division by zero should return ok=false")
	}
}

func TestEvalF64_Fibonacci(t *testing.T) {
	cases := []struct {
		n    int64
		want float64
	}{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 3}, {5, 5}, {6, 8}, {10, 55},
	}
	for _, tc := range cases {
		node := &UnaryNode{Op: OpFibonacci, Child: &ConstNode{Val: tc.n}}
		assertEvalF64(t, node, 0, tc.want, 0)
	}
}

func TestEvalF64_Pow(t *testing.T) {
	// 2^3 = 8
	node := &BinaryNode{Op: OpPow, Left: &ConstNode{Val: 2}, Right: &ConstNode{Val: 3}}
	assertEvalF64(t, node, 0, 8, 0)

	// 2^(-1) = 0.5
	node = &BinaryNode{Op: OpPow, Left: &ConstNode{Val: 2}, Right: &ConstNode{Val: -1}}
	assertEvalF64(t, node, 0, 0.5, 1e-15)

	// exp > 20 should fail
	node = &BinaryNode{Op: OpPow, Left: &ConstNode{Val: 2}, Right: &ConstNode{Val: 21}}
	_, ok := node.EvalF64(0)
	if ok {
		t.Error("2^21 should fail (intPowF64 cap)")
	}
}

func TestEvalF64_Binomial(t *testing.T) {
	// C(5,2) = 10
	node := &BinaryNode{Op: OpBinomial, Left: &ConstNode{Val: 5}, Right: &ConstNode{Val: 2}}
	assertEvalF64(t, node, 0, 10, 0)

	// C(10,0) = 1
	node = &BinaryNode{Op: OpBinomial, Left: &ConstNode{Val: 10}, Right: &ConstNode{Val: 0}}
	assertEvalF64(t, node, 0, 1, 0)
}
