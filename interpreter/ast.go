package interpreter

type AST struct {
}

type Node interface {
}

type UnaryOp struct {
	Op   *Token
	Expr Node
}

type BinOp struct {
	Left  Node
	Op    *Token
	Right Node
}

type Num struct {
	Token *Token
	Value string
}
