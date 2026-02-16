package expr

import "math"

func (v *VarNode) NodeCount() int { return 1 }
func (c *ConstNode) NodeCount() int { return 1 }
func (u *UnaryNode) NodeCount() int { return 1 + u.Child.NodeCount() }
func (b *BinaryNode) NodeCount() int {
	return 1 + b.Left.NodeCount() + b.Right.NodeCount()
}

func (v *VarNode) Depth() int { return 1 }
func (c *ConstNode) Depth() int { return 1 }
func (u *UnaryNode) Depth() int { return 1 + u.Child.Depth() }
func (b *BinaryNode) Depth() int {
	ld := b.Left.Depth()
	rd := b.Right.Depth()
	if ld > rd {
		return 1 + ld
	}
	return 1 + rd
}

// WeightedComplexity returns a complexity score with heavier weight for
// operations that are more "expensive" (factorial, trig, etc.).
func WeightedComplexity(node ExprNode) float64 {
	switch n := node.(type) {
	case *VarNode:
		return 1.0
	case *ConstNode:
		v := n.Val
		if v < 0 {
			v = -v
		}
		if v <= 10 {
			return 1.0
		}
		return 1.0 + math.Log10(float64(v))
	case *UnaryNode:
		w := unaryWeight(n.Op)
		return w + WeightedComplexity(n.Child)
	case *BinaryNode:
		w := binaryWeight(n.Op)
		return w + WeightedComplexity(n.Left) + WeightedComplexity(n.Right)
	default:
		return 1.0
	}
}

func unaryWeight(op UnaryOp) float64 {
	switch op {
	case OpNeg, OpAbs:
		return 1.0
	case OpFactorial, OpAltSign:
		return 2.0
	case OpDoubleFactorial, OpFibonacci:
		return 3.0
	case OpSin, OpCos, OpLn:
		return 3.0
	case OpFloor, OpCeil:
		return 2.0
	case OpSqrt:
		return 2.0
	default:
		return 2.0
	}
}

func binaryWeight(op BinaryOp) float64 {
	switch op {
	case OpAdd, OpSub:
		return 1.0
	case OpMul, OpDiv:
		return 1.5
	case OpPow:
		return 2.0
	case OpBinomial:
		return 3.0
	default:
		return 1.5
	}
}
