package pool

import (
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

func init() {
	Register("kitchensink", func() Pool { return &KitchenSinkPool{} })
}

// KitchenSinkPool extends moderate with trig, ln, floor, ceil.
type KitchenSinkPool struct{}

func (p *KitchenSinkPool) Name() string { return "kitchensink" }

func (p *KitchenSinkPool) RandomLeaf(rng *rand.Rand) expr.ExprNode {
	r := rng.Float64()
	switch {
	case r < 0.35:
		return &expr.VarNode{}
	case r < 0.75:
		return &expr.ConstNode{Val: int64(rng.Intn(10) + 1)}
	case r < 0.875:
		exp := rng.Intn(4) + 1
		val := int64(1) << uint(exp)
		return &expr.ConstNode{Val: val}
	default:
		vals := []int64{3, 9, 27}
		return &expr.ConstNode{Val: vals[rng.Intn(len(vals))]}
	}
}

var kitchenSinkUnary = []expr.UnaryOp{
	expr.OpFactorial,
	expr.OpAltSign,
	expr.OpNeg,
	expr.OpDoubleFactorial,
	expr.OpFibonacci,
	expr.OpSqrt,
	expr.OpSin,
	expr.OpCos,
	expr.OpLn,
	expr.OpFloor,
	expr.OpCeil,
}

func (p *KitchenSinkPool) RandomUnary(rng *rand.Rand) expr.UnaryOp {
	return kitchenSinkUnary[rng.Intn(len(kitchenSinkUnary))]
}

var kitchenSinkBinary = []expr.BinaryOp{
	expr.OpAdd,
	expr.OpSub,
	expr.OpMul,
	expr.OpDiv,
	expr.OpPow,
	expr.OpBinomial,
}

func (p *KitchenSinkPool) RandomBinary(rng *rand.Rand) expr.BinaryOp {
	return kitchenSinkBinary[rng.Intn(len(kitchenSinkBinary))]
}

func (p *KitchenSinkPool) RandomTree(rng *rand.Rand, maxDepth int) expr.ExprNode {
	return randomTree(p, rng, maxDepth)
}
