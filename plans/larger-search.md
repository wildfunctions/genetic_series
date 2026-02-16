# Larger Search: Scale Up Population and Attempts

## Problem
With population 1000 and conservative pool, we explore a tiny fraction of the search space per attempt. Many good series may be just a few mutations away from a candidate that got discarded.

## Idea
Simply throw more compute at the problem. Larger populations, more attempts, longer runs. No algorithmic changes — just scale.

## Configurations to Try

### Brute force sweep
```bash
# Run 3 processes in parallel, each targeting a different constant
GOMAXPROCS=5 make run TARGET=pi POOL=conservative WORKERS=5 &
GOMAXPROCS=5 make run TARGET=e POOL=conservative WORKERS=5 &
GOMAXPROCS=5 make run TARGET=euler_gamma POOL=conservative WORKERS=5 &
```

### Large population, short attempts
```bash
# 5000 population, aggressive stagnation = many diverse attempts
make run TARGET=pi POOL=conservative GENERATIONS=0 WORKERS=16 \
  -population 5000 -stagnation 50
```

### Small population, many attempts
```bash
# 200 population, very short stagnation = rapid restarts, lottery-style
make run TARGET=pi POOL=conservative GENERATIONS=0 WORKERS=16 \
  -population 200 -stagnation 20
```

### Multi-pool sweep
```bash
# Try all pools simultaneously
GOMAXPROCS=5 make run TARGET=pi POOL=conservative WORKERS=5 &
GOMAXPROCS=5 make run TARGET=pi POOL=moderate WORKERS=5 &
GOMAXPROCS=5 make run TARGET=pi POOL=kitchensink WORKERS=5 &
```

## What to Measure
- Digits achieved vs wall-clock time
- Digits achieved vs total generations consumed
- Number of attempts before first 5+ digit match
- Whether larger population or more restarts is more effective

## Tradeoffs
- Pure compute cost — electricity and time
- Diminishing returns: 10x compute probably doesn't give 10x digits
- Doesn't address fundamental search space limitations
- But it's the easiest thing to try and requires zero code changes

## Open Questions
- What's the optimal population size for conservative pool? (Maybe 500 is better than 5000)
- What's the optimal stagnation limit? (Short = more diversity, long = more exploitation)
- Is there a point where scaling population has zero marginal return?
- Should we run the same config with different seeds and compare, to understand variance?
