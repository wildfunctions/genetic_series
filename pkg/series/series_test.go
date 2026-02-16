package series

import (
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
