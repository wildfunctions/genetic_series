package expr

func (v *VarNode) Clone() ExprNode {
	return &VarNode{}
}

func (c *ConstNode) Clone() ExprNode {
	return &ConstNode{Val: c.Val}
}

func (u *UnaryNode) Clone() ExprNode {
	return &UnaryNode{
		Op:    u.Op,
		Child: u.Child.Clone(),
	}
}

func (b *BinaryNode) Clone() ExprNode {
	return &BinaryNode{
		Op:    b.Op,
		Left:  b.Left.Clone(),
		Right: b.Right.Clone(),
	}
}
