package series

import (
	"math"
	"math/big"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

// FitnessWeights controls the relative importance of fitness components.
type FitnessWeights struct {
	Accuracy    float64
	Complexity  float64 // penalty weight (subtracted)
	Convergence float64
}

// DefaultWeights returns the default fitness weights.
func DefaultWeights() FitnessWeights {
	return FitnessWeights{
		Accuracy:    10.0,
		Complexity:  2.0,
		Convergence: 1.0,
	}
}

// Fitness holds the multi-objective fitness score for a candidate.
type Fitness struct {
	Combined        float64
	CorrectDigits   float64
	Simplicity      float64
	ConvergenceRate float64
}

// WorstFitness returns a fitness score for invalid/failed candidates.
func WorstFitness() Fitness {
	return Fitness{
		Combined:        -1e9,
		CorrectDigits:   0,
		Simplicity:      0,
		ConvergenceRate: 0,
	}
}

// ComputeFitness scores a candidate against a target value.
func ComputeFitness(c *Candidate, result EvalResult, target *big.Float, weights FitnessWeights) Fitness {
	if !result.OK {
		return WorstFitness()
	}

	// A series whose terms don't depend on n is just a constant times infinity — reject it.
	if !expr.ContainsVar(c.Numerator) && !expr.ContainsVar(c.Denominator) {
		return WorstFitness()
	}

	// Denominator must depend on n — otherwise terms don't shrink to zero and the series diverges.
	if !expr.ContainsVar(c.Denominator) {
		return WorstFitness()
	}

	// Reject non-convergent series — a partial sum that doesn't converge is meaningless.
	if !result.Converged {
		return WorstFitness()
	}

	// Reject divergent series — if partial sum is wildly off (>1e50 times target), it's garbage.
	if result.PartialSum != nil {
		absDiff := new(big.Float).Sub(result.PartialSum, target)
		absDiff.Abs(absDiff)
		absTgt := new(big.Float).Abs(target)
		if absTgt.Sign() > 0 {
			ratio, _ := new(big.Float).Quo(absDiff, absTgt).Float64()
			if math.IsInf(ratio, 0) || math.IsNaN(ratio) || ratio > 1e50 {
				return WorstFitness()
			}
		} else {
			f, _ := absDiff.Float64()
			if math.IsInf(f, 0) || math.IsNaN(f) || f > 1e50 {
				return WorstFitness()
			}
		}
	}

	correctDigits := countCorrectDigits(result.PartialSum, target)
	complexity := c.Complexity()
	simplicity := 1.0 / math.Max(complexity, 1.0)

	// Scale complexity penalty by accuracy: no penalty at 0 digits (allow exploration),
	// full penalty at 5+ digits (prevent bloat once candidates are accurate).
	penaltyScale := math.Min(correctDigits, 5.0) / 5.0

	combined := weights.Accuracy*correctDigits -
		weights.Complexity*complexity*penaltyScale

	return Fitness{
		Combined:        combined,
		CorrectDigits:   correctDigits,
		Simplicity:      simplicity,
		ConvergenceRate: result.ConvergenceRate,
	}
}

// MaxDigits is the cap on correct digits (limited by precision).
const MaxDigits = 50

// countCorrectDigits returns the number of matching decimal digits between two values.
func countCorrectDigits(computed, target *big.Float) float64 {
	if computed == nil || target == nil {
		return 0
	}

	diff := new(big.Float).Sub(computed, target)
	diff.Abs(diff)

	// If exact match
	if diff.Sign() == 0 {
		return MaxDigits // cap at 50 digits
	}

	// -log10(|computed - target| / |target|)
	absTgt := new(big.Float).Abs(target)
	if absTgt.Sign() == 0 {
		// Target is zero, use absolute error
		f, _ := diff.Float64()
		if f == 0 {
			return MaxDigits
		}
		return math.Max(0, -math.Log10(f))
	}

	relErr := new(big.Float).Quo(diff, absTgt)
	errFloat, _ := relErr.Float64()
	if errFloat == 0 {
		return MaxDigits
	}

	digits := -math.Log10(errFloat)
	return math.Max(0, digits)
}
