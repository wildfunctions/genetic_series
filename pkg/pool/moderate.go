package pool

import (
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

func init() {
	Register("moderate", func() Pool { return &ModeratePool{} })
}

// ModeratePool extends conservative with powers of 2/3 as leaves,
// sqrt as unary, and power as binary.
type ModeratePool struct{}

func (p *ModeratePool) Name() string { return "moderate" }

func (p *ModeratePool) RandomLeaf(rng *rand.Rand) expr.ExprNode {
	r := rng.Float64()
	switch {
	case r < 0.35:
		return &expr.VarNode{}
	case r < 0.75:
		return &expr.ConstNode{Val: int64(rng.Intn(10) + 1)}
	case r < 0.875:
		// powers of 2: 2, 4, 8, 16
		exp := rng.Intn(4) + 1
		val := int64(1) << uint(exp)
		return &expr.ConstNode{Val: val}
	default:
		// powers of 3: 3, 9, 27
		vals := []int64{3, 9, 27}
		return &expr.ConstNode{Val: vals[rng.Intn(len(vals))]}
	}
}

var moderateUnary = []expr.UnaryOp{
	expr.OpFactorial,
	expr.OpAltSign,
	expr.OpNeg,
	expr.OpSqrt,
}

func (p *ModeratePool) RandomUnary(rng *rand.Rand) expr.UnaryOp {
	return moderateUnary[rng.Intn(len(moderateUnary))]
}

var moderateBinary = []expr.BinaryOp{
	expr.OpAdd,
	expr.OpSub,
	expr.OpMul,
	expr.OpDiv,
	expr.OpPow,
}

func (p *ModeratePool) RandomBinary(rng *rand.Rand) expr.BinaryOp {
	return moderateBinary[rng.Intn(len(moderateBinary))]
}

func (p *ModeratePool) RandomTree(rng *rand.Rand, maxDepth int) expr.ExprNode {
	return randomTree(p, rng, maxDepth)
}
