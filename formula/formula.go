package formula

import (
	"fmt"
	"strconv"

	"github.com/xuri/efp"
)

// NodeType Descriptor of node in Formula1 AST
type NodeType int8

const (
	// NodeTypeRoot Root of AST
	NodeTypeRoot NodeType = iota + 1
	// NodeTypeLiteral Literal node, value should be understood as is
	NodeTypeLiteral
	// NodeTypeInteger Literal node, value should be understood as is
	NodeTypeInteger
	// NodeTypeFloat Literal node, value should be understood as is
	NodeTypeFloat
	// NodeTypeRef Reference node, value should be dereffed
	NodeTypeRef
	// NodeTypeFunc Function call node, value must be executed
	NodeTypeFunc
	// NodeTypeOperator Infix operator
	NodeTypeOperator
)

var PRECEDENCE = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

// Node AST node
type Node struct {
	value      interface{}
	nodeType   NodeType
	children   []*Node
	parent     *Node
	infixChild *Node
}

// Formula Formula1 executable formula
type Formula struct {
	root Node
}

// NewFormula Create a new formula instance
func NewFormula(text string) *Formula {
	parser := efp.ExcelParser()
	parser.Parse(text)
	fmt.Printf("%s\n=====\n", parser.PrettyPrint())

	tokens := parser.Tokens.Items
	root := Node{
		value:    "root",
		nodeType: NodeTypeRoot,
		children: nil,
	}

	current := &root
	index := 0
	count := len(tokens)

	var token *efp.Token

	for index <= count {
		if index >= count {
			if current.infixChild != nil {
				if current.children == nil {
					current.children = []*Node{current.infixChild}
				} else {
					current.children = append(current.children, current.infixChild)
				}
				current.resetInfixChild()
			}

			if root.infixChild != nil {
				if root.children == nil {
					root.children = []*Node{root.infixChild}
				} else {
					root.children = append(root.children, root.infixChild)
					root.resetInfixChild()
				}
			}
			break
		}

		token = &tokens[index]
		tvalue := token.TValue
		ttype := token.TType
		tsubtype := token.TSubType
		var value interface{}
		var nodeType NodeType

		if tsubtype == efp.TokenSubTypeStart {
			value, nodeType = resolveNodeType(ttype, tsubtype, tvalue)
			if current.infixChild != nil {
				parent := current
				current = current.infixChild.makeNode(nodeType, value)
				current.parent = parent
			} else {
				current = current.makeNode(nodeType, value)
			}

			index++
			continue
		} else if tsubtype == efp.TokenSubTypeStop {
			if current.infixChild != nil {
				current.children = append(current.children, current.infixChild)
				current.resetInfixChild()
			}

			current = current.parent

			index++
			continue
		} else if ttype == efp.TokenTypeArgument {
			if current.infixChild != nil {
				current.children = append(current.children, current.infixChild)
				current.resetInfixChild()
			}
			index++
			continue
		} else if ttype == efp.TokenTypeOperatorInfix {
			current.makeInfixChild(tvalue)

			index++
			continue
		} else if ttype == efp.TokenTypeOperand {
			value, nodeType = resolveNodeType(ttype, tsubtype, tvalue)
			if current.infixChild != nil {
				current.infixChild.makeNode(nodeType, value)
			} else {
				current.makeNode(nodeType, value)
			}

			index++
			continue
		} else {
			index++

			continue
		}
	}

	formula := Formula{
		root: root,
	}
	return &formula
}

func (formula *Formula) GetEntryNode() *Node {
	return formula.root.children[0]
}

func (parent *Node) makeNode(nodeType NodeType, value interface{}) *Node {
	node := Node{
		value:    value,
		nodeType: nodeType,
		parent:   parent,
		children: nil,
	}
	if parent.children == nil {
		parent.children = []*Node{&node}
	} else {
		parent.children = append(parent.children, &node)
	}

	return &node
}

func (parent *Node) makeInfixChild(value string) *Node {
	if parent.infixChild == nil {
		parent.infixChild = &Node{
			value:    value,
			nodeType: NodeTypeOperator,
			parent:   parent,
			children: []*Node{},
		}
		if len(parent.children) > 0 {
			lastChild := parent.LastChild()
			parent.children = parent.children[0 : parent.ChildCount()-1]
			parent.infixChild.children = []*Node{lastChild}
		}
	} else if parent.infixChild.value != value {
		// Infix precedence resolution
		if PRECEDENCE[value] > PRECEDENCE[parent.infixChild.value.(string)] {
			// Detach the last child and append it to the new node's children
			temp := parent.infixChild
			node := &Node{
				value:    value,
				nodeType: NodeTypeOperator,
				parent:   parent,
				children: []*Node{
					temp,
					temp.children[temp.ChildCount()-1],
				},
			}
			temp.children[temp.ChildCount()-1].parent = node
			temp.children = temp.children[:temp.ChildCount()-2]

			parent.infixChild = node
		} else {
			// Simply wrap the existing infixChild inside the new one
			temp := parent.infixChild
			node := &Node{
				value:    value,
				nodeType: NodeTypeOperator,
				parent:   parent,
				children: []*Node{
					temp,
				},
			}
			parent.infixChild = node
		}
	}

	return parent.infixChild
}

func (parent *Node) resetInfixChild() {
	parent.infixChild = nil
}

func resolveNodeType(ttype string, tsubtype string, tvalue string) (value interface{}, nodeType NodeType) {
	var _err error
	if ttype == efp.TokenTypeFunction && tsubtype == efp.TokenSubTypeStart {
		nodeType = NodeTypeFunc
		value = tvalue
		return
	} else if ttype == efp.TokenTypeSubexpression && tsubtype == efp.TokenSubTypeStart {
		nodeType = NodeTypeFunc
		value = "IDENTITY"
		return
	} else if ttype == efp.TokenTypeOperand && tsubtype == efp.TokenSubTypeRange {
		nodeType = NodeTypeRef
		value = tvalue
		return
	} else if ttype == efp.TokenTypeOperand && tsubtype == efp.TokenSubTypeText {
		nodeType = NodeTypeLiteral
		value = tvalue
		return
	} else if ttype == efp.TokenTypeOperand && tsubtype == efp.TokenSubTypeNumber {
		nodeType = NodeTypeFloat
		value, _err = strconv.ParseFloat(tvalue, 64)
		if _err != nil {
			return
		}
		return
	}

	value = tvalue
	nodeType = NodeTypeLiteral
	return
}

func (node *Node) Value() interface{} {
	return node.value
}

func (node *Node) NodeType() NodeType {
	return node.nodeType
}

func (parent *Node) ChildCount() int {
	if parent.children == nil {
		return 0
	}
	return len(parent.children)
}

func (parent *Node) FirstChild() *Node {
	if parent.children == nil || len(parent.children) == 0 {
		return nil
	}

	return parent.children[0]
}

func (parent *Node) LastChild() *Node {
	if parent.children == nil || len(parent.children) == 0 {
		return nil
	}

	return parent.children[len(parent.children)-1]
}

func (parent *Node) ChildAt(index int) *Node {
	if index < 0 {
		return nil
	} else if index >= len(parent.children) {
		return nil
	}

	return parent.children[index]
}

func (parent *Node) HasChildren() bool {
	if parent.children == nil {
		return false
	} else if len(parent.children) == 0 {
		return false
	}

	return true
}

func (parent *Node) Children() []*Node {
	if parent.children == nil {
		return make([]*Node, 0)
	}
	return parent.children[:]
}
