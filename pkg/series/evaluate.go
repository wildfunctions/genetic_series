package series

import (
	"math"
	"math/big"
	"time"
)

// EvalResult holds the result of evaluating a candidate's partial sum.
type EvalResult struct {
	PartialSum      *big.Float
	TermsComputed   int64
	Converged       bool
	ConvergenceRate float64 // average ratio of |S_{2N} - S_N| decrease per doubling
	OK              bool
}

// evalTimeout is the maximum time allowed for evaluating a single candidate.
const evalTimeout = 100 * time.Millisecond

// EvaluateCandidate computes the partial sum of a candidate series up to maxTerms,
// using checkpoints at powers of 2 for convergence detection.
func EvaluateCandidate(c *Candidate, maxTerms int64, prec uint) EvalResult {
	sum := new(big.Float).SetPrec(prec)
	n := new(big.Float).SetPrec(prec)

	// Track partial sums at checkpoints (powers of 2)
	var checkpoints []checkpoint
	nextCheckpoint := int64(1)

	var termsComputed int64
	deadline := time.Now().Add(evalTimeout)

	for i := c.Start; i < c.Start+maxTerms; i++ {
		if time.Now().After(deadline) {
			return EvalResult{OK: false}
		}

		n.SetInt64(i)

		num, ok := c.Numerator.Eval(n, prec)
		if !ok {
			break // term failed — use partial sum so far
		}

		den, ok := c.Denominator.Eval(n, prec)
		if !ok {
			break
		}

		if den.Sign() == 0 {
			break
		}

		term := new(big.Float).SetPrec(prec).Quo(num, den)
		sum.Add(sum, term)
		termsComputed++

		// Record checkpoint at powers of 2 (relative to start)
		offset := i - c.Start + 1
		if offset == nextCheckpoint {
			checkpoints = append(checkpoints, checkpoint{
				terms: offset,
				sum:   new(big.Float).SetPrec(prec).Copy(sum),
			})
			nextCheckpoint *= 2
		}
	}

	// Need at least a few terms for a meaningful result
	if termsComputed < 4 {
		return EvalResult{OK: false}
	}

	// Compute convergence rate from checkpoints
	converged, rate := analyzeConvergence(checkpoints, prec)

	return EvalResult{
		PartialSum:      sum,
		TermsComputed:   termsComputed,
		Converged:       converged,
		ConvergenceRate: rate,
		OK:              true,
	}
}

type checkpoint struct {
	terms int64
	sum   *big.Float
}

// analyzeConvergence checks if |S_{2N} - S_N| is decreasing by a consistent factor.
func analyzeConvergence(cps []checkpoint, prec uint) (bool, float64) {
	if len(cps) < 3 {
		return false, 0
	}

	var diffs []float64
	for i := 1; i < len(cps); i++ {
		diff := new(big.Float).SetPrec(prec).Sub(cps[i].sum, cps[i-1].sum)
		diff.Abs(diff)
		f, _ := diff.Float64()
		diffs = append(diffs, f)
	}

	// Check that differences are decreasing
	if len(diffs) < 2 {
		return false, 0
	}

	var totalRatio float64
	var validRatios int
	converging := true

	for i := 1; i < len(diffs); i++ {
		if diffs[i-1] == 0 {
			// Perfect convergence at this point
			continue
		}
		ratio := diffs[i] / diffs[i-1]
		if ratio >= 1.0 {
			converging = false
		}
		totalRatio += ratio
		validRatios++
	}

	if validRatios == 0 {
		return true, 1.0 // converged exactly
	}

	avgRatio := totalRatio / float64(validRatios)
	return converging && avgRatio < 0.99, avgRatio
}

// EvalResultF64 holds the result of a float64 candidate evaluation.
type EvalResultF64 struct {
	PartialSum    float64
	TermsComputed int64
	Converged     bool
	OK            bool
}

// EvaluateCandidateF64 evaluates a candidate series entirely in float64.
// No timeout — float64 on 1024 terms runs in microseconds.
func EvaluateCandidateF64(c *Candidate, maxTerms int64) EvalResultF64 {
	var sum float64
	var termsComputed int64

	// Ring buffer of 3 checkpoint sums for convergence detection.
	var cpSums [3]float64
	cpIdx := 0
	cpCount := 0
	nextCheckpoint := int64(1)

	for i := c.Start; i < c.Start+maxTerms; i++ {
		n := float64(i)

		num, ok := c.Numerator.EvalF64(n)
		if !ok {
			break
		}

		den, ok := c.Denominator.EvalF64(n)
		if !ok {
			break
		}

		if den == 0 {
			break
		}

		term := num / den
		sum += term
		termsComputed++

		if math.IsInf(sum, 0) || math.IsNaN(sum) {
			return EvalResultF64{OK: false}
		}

		offset := i - c.Start + 1
		if offset == nextCheckpoint {
			cpSums[cpIdx%3] = sum
			cpIdx++
			cpCount++
			nextCheckpoint *= 2
		}
	}

	if termsComputed < 4 {
		return EvalResultF64{OK: false}
	}

	converged := analyzeConvergenceF64(cpSums[:], cpCount)

	return EvalResultF64{
		PartialSum:    sum,
		TermsComputed: termsComputed,
		Converged:     converged,
		OK:            true,
	}
}

// analyzeConvergenceF64 checks convergence from a ring buffer of checkpoint sums.
func analyzeConvergenceF64(ring []float64, count int) bool {
	if count < 3 {
		return false
	}
	// We have the last 3 checkpoints in the ring buffer.
	// Extract them in order.
	n := len(ring) // always 3
	oldest := (count - 3) % n
	s0 := ring[(oldest)%n]
	s1 := ring[(oldest+1)%n]
	s2 := ring[(oldest+2)%n]

	d0 := math.Abs(s1 - s0)
	d1 := math.Abs(s2 - s1)

	if d0 == 0 {
		return true
	}
	ratio := d1 / d0
	return ratio < 0.99
}
