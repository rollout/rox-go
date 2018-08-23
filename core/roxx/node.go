package roxx

type NodeType int

const (
	NodeTypeRand NodeType = iota
	NodeTypeRator
	NodeTypeUnknown
)

type Node struct {
	Type  NodeType
	Value interface{}
}

func NewNode(nodeType NodeType, value interface{}) *Node {
	return &Node{Type: nodeType, Value: value}
}
