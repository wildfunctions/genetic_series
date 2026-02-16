package series

import (
	"math"
	"math/big"
	"testing"

	"github.com/wildfunctions/genetic_series/pkg/expr"
)

const testPrec = 512

func TestEvaluateCandidate_EMinusOne(t *testing.T) {
	// Sum_{n=0}^{inf} 1/n! = e
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.UnaryNode{Op: expr.OpFactorial, Child: &expr.VarNode{}},
		Start:       0,
	}

	result := EvaluateCandidate(c, 30, testPrec)
	if !result.OK {
		t.Fatal("EvaluateCandidate returned OK=false")
	}

	e, _ := new(big.Float).SetPrec(testPrec).SetString("2.71828182845904523536028747135266249775724709369995")
	diff := new(big.Float).Sub(result.PartialSum, e)
	diff.Abs(diff)

	eps := new(big.Float).SetPrec(testPrec).SetFloat64(1e-30)
	if diff.Cmp(eps) > 0 {
		t.Errorf("partial sum = %s, want ≈ e", result.PartialSum.Text('g', 40))
	}

	if !result.Converged {
		t.Log("Note: convergence not detected (expected for fast-converging series with few checkpoints)")
	}
}

func TestEvaluateCandidate_DivByZero(t *testing.T) {
	// 1/0 at n=0 should fail
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.ConstNode{Val: 0},
		Start:       0,
	}

	result := EvaluateCandidate(c, 10, testPrec)
	if result.OK {
		t.Error("Expected OK=false for 1/0")
	}
}

func TestFitness_KnownSeries(t *testing.T) {
	// 1/n! candidate targeting e
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.UnaryNode{Op: expr.OpFactorial, Child: &expr.VarNode{}},
		Start:       0,
	}

	result := EvaluateCandidate(c, 30, testPrec)
	if !result.OK {
		t.Fatal("evaluation failed")
	}

	e, _ := new(big.Float).SetPrec(testPrec).SetString("2.71828182845904523536028747135266249775724709369995")
	fitness := ComputeFitness(c, result, e, DefaultWeights())

	if fitness.CorrectDigits < 10 {
		t.Errorf("Expected many correct digits, got %.1f", fitness.CorrectDigits)
	}
	if fitness.Combined <= 0 {
		t.Errorf("Expected positive combined fitness, got %f", fitness.Combined)
	}

	t.Logf("1/n! fitness: combined=%.2f, digits=%.1f, simplicity=%.4f",
		fitness.Combined, fitness.CorrectDigits, fitness.Simplicity)
}

func TestFitness_BadCandidate(t *testing.T) {
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.ConstNode{Val: 0},
		Start:       0,
	}

	result := EvaluateCandidate(c, 10, testPrec)
	target := new(big.Float).SetPrec(testPrec).SetFloat64(2.718)
	fitness := ComputeFitness(c, result, target, DefaultWeights())

	if fitness.Combined != WorstFitness().Combined {
		t.Errorf("Expected worst fitness for bad candidate, got %f", fitness.Combined)
	}
}

func TestCandidateClone(t *testing.T) {
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.VarNode{},
		Start:       1,
	}

	clone := c.Clone()
	if clone.String() != c.String() {
		t.Errorf("Clone string mismatch: %q vs %q", clone.String(), c.String())
	}
	if clone.Start != c.Start {
		t.Errorf("Clone start mismatch: %d vs %d", clone.Start, c.Start)
	}
}

func TestCandidateComplexity(t *testing.T) {
	c := &Candidate{
		Numerator: &expr.BinaryNode{
			Op:    expr.OpAdd,
			Left:  &expr.VarNode{},
			Right: &expr.ConstNode{Val: 1},
		},
		Denominator: &expr.UnaryNode{
			Op:    expr.OpFactorial,
			Child: &expr.VarNode{},
		},
		Start: 0,
	}

	nc := c.NodeCount()
	if nc != 5 {
		t.Errorf("NodeCount() = %d, want 5", nc)
	}

	wc := c.Complexity()
	if wc <= 0 {
		t.Errorf("Complexity() = %f, want > 0", wc)
	}
}

// TestEvaluateCandidateF64_EMinusOne verifies 1/n! at float64 ≈ e with ~15 digits.
func TestEvaluateCandidateF64_EMinusOne(t *testing.T) {
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.UnaryNode{Op: expr.OpFactorial, Child: &expr.VarNode{}},
		Start:       0,
	}

	result := EvaluateCandidateF64(c, 30)
	if !result.OK {
		t.Fatal("EvaluateCandidateF64 returned OK=false")
	}

	e := 2.718281828459045
	diff := math.Abs(result.PartialSum - e)
	if diff > 1e-14 {
		t.Errorf("partial sum = %v, want ≈ e, diff = %v", result.PartialSum, diff)
	}
	if result.TermsComputed < 4 {
		t.Errorf("expected at least 4 terms, got %d", result.TermsComputed)
	}
	t.Logf("F64 partial sum of 1/n! = %.16g, terms=%d, converged=%v",
		result.PartialSum, result.TermsComputed, result.Converged)
}

// TestComputeFitnessF64_DegenerateRejection verifies constant series are still rejected.
func TestComputeFitnessF64_DegenerateRejection(t *testing.T) {
	// Both numerator and denominator are constants — should be rejected.
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.ConstNode{Val: 1},
		Start:       0,
	}

	result := EvaluateCandidateF64(c, 1024)
	fitness := ComputeFitnessF64(c, result, 2.718, DefaultWeights())

	if fitness.Combined != WorstFitness().Combined {
		t.Errorf("Expected worst fitness for constant series, got %f", fitness.Combined)
	}

	// Denominator doesn't depend on n — should be rejected.
	c2 := &Candidate{
		Numerator:   &expr.VarNode{},
		Denominator: &expr.ConstNode{Val: 2},
		Start:       0,
	}

	result2 := EvaluateCandidateF64(c2, 1024)
	fitness2 := ComputeFitnessF64(c2, result2, 2.718, DefaultWeights())

	if fitness2.Combined != WorstFitness().Combined {
		t.Errorf("Expected worst fitness for non-n denominator, got %f", fitness2.Combined)
	}
}

// TestComputeFitnessF64_KnownSeries verifies 1/n! gets good fitness against e.
func TestComputeFitnessF64_KnownSeries(t *testing.T) {
	c := &Candidate{
		Numerator:   &expr.ConstNode{Val: 1},
		Denominator: &expr.UnaryNode{Op: expr.OpFactorial, Child: &expr.VarNode{}},
		Start:       0,
	}

	result := EvaluateCandidateF64(c, 30)
	if !result.OK {
		t.Fatal("evaluation failed")
	}

	fitness := ComputeFitnessF64(c, result, 2.718281828459045, DefaultWeights())

	if fitness.CorrectDigits < 10 {
		t.Errorf("Expected many correct digits, got %.1f", fitness.CorrectDigits)
	}
	if fitness.Combined <= 0 {
		t.Errorf("Expected positive combined fitness, got %f", fitness.Combined)
	}

	t.Logf("F64 1/n! fitness: combined=%.2f, digits=%.1f", fitness.Combined, fitness.CorrectDigits)
}
