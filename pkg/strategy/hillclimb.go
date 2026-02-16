package strategy

import (
	"math/rand"
	"sort"

	"github.com/wildfunctions/genetic_series/pkg/expr"
	"github.com/wildfunctions/genetic_series/pkg/pool"
	"github.com/wildfunctions/genetic_series/pkg/series"
)

const (
	hillclimbMaxDepth      = 4
	hillclimbInjectionRate = 0.05 // fraction of population replaced with random each gen
)

func init() {
	Register("hillclimb", func() Strategy { return &HillClimbStrategy{} })
}

// HillClimbStrategy implements directed hill-climbing with population.
// For each candidate: clone + directed mutation, keep whichever is better.
// Periodically injects random candidates to escape local optima.
type HillClimbStrategy struct{}

func (s *HillClimbStrategy) Name() string { return "hillclimb" }

func (s *HillClimbStrategy) Initialize(p pool.Pool, rng *rand.Rand, popSize int) []*series.Candidate {
	pop := make([]*series.Candidate, popSize)
	for i := range pop {
		pop[i] = randomCandidate(p, rng, hillclimbMaxDepth)
	}
	return pop
}

func (s *HillClimbStrategy) Evolve(
	population []*series.Candidate,
	fitnesses []series.Fitness,
	p pool.Pool,
	rng *rand.Rand,
) []*series.Candidate {
	n := len(population)
	next := make([]*series.Candidate, n)

	for i := 0; i < n; i++ {
		// Clone and mutate
		child := population[i].Clone()
		MutateCandidate(child, p, rng)
		child.Numerator = expr.SimplifyBigFloat(child.Numerator, 128)
		child.Denominator = expr.SimplifyBigFloat(child.Denominator, 128)

		if !candidateOK(child) {
			child = randomCandidate(p, rng, hillclimbMaxDepth)
		}

		next[i] = child
	}

	// Sort by fitness to identify worst candidates for injection
	type indexed struct {
		idx     int
		fitness float64
	}
	ranked := make([]indexed, n)
	for i := range ranked {
		ranked[i] = indexed{i, fitnesses[i].Combined}
	}
	sort.Slice(ranked, func(a, b int) bool {
		return ranked[a].fitness < ranked[b].fitness
	})

	// Replace worst candidates with random injection
	injectionCount := int(float64(n) * hillclimbInjectionRate)
	if injectionCount < 1 {
		injectionCount = 1
	}
	for i := 0; i < injectionCount && i < n; i++ {
		idx := ranked[i].idx
		next[idx] = randomCandidate(p, rng, hillclimbMaxDepth)
	}

	// Elitism: keep the best from the old generation if it's better
	bestIdx := ranked[len(ranked)-1].idx
	next[bestIdx] = population[bestIdx].Clone()

	return next
}
