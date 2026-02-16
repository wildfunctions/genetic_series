# Guided Mutation: Bias Toward Convergent Structures

## Problem
Mutations are uniformly random. A point mutation is equally likely to replace a factorial with a negation as vice versa. But we know that certain structures are much more likely to produce convergent series — factorial denominators, alternating signs, decreasing terms. The GA rediscovers this every run through natural selection, but it's slow.

## Idea
Bias the mutation operators to favor structures that tend to produce convergent, elegant series. Not hard constraints — just weighted probabilities.

## Specific Biases
1. **Denominator-favoring mutations**: When mutating the denominator, bias toward:
   - Inserting factorial (`n!` makes almost anything converge)
   - Inserting powers (`n^k` for convergence of p-series)
   - Keeping variable `n` (denominator must depend on n)

2. **Numerator-favoring mutations**: When mutating the numerator, bias toward:
   - Alternating signs (`(-1)^n` for alternating series)
   - Small constants (elegant series use small integers)
   - Simple structures (fewer nodes)

3. **Constant perturbation bias**: When perturbing constants, prefer small values. Currently ±1 to ±3; could add a bias toward reducing magnitude.

4. **Complexity-aware mutation**: Higher chance of shrink/hoist mutations on complex trees, higher chance of grow/subtree on simple trees. Pushes toward a complexity sweet spot.

## Design Sketch
- Add weights to `mutateTree` selection (currently uniform `rng.Intn(6)`)
- Make weights context-dependent: different for numerator vs denominator
- Could make weights adaptive: track which mutations produce fitness improvements and increase their probability
- Add a "smart subtree" mutation that generates subtrees with known-good patterns instead of fully random

## Tradeoffs
- Reduces exploration of truly novel structures
- Harder to tune — many knobs to set
- Risk of the bias being wrong for some constants (maybe euler_gamma needs weird structures)
- Adaptive weights add complexity and could overfit to early results

## Open Questions
- How strong should the bias be? 2:1? 5:1? 10:1?
- Should bias strength decrease over time (anneal)?
- Should different targets get different biases?
- Could we use the hall of fame to learn biases? (If factorial appears in 80% of top results, increase factorial probability)
