# Smarter Seeding: Initialize with Known Series Structures

## Problem
Random initialization produces mostly garbage. The GA wastes early generations just discovering that factorials in denominators make things converge. Every mathematician already knows this — we should tell the GA.

## Idea
Seed a fraction of the initial population with "templates" — structural skeletons of known series patterns with randomized constants. The rest stays fully random for exploration.

## Template Examples
- `1 / n!` — basic factorial series (e)
- `(-1)^n / (2n+1)` — alternating odd denominator (pi/4 Leibniz)
- `1 / n^k` — power series (zeta functions)
- `(-1)^n * k / n!` — alternating factorial with constant
- `1 / (n! * k^n)` — factorial-exponential hybrid
- `C(2n, n) / k^n` — central binomial coefficient series
- `1 / n!!` — double factorial series
- `(-1)^n / (k*n + c)` — generalized alternating harmonic

## Design Sketch
- Define a `Template` type: function that takes an RNG and returns a `*series.Candidate` with randomized constants
- Add a `templates.go` file in `pkg/strategy/` with 10-15 templates
- Modify `Initialize` in both strategies: seed 20-30% of population from templates, rest random
- Templates provide structure, mutation provides variation
- Could also inject templates when restarting after stagnation (fresh attempt gets template boost)

## Tradeoffs
- Biases search toward known patterns (could miss truly novel series)
- Templates need to be general enough that mutation can reshape them
- Risk of the population converging on template variations instead of exploring
- The 20-30% ratio needs tuning — too much kills diversity, too little has no effect

## Open Questions
- Should templates be pool-specific? (Conservative templates vs kitchensink templates)
- Should we inject templates on restart, on init, or both?
- Could we learn templates from the hall of fame? (Meta-evolution)
