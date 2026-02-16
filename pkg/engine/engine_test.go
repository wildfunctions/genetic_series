package engine

import (
	"testing"

	_ "github.com/wildfunctions/genetic_series/pkg/pool"
	_ "github.com/wildfunctions/genetic_series/pkg/strategy"
)

func TestEngine_SmallRun(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "e"
	cfg.Population = 30
	cfg.Generations = 10
	cfg.MaxTerms = 128
	cfg.Seed = 42
	cfg.Verbose = false
	cfg.StagnationLimit = 5

	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	report := e.Run()

	if report.BestCandidate == "" {
		t.Error("Expected a best candidate")
	}
	if report.BestFitness.Combined <= -1e9 {
		t.Error("Expected non-worst fitness")
	}
	if len(report.Attempts) == 0 {
		t.Error("Expected at least one attempt in hall of fame")
	}

	t.Logf("Best after %d attempts: fitness=%.4f, digits=%.1f, candidate=%s",
		len(report.Attempts), report.BestFitness.Combined, report.BestFitness.CorrectDigits, report.BestCandidate)
}

func TestEngine_Restart(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "euler_gamma"
	cfg.Population = 10
	cfg.Generations = 50
	cfg.MaxTerms = 32
	cfg.Seed = 99
	cfg.StagnationLimit = 5

	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	report := e.Run()

	// With a hard target, tiny population, and short stagnation we expect restarts
	if len(report.Attempts) < 2 {
		t.Errorf("Expected multiple attempts with stagnation limit 5 and 50 gens, got %d", len(report.Attempts))
	}

	// Verify attempts are populated correctly
	for i, a := range report.Attempts {
		if a.Attempt != i+1 {
			t.Errorf("Attempt %d has wrong attempt number %d", i+1, a.Attempt)
		}
		if a.Generations == 0 {
			t.Errorf("Attempt %d has 0 generations", a.Attempt)
		}
		if a.BestCandidate == "" {
			t.Errorf("Attempt %d has empty best candidate", a.Attempt)
		}
		if a.BestFoundAtGen > a.Generations {
			t.Errorf("Attempt %d: BestFoundAtGen %d > Generations %d", a.Attempt, a.BestFoundAtGen, a.Generations)
		}
	}

	// Verify total gens used across attempts doesn't exceed budget
	totalGens := 0
	for _, a := range report.Attempts {
		totalGens += a.Generations
	}
	if totalGens > cfg.Generations {
		t.Errorf("Total generations %d exceeds budget %d", totalGens, cfg.Generations)
	}

	t.Logf("Completed %d attempts, total gens %d/%d, best=%.1f digits",
		len(report.Attempts), totalGens, cfg.Generations, report.BestFitness.CorrectDigits)
	for _, a := range report.Attempts {
		t.Logf("  Attempt %d: %d gens, best at gen %d, %.1f digits",
			a.Attempt, a.Generations, a.BestFoundAtGen, a.BestFitness.CorrectDigits)
	}
}

func TestEngine_Tournament(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "e"
	cfg.Strategy = "tournament"
	cfg.Population = 30
	cfg.Generations = 10
	cfg.MaxTerms = 128
	cfg.Seed = 42

	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	report := e.Run()

	if report.BestCandidate == "" {
		t.Error("Expected a best candidate")
	}
	t.Logf("Tournament best after 10 gens: fitness=%.4f, digits=%.1f",
		report.BestFitness.Combined, report.BestFitness.CorrectDigits)
}

func TestEngine_InvalidTarget(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "nonexistent"

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error for invalid target")
	}
}

func TestEngine_InvalidStrategy(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Strategy = "nonexistent"

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error for invalid strategy")
	}
}

func TestEngine_InvalidPool(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Pool = "nonexistent"

	_, err := New(cfg)
	if err == nil {
		t.Error("Expected error for invalid pool")
	}
}

// TestEngine_F64Disabled verifies that threshold=0 (no float64 fast path) still works.
func TestEngine_F64Disabled(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "e"
	cfg.Population = 30
	cfg.Generations = 5
	cfg.MaxTerms = 128
	cfg.Seed = 42
	cfg.StagnationLimit = 5
	cfg.F64PromotionThreshold = 0 // disabled

	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	report := e.Run()

	if report.BestCandidate == "" {
		t.Error("Expected a best candidate with F64 disabled")
	}
	if report.BestFitness.Combined <= -1e9 {
		t.Error("Expected non-worst fitness with F64 disabled")
	}
	if len(report.Attempts) == 0 {
		t.Error("Expected at least one attempt")
	}

	t.Logf("F64 disabled: fitness=%.4f, digits=%.1f, candidate=%s",
		report.BestFitness.Combined, report.BestFitness.CorrectDigits, report.BestCandidate)
}

func TestEngine_JSONFormat(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Target = "pi"
	cfg.Population = 20
	cfg.Generations = 3
	cfg.MaxTerms = 64
	cfg.Seed = 123
	cfg.Format = "json"

	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	report := e.Run()
	if report.BestCandidate == "" {
		t.Error("Expected a best candidate in JSON mode")
	}
}
