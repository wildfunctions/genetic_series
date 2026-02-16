# Genetic Series Engine — Project Context

**Status: v2 fully implemented and working.** Multi-attempt restart, hall of fame, LaTeX/PDF output, memoized evaluation, per-candidate timeout.

## What This Is

The goal of this project is to use evolution — the same blind process that designed eyes and brains — to discover new infinite series that sum to famous mathematical constants. Mathematicians have spent centuries hand-crafting series for pi, e, and other constants. We want to see if a computer, starting from random building blocks and guided only by "does this get closer to the target number?", can rediscover known series or stumble onto novel ones nobody has written down before. It's part math experiment, part AI exploration: can brute evolutionary pressure, applied to symbolic expressions rather than biological organisms, produce something genuinely beautiful?

In practice, it's a Go CLI tool that maintains a population of candidate series, each represented as a pair of expression trees (numerator and denominator). Every generation, candidates are mutated, crossed over, and selected based on how many correct digits their partial sums achieve against a target constant. The search is configurable — you pick the target constant, gene pool complexity, evolution strategy, and population size — then let it run and see what it finds.

Each candidate series has the form:

```
Sum_{n=start}^{inf} Numerator(n) / Denominator(n)
```

where Numerator and Denominator are expression trees built from configurable building blocks. The system evolves candidates via swappable mutation/selection strategies until a series converging to the target is found.

## Usage

```bash
# Build and run via Makefile
make run TARGET=pi STRATEGY=tournament POOL=conservative
make run TARGET=e POOL=moderate GENERATIONS=5000 WORKERS=5
make run TARGET=euler_gamma  # defaults: tournament, conservative, unlimited gens

# Direct binary usage
make build
./genetic_series -target pi -strategy tournament -pool conservative -population 1000 -generations 0 -maxterms 512 -seed 0 -workers 5

# Utilities
make test    # run all tests
make clean   # remove .tex, .pdf, .aux, .log files
```

CLI flags: `-target`, `-precision`, `-pool`, `-strategy`, `-population`, `-generations` (0=unlimited), `-maxterms`, `-seed` (0=random), `-format` (text/json), `-verbose`, `-maxdepth`, `-stagnation`, `-workers`, `-outdir`

## Project Structure

```
genetic_series/
├── main.go                        # CLI entry point (flag parsing, engine init, run, output)
├── Makefile                       # build, test, clean, run targets with configurable vars
├── .gitignore                     # ignores LaTeX/PDF output, build artifacts, IDE files
├── go.mod                         # module: github.com/wildfunctions/genetic_series, go 1.22.1
├── context.md                     # THIS FILE — session context for AI assistants
├── pkg/
│   ├── expr/                      # Expression tree system
│   │   ├── node.go                # ExprNode interface + VarNode, ConstNode, UnaryNode, BinaryNode
│   │   ├── eval.go                # big.Float evaluation, memoized factorial/fibonacci/double factorial
│   │   ├── print.go               # String() and LaTeX() rendering
│   │   ├── clone.go               # Deep copy
│   │   ├── complexity.go          # NodeCount, Depth, WeightedComplexity
│   │   ├── simplify.go            # Rewrite rules + constant folding (int and non-int)
│   │   └── expr_test.go           # Tests including known-series verification
│   ├── constants/
│   │   └── constants.go           # High-precision values (512-bit) for gamma, pi, e, ln2, catalan, apery
│   ├── series/
│   │   ├── candidate.go           # Candidate struct (two expr trees + start index)
│   │   ├── evaluate.go            # Partial sum with convergence detection, 2s timeout, graceful term failure
│   │   ├── fitness.go             # Accuracy + simplicity, degenerate series rejection
│   │   └── series_test.go
│   ├── pool/
│   │   ├── pool.go                # Pool interface + registry + shared randomTree helper
│   │   ├── conservative.go        # n, ints 1-10, factorial, (-1)^n, neg, +/-/*/÷
│   │   ├── moderate.go            # + powers of 2/3, sqrt, pow
│   │   ├── kitchensink.go         # + double factorial, fibonacci, sin, cos, ln, floor, ceil
│   │   └── pool_test.go
│   ├── strategy/
│   │   ├── strategy.go            # Strategy interface + registry + randomCandidate helper
│   │   ├── hillclimb.go           # Hill-climbing: clone+mutate, keep better, 5% random injection, elitism
│   │   ├── tournament.go          # Tournament: top 5% elite, tournament-select parents, crossover, 80% mutation
│   │   ├── mutation.go            # 7 mutation types: point, subtree, hoist, constPerturb, grow, shrink, start flip
│   │   ├── crossover.go           # Subtree crossover on both num/den trees
│   │   └── strategy_test.go
│   └── engine/
│       ├── engine.go              # Multi-attempt evolutionary loop with stagnation restart
│       ├── config.go              # Config struct + DefaultConfig()
│       ├── output.go              # Reports, hall of fame, LaTeX/PDF generation
│       └── engine_test.go
```

