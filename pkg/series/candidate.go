package series

import (
	"fmt"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

// Candidate represents a candidate series: Sum_{n=Start}^{inf} Numerator(n) / Denominator(n)
type Candidate struct {
	Numerator   expr.ExprNode
	Denominator expr.ExprNode
	Start       int64 // starting index (0 or 1 typically)
}

// Clone returns a deep copy of the candidate.
func (c *Candidate) Clone() *Candidate {
	return &Candidate{
		Numerator:   c.Numerator.Clone(),
		Denominator: c.Denominator.Clone(),
		Start:       c.Start,
	}
}

// String returns a human-readable representation.
func (c *Candidate) String() string {
	return fmt.Sprintf("Sum_{n=%d}^{inf} (%s) / (%s)", c.Start, c.Numerator.String(), c.Denominator.String())
}

// LaTeX returns a LaTeX representation.
func (c *Candidate) LaTeX() string {
	return fmt.Sprintf("\\sum_{n=%d}^{\\infty} \\frac{%s}{%s}", c.Start, c.Numerator.LaTeX(), c.Denominator.LaTeX())
}

// Complexity returns combined complexity of both trees.
func (c *Candidate) Complexity() float64 {
	return expr.WeightedComplexity(c.Numerator) + expr.WeightedComplexity(c.Denominator)
}

// NodeCount returns the total node count of both trees.
func (c *Candidate) NodeCount() int {
	return c.Numerator.NodeCount() + c.Denominator.NodeCount()
}
