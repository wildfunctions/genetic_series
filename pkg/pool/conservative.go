package pool

import (
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

func init() {
	Register("conservative", func() Pool { return &ConservativePool{} })
}

// ConservativePool provides basic building blocks: n, ints 1-10,
// factorial, (-1)^n, negation, and basic arithmetic.
type ConservativePool struct{}

func (p *ConservativePool) Name() string { return "conservative" }

func (p *ConservativePool) RandomLeaf(rng *rand.Rand) expr.ExprNode {
	if rng.Float64() < 0.4 {
		return &expr.VarNode{}
	}
	return &expr.ConstNode{Val: int64(rng.Intn(10) + 1)}
}

var conservativeUnary = []expr.UnaryOp{
	expr.OpFactorial,
	expr.OpAltSign,
	expr.OpNeg,
}

func (p *ConservativePool) RandomUnary(rng *rand.Rand) expr.UnaryOp {
	return conservativeUnary[rng.Intn(len(conservativeUnary))]
}

var conservativeBinary = []expr.BinaryOp{
	expr.OpAdd,
	expr.OpSub,
	expr.OpMul,
	expr.OpDiv,
}

func (p *ConservativePool) RandomBinary(rng *rand.Rand) expr.BinaryOp {
	return conservativeBinary[rng.Intn(len(conservativeBinary))]
}

func (p *ConservativePool) RandomTree(rng *rand.Rand, maxDepth int) expr.ExprNode {
	return randomTree(p, rng, maxDepth)
}
