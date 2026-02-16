# Genetic Series Engine — v1 Implementation Plan

Status: **Implemented**

## Overview
Evolutionary algorithm tool that discovers mathematical series representations for known constants.

## Architecture
- `pkg/expr/` — Expression tree nodes with big.Float evaluation, printing, cloning, simplification, complexity metrics
- `pkg/constants/` — High-precision values for γ, π, e, ln(2), Catalan, Apéry
- `pkg/series/` — Candidate series, partial sum evaluation with convergence detection, multi-objective fitness
- `pkg/pool/` — Gene pool interface with conservative, moderate, kitchen-sink implementations
- `pkg/strategy/` — Strategy interface with hill-climbing and tournament selection
- `pkg/engine/` — Main evolutionary loop with parallel evaluation, config, output
- `main.go` — CLI entry point

## Usage
```bash
go run main.go -target e -generations 1000 -verbose
go run main.go -target pi -pool moderate -strategy tournament -population 500 -seed 42
go run main.go -target ln2 -format json
```
