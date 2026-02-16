package strategy

import (
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/expr"
	"github.com/wildfunctions/genetic_series/pkg/pool"
	"github.com/wildfunctions/genetic_series/pkg/series"
)

// MutationType identifies a kind of mutation.
type MutationType int

const (
	MutPoint         MutationType = iota // replace a random node with a new one
	MutSubtree                           // replace a random subtree with a new random tree
	MutHoist                             // replace tree with one of its subtrees
	MutConstPerturb                      // adjust a constant value by ±1-3
	MutGrow                              // wrap a leaf in a new operation
	MutShrink                            // replace a node with one of its children
)

const maxMutationDepth = 4

// MutateCandidate applies a random mutation to a candidate (modifies in place).
func MutateCandidate(c *series.Candidate, p pool.Pool, rng *rand.Rand) {
	r := rng.Float64()
	switch {
	case r < 0.1:
		// Flip start index (0 ↔ 1)
		c.Start = 1 - c.Start
	case r < 0.55:
		c.Numerator = mutateTree(c.Numerator, p, rng)
	default:
		c.Denominator = mutateTree(c.Denominator, p, rng)
	}
}

func mutateTree(root expr.ExprNode, p pool.Pool, rng *rand.Rand) expr.ExprNode {
	mut := MutationType(rng.Intn(6))
	switch mut {
	case MutPoint:
		return pointMutate(root, p, rng)
	case MutSubtree:
		return subtreeMutate(root, p, rng)
	case MutHoist:
		return hoistMutate(root, rng)
	case MutConstPerturb:
		return constPerturb(root, rng)
	case MutGrow:
		return growMutate(root, p, rng)
	case MutShrink:
		return shrinkMutate(root, rng)
	default:
		return root
	}
}

// pointMutate replaces a random node's operation (keeping children).
func pointMutate(root expr.ExprNode, p pool.Pool, rng *rand.Rand) expr.ExprNode {
	nodes := collectNodes(root)
	if len(nodes) == 0 {
		return root
	}
	idx := rng.Intn(len(nodes))
	target := nodes[idx]

	switch n := (*target).(type) {
	case *expr.VarNode:
		*target = p.RandomLeaf(rng)
	case *expr.ConstNode:
		*target = p.RandomLeaf(rng)
	case *expr.UnaryNode:
		n.Op = p.RandomUnary(rng)
	case *expr.BinaryNode:
		n.Op = p.RandomBinary(rng)
	}
	return root
}

// subtreeMutate replaces a random subtree with a new random tree.
func subtreeMutate(root expr.ExprNode, p pool.Pool, rng *rand.Rand) expr.ExprNode {
	nodes := collectNodes(root)
	if len(nodes) == 0 {
		return p.RandomTree(rng, maxMutationDepth)
	}
	idx := rng.Intn(len(nodes))
	*nodes[idx] = p.RandomTree(rng, maxMutationDepth)
	return root
}

// hoistMutate replaces the tree with one of its subtrees.
func hoistMutate(root expr.ExprNode, rng *rand.Rand) expr.ExprNode {
	nodes := collectNodes(root)
	if len(nodes) <= 1 {
		return root
	}
	idx := rng.Intn(len(nodes))
	return (*nodes[idx]).Clone()
}

// constPerturb adjusts a random constant by ±1 to ±3.
func constPerturb(root expr.ExprNode, rng *rand.Rand) expr.ExprNode {
	consts := collectConsts(root)
	if len(consts) == 0 {
		return root
	}
	target := consts[rng.Intn(len(consts))]
	delta := int64(rng.Intn(3) + 1)
	if rng.Float64() < 0.5 {
		delta = -delta
	}
	target.Val += delta
	if target.Val == 0 {
		target.Val = 1 // avoid zero constants
	}
	return root
}

// growMutate wraps a random leaf in a new unary or binary operation.
func growMutate(root expr.ExprNode, p pool.Pool, rng *rand.Rand) expr.ExprNode {
	nodes := collectNodes(root)
	if len(nodes) == 0 {
		return root
	}
	idx := rng.Intn(len(nodes))
	old := *nodes[idx]

	if rng.Float64() < 0.5 {
		*nodes[idx] = &expr.UnaryNode{Op: p.RandomUnary(rng), Child: old}
	} else {
		if rng.Float64() < 0.5 {
			*nodes[idx] = &expr.BinaryNode{Op: p.RandomBinary(rng), Left: old, Right: p.RandomLeaf(rng)}
		} else {
			*nodes[idx] = &expr.BinaryNode{Op: p.RandomBinary(rng), Left: p.RandomLeaf(rng), Right: old}
		}
	}
	return root
}

// shrinkMutate replaces a non-leaf node with one of its children.
func shrinkMutate(root expr.ExprNode, rng *rand.Rand) expr.ExprNode {
	nodes := collectNodes(root)
	if len(nodes) == 0 {
		return root
	}
	idx := rng.Intn(len(nodes))
	switch n := (*nodes[idx]).(type) {
	case *expr.UnaryNode:
		*nodes[idx] = n.Child
	case *expr.BinaryNode:
		if rng.Float64() < 0.5 {
			*nodes[idx] = n.Left
		} else {
			*nodes[idx] = n.Right
		}
	}
	return root
}

// collectNodes returns pointers to all nodes in the tree (for in-place mutation).
func collectNodes(root expr.ExprNode) []*expr.ExprNode {
	var result []*expr.ExprNode
	collectNodesHelper(&root, &result)
	return result
}

func collectNodesHelper(node *expr.ExprNode, result *[]*expr.ExprNode) {
	*result = append(*result, node)
	switch n := (*node).(type) {
	case *expr.UnaryNode:
		collectNodesHelper(&n.Child, result)
	case *expr.BinaryNode:
		collectNodesHelper(&n.Left, result)
		collectNodesHelper(&n.Right, result)
	}
}

// collectConsts returns pointers to all ConstNodes in the tree.
func collectConsts(root expr.ExprNode) []*expr.ConstNode {
	var result []*expr.ConstNode
	collectConstsHelper(root, &result)
	return result
}

func collectConstsHelper(node expr.ExprNode, result *[]*expr.ConstNode) {
	switch n := node.(type) {
	case *expr.ConstNode:
		*result = append(*result, n)
	case *expr.UnaryNode:
		collectConstsHelper(n.Child, result)
	case *expr.BinaryNode:
		collectConstsHelper(n.Left, result)
		collectConstsHelper(n.Right, result)
	}
}
