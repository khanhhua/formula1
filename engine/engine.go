package engine

import (
	"fmt"

	f1F "github.com/khanhhua/formula1/formula"

	"github.com/golang-collections/collections/stack"
	"github.com/tealeg/xlsx"
)

// Engine the F1 Engine
type Engine struct {
	xlFile *xlsx.File
	// Register AX
	ax interface{}
	// Current Code Index
	cci int
	// Max CCI
	maxCCI int
	// Execution stack
	callstack *stack.Stack
	// Execute and remember stuff here
	cache map[string]interface{}
}

type Invoke struct {
	fn 		string
	arity int
}

type Cell struct {
	value  string
	formula string
}

// NewEngine Create a new g to execute formula suitable for xlFile
func NewEngine(xlFile *xlsx.File) *Engine {
	return &Engine{
		xlFile: xlFile,
		callstack: stack.New(),
	}
}

func (g *Engine) GetCell (cellIDString string) (cell Cell,err error) {
	if len(cellIDString) == 0 {
		err = Error("Invalid address")
		return
	}

	if col, row, xlError := xlsx.GetCoordsFromCellIDString(cellIDString); xlError != nil {
		err = xlError
		return
	} else {
		xlCell := g.xlFile.Sheets[0].Cell(row, col)
		cell.value = xlCell.Value
		cell.formula = xlCell.Formula()
		return
	}
}

// Execute the g
func (g *Engine) Execute(inputs map[string]string) (names []string, outputs map[string]string) {
	// formula := NewFormula(formulaText)

	return
}

func (g *Engine) EvalFormula(f *f1F.Formula) interface{} {
	// The registers
	// var ax *interface{} = &g.ax
	// var cci *int = &g.cci
	// var maxCCI *int = &g.maxCCI
	var currentNode *f1F.Node
	var jump int

	currentNode = f.GetEntryNode()
	codes := []*f1F.Node{currentNode}
	if currentNode.NodeType() == f1F.NodeTypeFunc {
		if currentNode.Value().(string) == "IF" && currentNode.HasChildren() {
			codes = append(codes, currentNode.Children()...)
		}
	}
	g.maxCCI = len(codes) - 1
	fmt.Printf("Evaluating formula...\n")
	fmt.Printf("- Entry: %v\n", currentNode.Value())
	fmt.Printf("- MaxCCI: %d\n", g.maxCCI)

	// The mighty loop :D... LOOP LOOP LOOP until cci reaches maxCCI
	for g.cci <= g.maxCCI {
		// Phase 1 load node onto execution frame
		currentNode = codes[g.cci]
		nodeType := currentNode.NodeType()

		if nodeType == f1F.NodeTypeFunc {
			invoke := &Invoke{fn: currentNode.Value().(string)}
			jump = g.runStack(invoke)
		} else if nodeType == f1F.NodeTypeRef {
			invoke := &Invoke{fn: currentNode.Value().(string)}
			jump = g.runStack(invoke)
		} else if nodeType == f1F.NodeTypeLiteral ||
		 		nodeType == f1F.NodeTypeFloat ||
				nodeType == f1F.NodeTypeInteger {
			g.ax = currentNode.Value()
		} else if nodeType == f1F.NodeTypeOperator {
			g.callOperator(currentNode)
		}

		jump = 1
		g.cci += jump
	}

	if g.ax == nil {
		return "#NA"
	} else {
		switch g.ax.(type) {
		case string:
			return g.ax.(string)
		case int32:
			return g.ax.(int32)
		case float32:
			return g.ax.(float32)
		case float64:
			return g.ax.(float64)
		}
	}

	return "#NA"
}

// push Push whatever onto top of the g callstack
func (g *Engine) push(object interface{}) {
	g.callstack.Push(object)
}

// pop Pop top of the stack and store it in storage
func (g *Engine) pop(storage *interface{}) {
	*storage = g.callstack.Pop()
}

func (g *Engine) leave() {
	g.callstack.Pop()
}

// runStack Execute an invoke and store output in ax
// Output: ax register
func (g *Engine) runStack(invoke *Invoke) (jump int) {
	if invoke.fn == "+" {
		var operand interface{}
		var ax float64

		for i := 0; i < invoke.arity; i++ {
			g.pop(&operand)
			ax += operand.(float64)
		}
		g.ax = ax
	}

	jump = 1
	return
}

// run Expand an AST node and push result (ax) onto stack
func (g *Engine) callOperator(node *f1F.Node) {
	var invoke interface{}
	invoke = &Invoke{
		fn: node.Value().(string),
		arity: node.ChildCount(),
	}

	for _, childNode := range node.Children() {
		value := childNode.Value()
		g.push(value)
	}
	g.runStack(invoke.(*Invoke)) // Leave the stack
	g.push(g.ax)
}
