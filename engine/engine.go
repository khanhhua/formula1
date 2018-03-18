package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	f1F "github.com/khanhhua/formula1/formula"
	funs "github.com/khanhhua/formula1/funs"

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

// EvalFormula Execute formula by running AST nodes as necessary
func (g *Engine) EvalFormula(f *f1F.Formula) (value interface{}, valueType f1F.NodeType) {
	var currentNode *f1F.Node

	currentNode = f.GetEntryNode()
	fmt.Printf("Evaluating formula...\n")
	fmt.Printf("- Entry: %v\n", currentNode.Value())

	stackHeight := g.callstack.Len()
	g.evalNode(currentNode)
	if g.callstack.Len() != stackHeight {
		panic(errors.New(fmt.Sprintf("Stack not disposed properly: was %d, now %d",
			stackHeight, g.callstack.Len())))
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
	var ret interface{}

	if invoke.fn == "+" {
		var operand interface{}
		var ax float64

		for i := invoke.arity; i >= 1; i-- {
			g.pop(&operand)
			ax += operand.(float64)
		}
		// g.ax = ax
		ret = ax
	} else if invoke.fn == "-" {
		var operand interface{}
		var ax float64 = 0.0

		for i := invoke.arity; i >= 2; i-- {
			g.pop(&operand)
			ax += operand.(float64)
		}
		g.pop(&operand)
		// g.ax = operand.(float64) - ax
		ret = operand.(float64) - ax
	} else {
		// Non primitive operators: + - * /
		// NOTE: DO NOT REFACTOR INTO DYNAMIC METHOD CALLING WITH ARGS...
		if invoke.arity == 1 {
			var operand1 interface{}
			g.pop(&operand1)
			if output, err := funs.Call1(invoke.fn, operand1); err != nil {
				// g.ax = err
				ret = err
			} else {
				// g.ax = output
				ret = output
			}
		} else if invoke.arity == 2 {
			var operand1, operand2 interface{}
			g.pop(&operand2)
			g.pop(&operand1)
			if output, err := funs.Call2(invoke.fn, operand1, operand2); err != nil {
				// g.ax = err
				ret = err
			} else {
				// g.ax = output
				ret = output
			}
		}
	}
	// NOTE: Remember to g.pop after g.runStack
	g.push(ret)
}

func (g *Engine) evalNode(node *f1F.Node) {
	switch node.NodeType() {
	case f1F.NodeTypeOperator:
		g.callFunc(node)
		break
	case f1F.NodeTypeFunc:
		g.callFunc(node)
		break
	case f1F.NodeTypeRef:
		g.callDeref(node)
		break
	case f1F.NodeTypeLiteral, f1F.NodeTypeFloat, f1F.NodeTypeInteger:
		g.ax = node.Value()
		break
	}
}

// run Expand an AST node and push result (ax) onto stack
func (g *Engine) callFunc(node *f1F.Node) {
	var fn string
	fn = node.Value().(string)

	if fn == "IF" { // The IF-JUMP
		if g.callIf(node.FirstChild()) {
			g.evalNode(node.ChildAt(1))
		} else if falseBranch := node.ChildAt(2); falseBranch != nil {
			g.evalNode(falseBranch)
		} else {
			g.ax = false
			g.push(g.ax)
		}
	}

	var invoke interface{}
	invoke = &Invoke{
		fn:    fn,
		arity: node.ChildCount(),
	}

	for _, childNode := range node.Children() {
		switch childNode.NodeType() {
		case f1F.NodeTypeRef:
			g.callDeref(childNode)
			g.push(g.ax)
		case f1F.NodeTypeFloat:
			value := childNode.Value()
			g.push(value)
			break
		case f1F.NodeTypeOperator:
			g.callFunc(childNode)
			g.push(g.ax)
			break
		case f1F.NodeTypeFunc:
			g.callFunc(childNode)
			g.push(g.ax)
			break
		}
	}
	g.runStack(invoke.(*Invoke))
	g.pop(&g.ax) // Leave the stack
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
		// g.push(g.ax)
	} else if cell.value != "" {
		g.ax = cell.value
		// g.push(g.ax)
	}
}

func (g *Engine) callIf(node *f1F.Node) bool {
	g.evalNode(node)
	if g.ax == nil || g.ax == false {
		return false
	}
	return true
}
