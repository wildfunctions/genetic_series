package strategy

import (
	"fmt"
	"math/rand"

	"github.com/wildfunctions/genetic_series/pkg/pool"
	"github.com/wildfunctions/genetic_series/pkg/series"
)

// Strategy defines an evolutionary strategy for evolving candidate series.
type Strategy interface {
	Name() string
	Initialize(p pool.Pool, rng *rand.Rand, popSize int) []*series.Candidate
	Evolve(population []*series.Candidate, fitnesses []series.Fitness, p pool.Pool, rng *rand.Rand) []*series.Candidate
}

var registry = map[string]func() Strategy{}

// Register adds a strategy constructor to the registry.
func Register(name string, constructor func() Strategy) {
	registry[name] = constructor
}

// Get returns a strategy by name.
func Get(name string) (Strategy, error) {
	ctor, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
	return ctor(), nil
}

// Names returns all registered strategy names.
func Names() []string {
	names := make([]string, 0, len(registry))
	for k := range registry {
		names = append(names, k)
	}
	return names
}

const (
	maxTreeDepth = 10 // reject trees deeper than this
	maxNodeCount = 25 // reject candidates with more total nodes than this
)

// candidateOK checks that a candidate isn't too deep or bloated.
func candidateOK(c *series.Candidate) bool {
	return c.Numerator.Depth() <= maxTreeDepth &&
		c.Denominator.Depth() <= maxTreeDepth &&
		c.NodeCount() <= maxNodeCount
}

// randomCandidate creates a random candidate with trees of given max depth.
func randomCandidate(p pool.Pool, rng *rand.Rand, maxDepth int) *series.Candidate {
	return &series.Candidate{
		Numerator:   p.RandomTree(rng, maxDepth),
		Denominator: p.RandomTree(rng, maxDepth),
		Start:       int64(rng.Intn(2)), // 0 or 1
	}
}
