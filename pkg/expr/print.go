package expr

import "fmt"

var unaryOpNames = map[UnaryOp]string{
	OpNeg:              "-",
	OpFactorial:        "!",
	OpAltSign:          "(-1)^",
	OpDoubleFactorial:  "!!",
	OpFibonacci:        "fib",
	OpSin:              "sin",
	OpCos:              "cos",
	OpLn:               "ln",
	OpFloor:            "floor",
	OpCeil:             "ceil",
	OpAbs:              "abs",
	OpSqrt:             "sqrt",
}

var binaryOpSymbols = map[BinaryOp]string{
	OpAdd:      "+",
	OpSub:      "-",
	OpMul:      "*",
	OpDiv:      "/",
	OpPow:      "^",
	OpBinomial: "C",
}

// String methods

func (v *VarNode) String() string {
	return "n"
}

func (c *ConstNode) String() string {
	return fmt.Sprintf("%d", c.Val)
}

func (u *UnaryNode) String() string {
	child := u.Child.String()
	switch u.Op {
	case OpNeg:
		return fmt.Sprintf("(-%s)", child)
	case OpFactorial:
		return fmt.Sprintf("(%s)!", child)
	case OpAltSign:
		return fmt.Sprintf("(-1)^(%s)", child)
	case OpDoubleFactorial:
		return fmt.Sprintf("(%s)!!", child)
	default:
		name := unaryOpNames[u.Op]
		return fmt.Sprintf("%s(%s)", name, child)
	}
}

func (b *BinaryNode) String() string {
	left := b.Left.String()
	right := b.Right.String()
	sym := binaryOpSymbols[b.Op]
	switch b.Op {
	case OpBinomial:
		return fmt.Sprintf("C(%s, %s)", left, right)
	case OpPow:
		return fmt.Sprintf("(%s)^(%s)", left, right)
	default:
		return fmt.Sprintf("(%s %s %s)", left, sym, right)
	}
}

// LaTeX methods

func (v *VarNode) LaTeX() string {
	return "n"
}

func (c *ConstNode) LaTeX() string {
	return fmt.Sprintf("%d", c.Val)
}

func (u *UnaryNode) LaTeX() string {
	child := u.Child.LaTeX()
	switch u.Op {
	case OpNeg:
		return fmt.Sprintf("-{%s}", child)
	case OpFactorial:
		return fmt.Sprintf("{%s}!", child)
	case OpAltSign:
		return fmt.Sprintf("(-1)^{%s}", child)
	case OpDoubleFactorial:
		return fmt.Sprintf("{%s}!!", child)
	case OpFibonacci:
		return fmt.Sprintf("F_{%s}", child)
	case OpSin:
		return fmt.Sprintf("\\sin{(%s)}", child)
	case OpCos:
		return fmt.Sprintf("\\cos{(%s)}", child)
	case OpLn:
		return fmt.Sprintf("\\ln{(%s)}", child)
	case OpFloor:
		return fmt.Sprintf("\\lfloor %s \\rfloor", child)
	case OpCeil:
		return fmt.Sprintf("\\lceil %s \\rceil", child)
	case OpAbs:
		return fmt.Sprintf("|%s|", child)
	case OpSqrt:
		return fmt.Sprintf("\\sqrt{%s}", child)
	default:
		return child
	}
}

func (b *BinaryNode) LaTeX() string {
	left := b.Left.LaTeX()
	right := b.Right.LaTeX()
	switch b.Op {
	case OpAdd:
		return fmt.Sprintf("{%s} + {%s}", left, right)
	case OpSub:
		return fmt.Sprintf("{%s} - {%s}", left, right)
	case OpMul:
		return fmt.Sprintf("{%s} \\cdot {%s}", left, right)
	case OpDiv:
		return fmt.Sprintf("\\frac{%s}{%s}", left, right)
	case OpPow:
		return fmt.Sprintf("{%s}^{%s}", left, right)
	case OpBinomial:
		return fmt.Sprintf("\\binom{%s}{%s}", left, right)
	default:
		return ""
	}
}
