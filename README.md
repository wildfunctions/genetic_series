# Genetic Series

The goal is to construct infinite series with genetic components. Under this construction, we can mutate the series under conditions that give rise to a selective pressure. We hope to discover series that closely approximate constants of our choosing, if not converge entirely.

Each candidate series has the form:

$$\sum_{n=s}^{\infty} \frac{\text{Numerator}(n)}{\text{Denominator}(n)}$$

where Numerator and Denominator are expression trees built from configurable building blocks — integers, factorials, powers, alternating signs, and more. The system evolves populations of these candidates, selecting for decimal digit accuracy against a target constant while penalizing complexity to favor elegant results.

## Quick Start

```bash
# Build
make build

# Search for pi (runs indefinitely, Ctrl+C to stop)
make run TARGET=pi

# Search for e with moderate pool and 5000 generation budget
make run TARGET=e POOL=moderate GENERATIONS=5000

# Run multiple targets in parallel (16-core machine)
GOMAXPROCS=5 make run TARGET=pi WORKERS=5 &
GOMAXPROCS=5 make run TARGET=e WORKERS=5 &
GOMAXPROCS=5 make run TARGET=euler_gamma WORKERS=5 &
```

## Available Targets

`pi`, `e`, `euler_gamma`, `ln2`, `catalan`, `apery`

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `-target` | `e` | Target constant |
| `-pool` | `conservative` | Gene pool: `conservative`, `moderate`, `kitchensink` |
| `-strategy` | `hillclimb` | Evolution strategy: `hillclimb`, `tournament` |
| `-population` | `200` | Population size |
| `-generations` | `1000` | Generation budget (0 = unlimited) |
| `-maxterms` | `1024` | Max terms to sum per series |
| `-stagnation` | `200` | Generations without improvement before restart |
| `-workers` | `NumCPU` | Parallel evaluation workers |
| `-seed` | `0` | Random seed (0 = random) |
| `-outdir` | `.` | Output directory for LaTeX/PDF |
| `-format` | `text` | Output format: `text`, `json` |
| `-verbose` | `false` | Per-generation output |

## Gene Pools

- **conservative** — `n`, integers 1-10, factorial, `(-1)^n`, negation, `+` `-` `*` `/`. Tight search space, most productive for common constants.
- **moderate** — Adds powers of 2/3, sqrt, exponentiation. Good middle ground.
- **kitchensink** — Adds double factorial, fibonacci, sin, cos, ln, floor, ceil. Large search space for exotic constants.

## How It Works

1. **Initialize** a random population of candidate series
2. **Evaluate** each candidate by summing terms and counting correct digits against the target
3. **Select** the fittest candidates (tournament selection or hill climbing)
4. **Evolve** via crossover and mutation (point, subtree, hoist, constant perturbation, grow, shrink)
5. **Repeat** until the generation budget is exhausted or the digit cap (50) is hit
6. **Restart** with a fresh population when stagnation is detected, preserving the best result in a hall of fame

The hall of fame is written to a LaTeX/PDF file after each restart attempt, so results survive long runs and Ctrl+C.

## Example Output

```
Starting target pi, pool conservative, strategy tournament, population 1000,
unlimited gen budget, stagnation 200, workers 16, seed 0

=== Attempt 1 ===
[gen 0] NEW BEST 1.2 digits | fitness 8.4321
  #1: Sum_{n=1}^{inf} (3) / ((n)^(5))
[gen 4] NEW BEST 3.1 digits | fitness 27.8900
  #1: Sum_{n=0}^{inf} ((-1)^(n)) / ((2 * n + 1))
[gen 47] Stagnated after 43 generations (3.1 digits, patience 20)

--- Hall of Fame ---
  #1: [attempt 1, gen 4]   3.1 digits | Sum_{n=0}^{inf} ((-1)^(n)) / ((2 * n + 1))

=== Attempt 2 ===
...
```

## Requirements

- Go 1.22+
- pdflatex (optional, for PDF generation)
