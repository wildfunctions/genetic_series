package pool

import (
	"math/big"
	"math/rand"
	"testing"
)

const testPrec = 512

func TestConservativePool(t *testing.T) {
	p, err := Get("conservative")
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(42))

	// Generate many random trees and verify they evaluate cleanly
	successes := 0
	total := 1000
	for i := 0; i < total; i++ {
		tree := p.RandomTree(rng, 3)
		n := new(big.Float).SetPrec(testPrec).SetInt64(int64(rng.Intn(10) + 1))
		_, ok := tree.Eval(n, testPrec)
		if ok {
			successes++
		}
	}

	// At least 50% should evaluate successfully
	if float64(successes)/float64(total) < 0.5 {
		t.Errorf("Only %d/%d trees evaluated successfully", successes, total)
	}
	t.Logf("Conservative pool: %d/%d trees evaluated cleanly", successes, total)
}

func TestModeratePool(t *testing.T) {
	p, err := Get("moderate")
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(42))

	successes := 0
	total := 1000
	for i := 0; i < total; i++ {
		tree := p.RandomTree(rng, 3)
		n := new(big.Float).SetPrec(testPrec).SetInt64(int64(rng.Intn(10) + 1))
		_, ok := tree.Eval(n, testPrec)
		if ok {
			successes++
		}
	}

	if float64(successes)/float64(total) < 0.3 {
		t.Errorf("Only %d/%d trees evaluated successfully", successes, total)
	}
	t.Logf("Moderate pool: %d/%d trees evaluated cleanly", successes, total)
}

func TestKitchenSinkPool(t *testing.T) {
	p, err := Get("kitchensink")
	if err != nil {
		t.Fatal(err)
	}

	rng := rand.New(rand.NewSource(42))

	successes := 0
	total := 1000
	for i := 0; i < total; i++ {
		tree := p.RandomTree(rng, 3)
		n := new(big.Float).SetPrec(testPrec).SetInt64(int64(rng.Intn(10) + 1))
		_, ok := tree.Eval(n, testPrec)
		if ok {
			successes++
		}
	}

	if float64(successes)/float64(total) < 0.2 {
		t.Errorf("Only %d/%d trees evaluated successfully", successes, total)
	}
	t.Logf("Kitchen sink pool: %d/%d trees evaluated cleanly", successes, total)
}

func TestPoolRegistry(t *testing.T) {
	names := Names()
	if len(names) < 3 {
		t.Errorf("Expected at least 3 registered pools, got %d", len(names))
	}

	for _, name := range names {
		p, err := Get(name)
		if err != nil {
			t.Errorf("Get(%q) failed: %v", name, err)
			continue
		}
		if p.Name() != name {
			t.Errorf("Pool name mismatch: %q vs %q", p.Name(), name)
		}
	}
}

func TestUnknownPool(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Error("Expected error for unknown pool")
	}
}