## Key Design Decisions

### Expression trees as interface
`ExprNode` interface with `VarNode`, `ConstNode`, `UnaryNode`, `BinaryNode` structs. Clean dispatch for `Eval`/`String`/`Clone`. New node kinds = new struct, no existing code modified.

### Eval returns `(*big.Float, bool)` not error
Performance in tight loops. Bool flag avoids GC pressure from error allocation across millions of evaluations.

### Precision
`math/big.Float` at 512 bits (~154 decimal digits) default. No CGo dependency. Configurable via CLI.

### Bad candidates handled at eval time
Division by zero, negative factorial, overflow all return `ok=false`. Fitness function assigns worst score (`-1e9`). Natural selection eliminates them. No need to constrain the search space.

### Convergence detection
At checkpoints N = 2^k, record S_N. Convergence = `|S_{2N} - S_N|` decreasing by consistent factor across doublings. Average ratio < 0.99 = converged.

### Parallelism
Candidate evaluation is embarrassingly parallel. Bounded worker pool (`-workers` flag, defaults to `runtime.NumCPU()`) with channels. Use `GOMAXPROCS` to truly pin OS threads when running multiple processes.

### Simplification
Runs after every mutation/crossover. Two-pass: algebraic rewrite rules (identity elimination, constant folding, double negation, etc.) then big.Float constant subtree evaluation. Non-integer constant subtrees (e.g. `1/(-13) + 9`) are rounded to nearest integer. Capped at 20 iterations.

TODO: support rational constants (e.g. `RatNode{Num, Den}`) so we can fold `1/3 + 1` to `4/3` instead of rounding.

### Fitness function
```
penaltyScale = min(CorrectDigits, 5) / 5
Combined = 10.0 * CorrectDigits - 2.0 * WeightedComplexity * penaltyScale
```
- CorrectDigits: `-log10(relative_error)`, capped at 50
- WeightedComplexity: sum of node weights for both trees. Large constants (`|val| > 10`) cost `1 + log10(|val|)` instead of 1.0. Complexity is **subtracted** as a penalty, scaled by accuracy.
- Convergence is NOT part of fitness — only accuracy and simplicity matter.
- **penaltyScale**: At 0 digits, complexity penalty is zero (free exploration). Ramps linearly to full penalty at 5+ digits (anti-bloat).

### Degenerate series rejection
ComputeFitness returns WorstFitness for:
- Failed evaluation (`!result.OK`)
- Constant series: neither numerator nor denominator contains variable `n`
- Non-shrinking terms: denominator doesn't contain `n` (terms won't approach zero)
- Non-convergent series: checkpoint analysis shows divergence
- Wildly divergent: partial sum >1e50 times target value

### Multi-attempt restart with stagnation detection
`-generations` is the **total budget** across all attempts (0=unlimited). When no improvement for `stagnation` generations, the attempt ends and a fresh population starts. Adaptive patience: `effectiveLimit = max(20, stagnationLimit * min(1.0, bestDigits / 10.0))`. Low-digit matches get short patience; high-digit matches get full patience. Early exit when 50-digit cap is hit.

### Hall of Fame
Best candidate from each attempt is saved. Sorted by CorrectDigits descending. Printed to stderr after each attempt. Written to LaTeX/PDF (if `-outdir` set) after each attempt so results survive Ctrl+C.

