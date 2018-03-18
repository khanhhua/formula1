package engine

import (
	"fmt"
	"strconv"
	"strings"

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
	fn    string
	arity int
}

type Cell struct {
	value   interface{}
	formula string
}

// NewEngine Create a new g to execute formula suitable for xlFile
func NewEngine(xlFile *xlsx.File) *Engine {
	return &Engine{
		xlFile:    xlFile,
		callstack: stack.New(),
	}
}

func (g *Engine) GetCell(cellIDString string) (cell Cell, err error) {
	if len(cellIDString) == 0 {
		err = Error("Invalid address")
		return
	}
	var sheet *xlsx.Sheet
	if strings.Contains(cellIDString, "!") {
		splat := strings.Split(cellIDString, "!")
		sheet = g.xlFile.Sheet[splat[0]]
		cellIDString = splat[1]
	} else {
		sheet = g.xlFile.Sheets[0]
	}

	if col, row, xlError := xlsx.GetCoordsFromCellIDString(cellIDString); xlError != nil {
		err = xlError
		return
	} else {
		xlCell := sheet.Cell(row, col)
		if f, err := strconv.ParseFloat(xlCell.Value, 64); err != nil {
			cell.value = xlCell.Value
		} else {
			cell.value = f
		}

		if formula := xlCell.Formula(); formula != "" {
			cell.formula = `=` + formula
		}

		return
	}
}

// Execute the g
func (g *Engine) Execute(inputs map[string]string) (names []string, outputs map[string]string) {
	// formula := NewFormula(formulaText)

	return
}

func (g *Engine) EvalFormula(f *f1F.Formula) (value interface{}, valueType f1F.NodeType) {
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
			g.runStack(invoke)
		} else if nodeType == f1F.NodeTypeRef {
			g.callDeref(currentNode)
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
		value = "#NA"
		valueType = 0
		return
	}

	value = g.ax

	switch g.ax.(type) {
	case string:
		valueType = f1F.NodeTypeLiteral
	case int32:
		valueType = f1F.NodeTypeInteger
	case float32:
		valueType = f1F.NodeTypeFloat
	case float64:
		valueType = f1F.NodeTypeFloat
	}
	return
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
func (g *Engine) runStack(invoke *Invoke) {
	if invoke.fn == "+" {
		var operand interface{}
		var ax float64

		for i := invoke.arity; i >= 1; i-- {
			g.pop(&operand)
			ax += operand.(float64)
		}
		g.ax = ax
	} else if invoke.fn == "-" {
		var operand interface{}
		var ax float64 = 0.0

		for i := invoke.arity; i >= 2; i-- {
			g.pop(&operand)
			ax += operand.(float64)
		}
		g.pop(&operand)
		g.ax = operand.(float64) - ax
	}
}

// run Expand an AST node and push result (ax) onto stack
func (g *Engine) callOperator(node *f1F.Node) {
	var invoke interface{}
	invoke = &Invoke{
		fn:    node.Value().(string),
		arity: node.ChildCount(),
	}

	for _, childNode := range node.Children() {
		switch childNode.NodeType() {
		case f1F.NodeTypeRef:
			g.callDeref(childNode)
		case f1F.NodeTypeFloat:
			value := childNode.Value()
			g.push(value)
			break
		case f1F.NodeTypeOperator:
			g.callOperator(childNode)
			break
		}
	}
	g.runStack(invoke.(*Invoke)) // Leave the stack
	g.push(g.ax)
}

func (g *Engine) callDeref(node *f1F.Node) {
	cellIDString := node.Value().(string)

	if cell, err := g.GetCell(cellIDString); err != nil {
		fmt.Printf("Could not deref %s. Reason: %v", cellIDString, err)
		return
	} else if cell.formula != "" {
		newEngine := NewEngine(g.xlFile)
		formula := f1F.NewFormula(cell.formula)
		result, _ := newEngine.EvalFormula(formula)
		g.ax = result
		g.push(g.ax)
	} else if cell.value != "" {
		g.ax = cell.value
		g.push(g.ax)
	}
}
