# Two-Phase Search: Float64 Sweep + Big.Float Verification

## Problem
Big.Float arithmetic is slow. Most candidates are garbage that gets rejected anyway. We're spending expensive precision on expressions that don't deserve it.

## Idea
Phase 1: Run the evolutionary search at float64 precision. Evaluation becomes native CPU floating point — no allocation, no big.Int, just raw speed. This lets us churn through massive populations and many attempts quickly. Any candidate that hits ~10+ digits at float64 is a genuine signal, not noise.

Phase 2: Take the top candidates from phase 1 and re-evaluate them at full big.Float precision (512-bit). Verify the digit count holds up. Candidates that were curve-fitting noise will collapse; real series will maintain or improve their digit count.

## Design Sketch
- Add a `float64` eval path alongside the existing `big.Float` path (or make precision switchable)
- Phase 1 engine runs with float64, large population (5000+), short attempts, many restarts
- Candidates that cross a digit threshold (e.g. 10) are saved to a "promotion queue"
- Phase 2 re-evaluates promoted candidates at full precision and ranks them
- Could run phase 1 continuously, feeding phase 2 as candidates appear

## Tradeoffs
- Float64 gives ~15 digits max — enough to identify real series but can't distinguish 15-digit matches from 50-digit ones
- Need to be careful: some series converge slowly and only look good at high precision (false negatives)
- Two codepaths to maintain (or a generic eval interface)
- GPU becomes viable for phase 1 if we ever go there

## Open Questions
- Should phase 1 use the same fitness function, just at lower precision?
- How many digits at float64 is a reliable signal? 8? 10? 12?
- Could we run phase 1 and phase 2 as separate processes communicating via files?