### LaTeX/PDF output
- File naming: `{target}_{pool}_{strategy}.tex/.pdf`
- Header shows all run parameters (target, pool, strategy, population, gen budget, stagnation, workers, seed)
- Each entry shows: rank, digits, attempt number, generation found, UTC timestamp, LaTeX formula, partial sum, error
- Compiled via pdflatex (if available) using /tmp staging to avoid WSL filesystem issues

### Per-candidate evaluation timeout
2-second deadline per candidate. Checked every 64 terms. Prevents pathological expressions (deeply nested factorial/fibonacci compositions) from blocking the entire generation.

### Graceful term failure
When a single term fails to evaluate (e.g. `factorial(21)` when cap is 20), evaluation stops and uses the partial sum computed so far instead of failing the entire candidate. Requires at least 4 successful terms.

### Memoized expensive operations
Factorial, double factorial, and fibonacci use thread-safe growing lookup tables (`sync.RWMutex`). Precomputed for inputs 0-20 at startup. On first access to a larger input, values are computed incrementally and cached. All subsequent accesses are a single slice lookup. Hard cap at input=1000.

### Hard complexity cap
Candidates with >25 total nodes or tree depth >10 are rejected and replaced with random candidates during evolution. This prevents runaway bloat from crossover/grow mutations.

### Strategy details
- **HillClimb**: Clone+mutate each candidate. Keep mutant (parent compared in next gen). Replace worst 5% with random. Best candidate preserved via elitism.
- **Tournament**: Top 5% elite carried forward. Rest: tournament-select 2 parents (size=5), subtree crossover on both trees, 80% chance of mutation, simplify, reject trees deeper than 10. Replace rejected with random.

### Mutation types (7)
1. **Start flip** (10%): Toggle start index between 0 and 1
2. **Numerator mutation** (45%): One of the 6 tree mutations below
3. **Denominator mutation** (45%): One of the 6 tree mutations below

Tree mutations (equal probability):
1. **Point**: Replace a random node's operation (keep children)
2. **Subtree**: Replace a random subtree with a new random tree
3. **Hoist**: Replace tree with one of its subtrees
4. **ConstPerturb**: Adjust a constant by ±1 to ±3
5. **Grow**: Wrap a node in a new unary or binary operation
6. **Shrink**: Replace a non-leaf node with one of its children

### Expression operations supported
- **Unary**: Neg, Factorial, AltSign `(-1)^n`, DoubleFactorial, Fibonacci, Sqrt, Sin, Cos, Ln, Floor, Ceil, Abs
- **Binary**: Add, Sub, Mul, Div, Pow, Binomial `C(n,k)`
- Sin/Cos/Ln/Sqrt fall back to float64 (not arbitrary precision)
- Factorial/DoubleFactorial/Fibonacci memoized, hard cap at input=1000
- IntPow uses binary exponentiation, capped at exp=200
- Large constants (`|val| > 10`) have higher complexity weight: `1 + log10(|val|)`
- Sqrt of perfect square constants folds during simplification (e.g. `sqrt(9)` → `3`)

