package pool

import (
	"fmt"
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

// Pool provides random building blocks for constructing expression trees.
type Pool interface {
	Name() string
	RandomLeaf(rng *rand.Rand) expr.ExprNode
	RandomUnary(rng *rand.Rand) expr.UnaryOp
	RandomBinary(rng *rand.Rand) expr.BinaryOp
	RandomTree(rng *rand.Rand, maxDepth int) expr.ExprNode
}

var registry = map[string]func() Pool{}

// Register adds a pool constructor to the registry.
func Register(name string, constructor func() Pool) {
	registry[name] = constructor
}

// Get returns a pool by name.
func Get(name string) (Pool, error) {
	ctor, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown pool: %s", name)
	}
	return ctor(), nil
}

// Names returns all registered pool names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}

// randomTree is a shared helper for building random trees.
func randomTree(p Pool, rng *rand.Rand, maxDepth int) expr.ExprNode {
	if maxDepth <= 1 {
		return p.RandomLeaf(rng)
	}
	// Bias toward leaves at shallow depths to keep trees small
	r := rng.Float64()
	switch {
	case r < 0.4:
		return p.RandomLeaf(rng)
	case r < 0.6:
		return &expr.UnaryNode{
			Op:    p.RandomUnary(rng),
			Child: randomTree(p, rng, maxDepth-1),
		}
	default:
		return &expr.BinaryNode{
			Op:    p.RandomBinary(rng),
			Left:  randomTree(p, rng, maxDepth-1),
			Right: randomTree(p, rng, maxDepth-1),
		}
	}
}
