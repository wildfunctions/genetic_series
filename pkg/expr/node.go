package expr

import "math/big"

// ExprNode is the interface for all expression tree nodes.
type ExprNode interface {
	Eval(n *big.Float, prec uint) (*big.Float, bool)
	EvalF64(n float64) (float64, bool)
	String() string
	LaTeX() string
	Clone() ExprNode
	NodeCount() int
	Depth() int
}

// UnaryOp identifies a unary operation.
type UnaryOp int

const (
	OpNeg UnaryOp = iota
	OpFactorial
	OpAltSign // (-1)^n
	OpDoubleFactorial
	OpFibonacci
	OpSin
	OpCos
	OpLn
	OpFloor
	OpCeil
	OpAbs
	OpSqrt
)

// BinaryOp identifies a binary operation.
type BinaryOp int

const (
	OpAdd BinaryOp = iota
	OpSub
	OpMul
	OpDiv
	OpPow
	OpBinomial // C(a, b)
)

// VarNode represents the variable n.
type VarNode struct{}

// ConstNode represents an integer constant.
type ConstNode struct {
	Val int64
}

// UnaryNode applies a unary operation to a child expression.
type UnaryNode struct {
	Op    UnaryOp
	Child ExprNode
}

// BinaryNode applies a binary operation to two child expressions.
type BinaryNode struct {
	Op          BinaryOp
	Left, Right ExprNode
}
