package expr

import "math"

// Float64 lookup tables — fixed-size, computed at init, read-only.
var (
	factorialF64    [171]float64  // 170! is the last finite float64 factorial
	dblFactorialF64 [301]float64  // overflow around ~300
	fibonacciF64    [1477]float64 // fib(1476) is the last finite float64
)

func init() {
	// Factorials
	factorialF64[0] = 1
	for i := 1; i < len(factorialF64); i++ {
		factorialF64[i] = factorialF64[i-1] * float64(i)
	}

	// Double factorials: n!! = n * (n-2) * (n-4) * ... * (2 or 1)
	dblFactorialF64[0] = 1
	dblFactorialF64[1] = 1
	for i := 2; i < len(dblFactorialF64); i++ {
		dblFactorialF64[i] = dblFactorialF64[i-2] * float64(i)
	}

	// Fibonacci
	fibonacciF64[0] = 0
	fibonacciF64[1] = 1
	for i := 2; i < len(fibonacciF64); i++ {
		fibonacciF64[i] = fibonacciF64[i-1] + fibonacciF64[i-2]
	}
}

// EvalF64 for VarNode returns n.
func (v *VarNode) EvalF64(n float64) (float64, bool) {
	return n, true
}

// EvalF64 for ConstNode returns the constant value.
func (c *ConstNode) EvalF64(n float64) (float64, bool) {
	return float64(c.Val), true
}

// EvalF64 for UnaryNode dispatches on op.
func (u *UnaryNode) EvalF64(n float64) (float64, bool) {
	child, ok := u.Child.EvalF64(n)
	if !ok {
		return 0, false
	}

	switch u.Op {
	case OpNeg:
		return -child, true

	case OpFactorial:
		iv := int64(child)
		if child != float64(iv) || iv < 0 || iv >= int64(len(factorialF64)) {
			return 0, false
		}
		return factorialF64[iv], true

	case OpAltSign:
		iv := int64(child)
		if child != float64(iv) || iv < 0 {
			return 0, false
		}
		if iv%2 == 0 {
			return 1, true
		}
		return -1, true

	case OpDoubleFactorial:
		iv := int64(child)
		if child != float64(iv) || iv < 0 || iv >= int64(len(dblFactorialF64)) {
			return 0, false
		}
		return dblFactorialF64[iv], true

	case OpFibonacci:
		iv := int64(child)
		if child != float64(iv) || iv < 0 || iv >= int64(len(fibonacciF64)) {
			return 0, false
		}
		return fibonacciF64[iv], true

	case OpSin:
		if math.IsInf(child, 0) || math.IsNaN(child) {
			return 0, false
		}
		return math.Sin(child), true

	case OpCos:
		if math.IsInf(child, 0) || math.IsNaN(child) {
			return 0, false
		}
		return math.Cos(child), true

	case OpLn:
		if child <= 0 || math.IsInf(child, 0) || math.IsNaN(child) {
			return 0, false
		}
		return math.Log(child), true

	case OpFloor:
		if math.IsInf(child, 0) || math.IsNaN(child) {
			return 0, false
		}
		return math.Floor(child), true

	case OpCeil:
		if math.IsInf(child, 0) || math.IsNaN(child) {
			return 0, false
		}
		return math.Ceil(child), true

	case OpAbs:
		return math.Abs(child), true

	case OpSqrt:
		if child < 0 || math.IsNaN(child) {
			return 0, false
		}
		return math.Sqrt(child), true

	default:
		return 0, false
	}
}

// EvalF64 for BinaryNode dispatches on op.
func (b *BinaryNode) EvalF64(n float64) (float64, bool) {
	left, ok := b.Left.EvalF64(n)
	if !ok {
		return 0, false
	}
	right, ok := b.Right.EvalF64(n)
	if !ok {
		return 0, false
	}

	switch b.Op {
	case OpAdd:
		r := left + right
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return 0, false
		}
		return r, true

	case OpSub:
		r := left - right
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return 0, false
		}
		return r, true

	case OpMul:
		r := left * right
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return 0, false
		}
		return r, true

	case OpDiv:
		if right == 0 {
			return 0, false
		}
		r := left / right
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return 0, false
		}
		return r, true

	case OpPow:
		return powF64(left, right)

	case OpBinomial:
		return binomialF64(left, right)

	default:
		return 0, false
	}
}

// powF64 computes base^exp in float64.
func powF64(base, exp float64) (float64, bool) {
	// For integer exponents, use intPowF64 for precision
	ei := int64(exp)
	if exp == float64(ei) {
		if ei < 0 {
			if base == 0 {
				return 0, false
			}
			pos, ok := intPowF64(base, -ei)
			if !ok {
				return 0, false
			}
			r := 1.0 / pos
			if math.IsInf(r, 0) || math.IsNaN(r) {
				return 0, false
			}
			return r, true
		}
		return intPowF64(base, ei)
	}
	// Non-integer exponent
	if base < 0 {
		return 0, false
	}
	r := math.Pow(base, exp)
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, false
	}
	return r, true
}

// intPowF64 computes base^exp using binary exponentiation, exp >= 0, capped at 20.
func intPowF64(base float64, exp int64) (float64, bool) {
	if exp > 20 {
		return 0, false
	}
	result := 1.0
	b := base
	e := exp
	for e > 0 {
		if e%2 == 1 {
			result *= b
		}
		b *= b
		e /= 2
	}
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, false
	}
	return result, true
}

// binomialF64 computes C(n, k) in float64.
func binomialF64(nf, kf float64) (float64, bool) {
	ni := int64(nf)
	ki := int64(kf)
	if nf != float64(ni) || kf != float64(ki) {
		return 0, false
	}
	if ni < 0 || ki < 0 || ki > ni || ni > 1000 {
		return 0, false
	}
	if ki > ni-ki {
		ki = ni - ki
	}
	result := 1.0
	for i := int64(0); i < ki; i++ {
		result *= float64(ni-i) / float64(i+1)
		if math.IsInf(result, 0) || math.IsNaN(result) {
			return 0, false
		}
	}
	return result, true
}
