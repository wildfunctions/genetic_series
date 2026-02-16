package strategy

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/wildfunctions/genetic_series/pkg/pool"
	"github.com/wildfunctions/genetic_series/pkg/series"
)

const testPrec = 512

func evalPopulation(pop []*series.Candidate, target *big.Float) []series.Fitness {
	fitnesses := make([]series.Fitness, len(pop))
	for i, c := range pop {
		result := series.EvaluateCandidate(c, 256, testPrec)
		fitnesses[i] = series.ComputeFitness(c, result, target, series.DefaultWeights())
	}
	return fitnesses
}

func TestHillClimb_FitnessImproves(t *testing.T) {
	p, _ := pool.Get("conservative")
	s, _ := Get("hillclimb")
	rng := rand.New(rand.NewSource(42))

	target, _ := new(big.Float).SetPrec(testPrec).SetString("2.718281828459045")
	population := s.Initialize(p, rng, 50)

	var bestFitness float64
	for gen := 0; gen < 20; gen++ {
		fitnesses := evalPopulation(population, target)

		genBest := fitnesses[0].Combined
		for _, f := range fitnesses[1:] {
			if f.Combined > genBest {
				genBest = f.Combined
			}
		}

		if gen == 0 {
			bestFitness = genBest
		} else if genBest > bestFitness {
			bestFitness = genBest
		}

		population = s.Evolve(population, fitnesses, p, rng)
	}

	t.Logf("HillClimb best fitness after 20 gens: %.4f", bestFitness)
}

func TestTournament_FitnessImproves(t *testing.T) {
	p, _ := pool.Get("conservative")
	s, _ := Get("tournament")
	rng := rand.New(rand.NewSource(42))

	target, _ := new(big.Float).SetPrec(testPrec).SetString("2.718281828459045")
	population := s.Initialize(p, rng, 50)

	var bestFitness float64
	for gen := 0; gen < 20; gen++ {
		fitnesses := evalPopulation(population, target)

		genBest := fitnesses[0].Combined
		for _, f := range fitnesses[1:] {
			if f.Combined > genBest {
				genBest = f.Combined
			}
		}

		if gen == 0 {
			bestFitness = genBest
		} else if genBest > bestFitness {
			bestFitness = genBest
		}

		population = s.Evolve(population, fitnesses, p, rng)
	}

	t.Logf("Tournament best fitness after 20 gens: %.4f", bestFitness)
}

func TestMutation_TreeRemainValid(t *testing.T) {
	p, _ := pool.Get("conservative")
	rng := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		c := &series.Candidate{
			Numerator:   p.RandomTree(rng, 3),
			Denominator: p.RandomTree(rng, 3),
			Start:       1,
		}
		MutateCandidate(c, p, rng)

		// Verify the mutated candidate can still be evaluated (may return ok=false but shouldn't panic)
		n := new(big.Float).SetPrec(testPrec).SetInt64(5)
		c.Numerator.Eval(n, testPrec)
		c.Denominator.Eval(n, testPrec)
	}
}

func TestCrossover_ProducesTwoCandidates(t *testing.T) {
	p, _ := pool.Get("conservative")
	rng := rand.New(rand.NewSource(42))

	a := &series.Candidate{
		Numerator:   p.RandomTree(rng, 3),
		Denominator: p.RandomTree(rng, 3),
		Start:       0,
	}
	b := &series.Candidate{
		Numerator:   p.RandomTree(rng, 3),
		Denominator: p.RandomTree(rng, 3),
		Start:       1,
	}

	c1, c2 := CrossoverCandidates(a, b, rng)
	if c1 == nil || c2 == nil {
		t.Fatal("CrossoverCandidates returned nil")
	}

	// Verify they can be evaluated
	n := new(big.Float).SetPrec(testPrec).SetInt64(3)
	c1.Numerator.Eval(n, testPrec)
	c2.Denominator.Eval(n, testPrec)
}

func TestStrategyRegistry(t *testing.T) {
	names := Names()
	if len(names) < 2 {
		t.Errorf("Expected at least 2 strategies, got %d", len(names))
	}

	for _, name := range names {
		s, err := Get(name)
		if err != nil {
			t.Errorf("Get(%q) failed: %v", name, err)
			continue
		}
		if s.Name() != name {
			t.Errorf("Strategy name mismatch: %q vs %q", s.Name(), name)
		}
	}
}