### Pool configurations
- **conservative**: n, ints 1-10, factorial, (-1)^n, neg, +/-/*/÷. Tight search space, most productive for common constants.
- **moderate**: Adds powers of 2/3 as leaves, sqrt as unary, power as binary. Good middle ground.
- **kitchensink**: Adds double factorial, fibonacci, sin, cos, ln, floor, ceil. Large search space — most random candidates are garbage. Better for constants that need exotic operations. Can be slow due to expensive evaluations (mitigated by timeout).

### Constants available (all 512-bit precision)
`euler_gamma`, `pi`, `e`, `ln2`, `catalan`, `apery`

### Output
- **Startup**: Prints all run parameters (target, pool, strategy, population, gen budget, stagnation, workers, seed)
- **Non-verbose (default)**: Prints to stderr on new best + every 20 gens. Shows top 2 candidates with digit counts.
- **Verbose**: Full generation report every gen.
- **Hall of fame**: Printed after each attempt, sorted by digit accuracy descending.
- **LaTeX/PDF**: Written after each attempt when `-outdir` is set.
- **Final report**: Printed to stdout in text or JSON format.

### Gene pool tree generation
Shared `randomTree` function with depth bias: 40% leaf, 20% unary, 40% binary at each level. Forces leaf at maxDepth=1.

### Pool/Strategy registries
Both use `init()` + map-based registry pattern. `Get(name)` returns instance, `Names()` lists all registered.

## Test Coverage

All packages have tests:
- `expr`: VarNode, ConstNode, factorial, altSign, binary ops, div-by-zero, pow, binomial, fibonacci, double factorial, clone, complexity, string, LaTeX, simplify rules, floor/ceil, **known series verification** (1/n! = e, Leibniz = pi/4)
- `series`: EvaluateCandidate for 1/n!, div-by-zero handling, fitness scoring, candidate clone/complexity
- `pool`: All 3 pools generate 1000 trees and check eval success rates (conservative >50%, moderate >30%, kitchensink >20%), registry, unknown pool error
- `strategy`: HillClimb and Tournament fitness over 20 gens, mutation validity (100 rounds, no panics), crossover produces 2 candidates, registry
- `engine`: Small runs (10 gens) with hillclimb and tournament, restart behavior with stagnation, invalid target/strategy/pool errors, JSON format

## No External Dependencies

Pure Go stdlib. No third-party packages.

## Change Log / Decision History

### v2.0 — Multi-attempt restart + hall of fame (Feb 2026)

**Problem observed:** The engine gets stuck in local optima. A bad random seed produces a population that converges to garbage for 1400+ generations, wasting the entire budget.

**Changes:**
- `-generations` is now total budget across all attempts (0=unlimited)
- `-stagnation` flag (default 200) controls when to restart
- Adaptive stagnation patience: scales with digit quality (`max(20, limit * min(1.0, digits/10.0))`)
- Hall of fame tracks best from each attempt, sorted by digits
- LaTeX/PDF output after each attempt (survives Ctrl+C)
- File naming: `{target}_{pool}_{strategy}.tex/.pdf`
- Header in PDF shows all run parameters + UTC timestamps per entry
- Per-candidate 2-second evaluation timeout
- Graceful term failure (break instead of hard-fail)

### v1.3 — Degenerate series rejection (Feb 2026)

**Problem observed:** Constant series like `Sum (-9)/(-7983)` getting 4.7 digits, and divergent series with partial sum 10^153 topping the hall of fame.

**Changes:**
- Reject series where neither numerator nor denominator contains `n`
- Reject series where denominator doesn't contain `n` (terms don't shrink to zero)
- Reject non-convergent series
- Reject divergent series (partial sum >1e50 times target)

### v1.2.1 — Sqrt operator + memoized evaluation (Feb 2026)

**Changes:**
- Added `OpSqrt` unary operator (evaluates via `math.Sqrt`, folds perfect squares during simplification)
- Available in moderate and kitchensink pools
- Non-integer constant subtrees rounded to nearest integer during folding
- Factorial/double factorial/fibonacci use thread-safe memoized lookup tables
- Start index mutation (10% chance to flip 0↔1)
- `-workers` flag for controlling parallelism across multiple processes

### v1.2 — Accuracy-scaled complexity penalty (Feb 2026)

**Problem observed:** The v1.1 complexity penalty completely killed exploration. With a flat `-2*Complexity` penalty, every candidate with 0 correct digits was ranked purely by simplicity. The system would actively strip away factorials and powers because adding them made the fitness score worse.

**Fix:** Scale the complexity penalty by accuracy level:
```
penaltyScale = min(CorrectDigits, 5) / 5
penalty = Complexity * penaltyScale * weight
```
- At 0 digits: penalty = 0 (explore freely)
- At 5+ digits: penalty = 100% (full anti-bloat)

### v1.1 — Anti-bloat overhaul (Feb 2026)

**Problem observed:** The system was producing bloated candidates stuffed with magic constants. Two failure modes: "new best" with same digits but trivially better convergence on more complex expressions, and runaway expression bloat.

**Changes:**
1. Large constants cost more in complexity scoring (`1.0 + log10(|val|)` for `|val| > 10`)
2. Fitness changed from additive simplicity bonus to subtractive complexity penalty
3. Hard complexity cap: max 25 nodes, max depth 10
4. Convergence removed from fitness — only accuracy and simplicity matter
