package strategy

import (
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
	"github.com/wildfunctions/genetic_series/pkg/series"
)

// CrossoverCandidates performs subtree crossover between two candidates,
// returning two new offspring. Both numerator and denominator trees are crossed.
func CrossoverCandidates(a, b *series.Candidate, rng *rand.Rand) (*series.Candidate, *series.Candidate) {
	c1 := a.Clone()
	c2 := b.Clone()

	// Cross numerators
	c1.Numerator, c2.Numerator = crossoverTrees(c1.Numerator, c2.Numerator, rng)

	// Cross denominators
	c1.Denominator, c2.Denominator = crossoverTrees(c1.Denominator, c2.Denominator, rng)

	return c1, c2
}

// crossoverTrees swaps random subtrees between two expression trees.
func crossoverTrees(a, b expr.ExprNode, rng *rand.Rand) (expr.ExprNode, expr.ExprNode) {
	nodesA := collectNodes(a)
	nodesB := collectNodes(b)

	if len(nodesA) == 0 || len(nodesB) == 0 {
		return a, b
	}

	idxA := rng.Intn(len(nodesA))
	idxB := rng.Intn(len(nodesB))

	// Swap the subtrees
	*nodesA[idxA], *nodesB[idxB] = *nodesB[idxB], *nodesA[idxA]

	return a, b
}
