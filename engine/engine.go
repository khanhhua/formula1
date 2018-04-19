package engine

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	f1F "github.com/khanhhua/formula1/formula"
	funs "github.com/khanhhua/formula1/funs"

	"github.com/golang-collections/collections/stack"
	"github.com/tealeg/xlsx"
)

// Engine the F1 Engine
type Engine struct {
	xlFile      *xlsx.File
	activeSheet *xlsx.Sheet
	// Register AX
	ax interface{}
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

type Range struct {
	cells []Cell
	// width Number of columns
	colCount int
	// height Number of rows
	rowCount int
}

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "ENGINE ", log.Lshortfile)
}

// ToSlice Get a copy of the 1D Range
func (cellRange *Range) ToSlice() (cells []Cell, ok bool) {
	if cellRange.colCount == 1 && cellRange.rowCount > 1 {
		cells = cellRange.cells[:]
		ok = true
		return
	} else if cellRange.colCount > 1 && cellRange.rowCount == 1 {
		cells = cellRange.cells[:]
		ok = true
		return
	}

	ok = false
	return
}

func (cellRange *Range) To2DSlice() (cells [][]Cell, ok bool) {
	cells = make([][]Cell, cellRange.rowCount)
	colCount := cellRange.colCount

	for i := range cells {
		cells[i] = make([]Cell, colCount)

		for j := range cells[i] {
			cells[i][j] = cellRange.cells[i*colCount+j]
		}
	}
	ok = true
	return
}

// NewEngine Create a new g to execute formula suitable for xlFile
func NewEngine(xlFile *xlsx.File) *Engine {
	return &Engine{
		xlFile:    xlFile,
		callstack: stack.New(),
	}
}

func (g *Engine) GetRange(rangeIDString string) (cellRange Range, err error) {
	if len(rangeIDString) == 0 {
		err = errors.New("Invalid address")
		return
	}

	var sheet *xlsx.Sheet
	var fromIDString, toIDString string

	if strings.Contains(rangeIDString, "!") {
		splat := strings.Split(rangeIDString, "!")
		sheet = g.xlFile.Sheet[splat[0]]
		g.activeSheet = sheet
		rangeIDString = splat[1]
	} else if g.activeSheet != nil {
		sheet = g.activeSheet
	} else {
		sheet = g.xlFile.Sheets[0]
	}

	var colFrom, rowFrom, colTo, rowTo int
	var xlError error
	if strings.Contains(rangeIDString, ":") {
		splat := strings.Split(rangeIDString, ":")
		fromIDString = splat[0]
		toIDString = splat[1]
	} else {
		fromIDString = rangeIDString
		toIDString = rangeIDString
	}

	if colFrom, rowFrom, xlError = xlsx.GetCoordsFromCellIDString(fromIDString); xlError != nil {
		err = xlError
		return
	}
	if colTo, rowTo, xlError = xlsx.GetCoordsFromCellIDString(toIDString); xlError != nil {
		err = xlError
		return
	}

	rowCount := rowTo - rowFrom + 1
	colCount := colTo - colFrom + 1
	cellRange = Range{
		cells:    make([]Cell, rowCount*colCount),
		rowCount: rowCount,
		colCount: colCount,
	}
	for i := 0; i < rowCount; i++ {
		for j := 0; j < colCount; j++ {
			xlCell := sheet.Cell(rowFrom+i, colFrom+j)
			cell := Cell{}
			if f, err := strconv.ParseFloat(xlCell.Value, 64); err != nil {
				cell.value = xlCell.Value
			} else {
				cell.value = f
			}

			if formula := xlCell.Formula(); formula != "" {
				cell.formula = `=` + formula
			}
			cellRange.cells[i*colCount+j] = cell
		}
	}

	return
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
		g.activeSheet = sheet
		cellIDString = splat[1]
	} else if g.activeSheet != nil {
		sheet = g.activeSheet
	} else {
		sheet = g.xlFile.Sheets[0]
	}

	if col, row, xlError := xlsx.GetCoordsFromCellIDString(cellIDString); xlError != nil {
		err = xlError
		return
	} else {
		xlCell := sheet.Cell(row, col)
		logger.Printf("Cell: %s, fmt: %s", cellIDString, xlCell.NumFmt)

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

func (g *Engine) Inspect() map[string]string {
	return map[string]string{
		"stackHeight": fmt.Sprintf("%d", g.callstack.Len()),
	}
}

// Execute the g
func (g *Engine) Execute(inputs map[string]string, outputs *map[string]string) (err error) {
	for cellID, value := range inputs {
		g.SetCell(cellID, value)
	}

	for cellID := range *outputs {
		var cell Cell
		cell, err = g.GetCell(cellID)
		if err != nil {
			logger.Printf("***Could not get cell %s. Reason: %v\n", cellID, err)
			return
		}
		if cell.formula != "" {
			logger.Printf("Formula: %s\n\n", cell.formula)
			formula := f1F.NewFormula(cell.formula)
			value, _ := g.EvalFormula(formula)

			(*outputs)[cellID] = fmt.Sprintf("%v", value)
		} else {
			(*outputs)[cellID] = fmt.Sprintf("%v", cell.value)
		}
	}

	return
}

// EvalFormula Execute formula by running AST nodes as necessary
func (g *Engine) EvalFormula(f *f1F.Formula) (value interface{}, valueType f1F.NodeType) {
	var currentNode *f1F.Node

	currentNode = f.GetEntryNode()
	logger.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n")
	logger.Printf("Evaluating formula...\n")
	logger.Printf("- Entry: %v\n", currentNode.Value())

	stackHeight := g.callstack.Len()
	err := g.evalNode(currentNode)
	if g.callstack.Len() != stackHeight {
		// panic(errors.New(fmt.Sprintf("Stack not disposed properly: was %d, now %d",
		// 	stackHeight, g.callstack.Len())))
		logger.Printf("***Stack not disposed properly: was %d, now %d ***\n",
			stackHeight, g.callstack.Len())
	}

	if err != nil {
		value = err
		valueType = 0
		return
	} else if g.ax == nil {
		value = "#NA"
		valueType = 0
		return
	}

	value = g.ax
	logger.Printf("f() = %v\n", value)
	logger.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n")
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

// SetCell Set value for a cell
func (g *Engine) SetCell(cellID string, value interface{}) {
	var sheetName string
	if strings.Contains(cellID, "!") {
		splat := strings.Split(cellID, "!")
		sheetName = splat[0]
		cellID = splat[1]
	} else {
		sheetName = "Input"
	}

	sheet := g.xlFile.Sheet[sheetName]
	if col, row, err := xlsx.GetCoordsFromCellIDString(cellID); err != nil {
		return
	} else {
		cell := sheet.Cell(row, col)
		cell.SetValue(value)
	}
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
	if invoke.fn == "IDENTITY" {
		var operand interface{}
		g.pop(&operand)
		ret = operand
	} else if invoke.fn == "+" {
		var operand interface{}
		var ax float64

		for i := invoke.arity; i >= 1; i-- {
			g.pop(&operand)
			ax += operand.(float64)
		}
		ret = ax
	} else if invoke.fn == "-" {
		var operand interface{}
		var ax float64 = 0.0

		for i := invoke.arity; i >= 2; i-- {
			g.pop(&operand)
			switch operand.(type) {
			case float64:
				ax += operand.(float64)
				break
			}
		}

		g.pop(&operand)
		switch operand.(type) {
		case float64:
			ret = operand.(float64) - ax
			break
		default:
			ret = 0.0
			break
		}
	} else if invoke.fn == "*" {
		var operand interface{}
		var ax float64 = 1.0

		for i := invoke.arity; i >= 1; i-- {
			g.pop(&operand)
			ax *= operand.(float64)
		}
		ret = ax
	} else if invoke.fn == "/" {
		var operand interface{}
		var ax float64 = 1

		for i := invoke.arity; i >= 2; i-- {
			g.pop(&operand)
			ax *= operand.(float64)
		}
		g.pop(&operand)
		ret = operand.(float64) / ax
	} else if strings.Contains(">=<=", invoke.fn) {
		var operand1, operand2 interface{}
		g.pop(&operand2)
		g.pop(&operand1)

		switch operand1.(type) {
		case error:
			ret = operand1.(error)
			break
		case float64:
			switch operand2.(type) {
			case error:
				ret = operand2.(error)
				break
			case float64:
				ret = g.logicalFloat64(invoke.fn, operand1.(float64), operand2.(float64))
				break
			case string:
				panic(fmt.Errorf("Operation not supported: float64 %f vs string '%s'",
					operand1.(float64), operand2.(string)))
				break
			}
			break
		case string:
			switch operand2.(type) {
			case error:
				ret = operand2.(error)
				break
			case float64:
				panic(fmt.Errorf("Operation not supported: string '%s' vs float64 %f",
					operand1.(string), operand2.(float64)))
				break
			case string:
				ret = g.logicalString(invoke.fn, operand1.(string), operand2.(string))
				break
			}
			break
		}

	} else {
		// Non primitive operators: + - * /
		// NOTE: DO NOT REFACTOR INTO DYNAMIC METHOD CALLING WITH ARGS...
		if invoke.arity == 1 {
			var operand1 interface{}
			g.pop(&operand1)
			logger.Printf("Call1: %s, %v\n", invoke.fn, operand1)
			if output, err := funs.Call1(invoke.fn, operand1); err != nil {
				ret = err
			} else {
				ret = output
			}
		} else if invoke.arity == 2 {
			var operand1, operand2 interface{}
			g.pop(&operand2)
			g.pop(&operand1)
			logger.Printf("Call2: %s, %v, %v\n", invoke.fn, operand1, operand2)
			if output, err := funs.Call2(invoke.fn, operand1, operand2); err != nil {
				ret = err
			} else {
				ret = output
			}
		} else if invoke.arity == 3 {
			var operand1, operand2, operand3 interface{}
			g.pop(&operand3)
			g.pop(&operand2)
			g.pop(&operand1)
			logger.Printf("Call3: %s, %v, %v, %v\n", invoke.fn, operand1, operand2, operand3)
			if output, err := funs.Call3(invoke.fn, operand1, operand2, operand3); err != nil {
				ret = err
			} else {
				ret = output
			}
		} else if invoke.arity == 4 {
			var operand1, operand2, operand3, operand4 interface{}
			g.pop(&operand4)
			g.pop(&operand3)
			g.pop(&operand2)
			g.pop(&operand1)
			logger.Printf("Call4: %s, %v, %v, %v, %v\n", invoke.fn, operand1, operand2, operand3, operand4)
			if output, err := funs.Call4(invoke.fn, operand1, operand2, operand3, operand4); err != nil {
				ret = err
			} else {
				ret = output
			}
		}
	}
	// NOTE: Remember to g.pop after g.runStack
	g.push(ret)
}

func (g *Engine) evalNode(node *f1F.Node) (err error) {
	switch node.NodeType() {
	case f1F.NodeTypeOperator:
		if err = g.callFunc(node); err != nil {
			return
		}
		break
	case f1F.NodeTypeFunc:
		if node.Value() == "TRUE" {
			g.ax = true
		} else if node.Value() == "FALSE" {
			g.ax = false
		} else {
			err = g.callFunc(node)
			if err != nil {
				return
			}
		}
		break
	case f1F.NodeTypeRef:
		g.callDeref(node)
		break
	case f1F.NodeTypeLiteral, f1F.NodeTypeFloat, f1F.NodeTypeInteger:
		g.ax = node.Value()
		break
	}

	return
}

func (g *Engine) logicalFloat64(fn string, operand1 float64, operand2 float64) bool {
	if fn == ">" {
		return operand1 > operand2
	} else if fn == "<" {
		return operand1 < operand2
	} else if fn == "=" {
		return operand1 == operand2
	} else if fn == ">=" {
		return operand1 >= operand2
	} else if fn == "<=" {
		return operand1 <= operand2
	} else {
		panic(errors.New(fmt.Sprintf("Invalid operator %s", fn)))
	}
}

func (g *Engine) logicalString(fn string, operand1 string, operand2 string) bool {
	if fn == ">" {
		return operand1 > operand2
	} else if fn == "<" {
		return operand1 < operand2
	} else if fn == "=" {
		return operand1 == operand2
	} else if fn == ">=" {
		return operand1 >= operand2
	} else if fn == "<=" {
		return operand1 <= operand2
	} else {
		panic(errors.New(fmt.Sprintf("Invalid operator %s", fn)))
	}
}

// run Expand an AST node and pop result from stack to ax
func (g *Engine) callFunc(node *f1F.Node) (err error) {
	var fn string
	fn = node.Value().(string)

	if fn == "IF" { // The IF-JUMP
		if g.callIf(node.FirstChild()) {
			err = g.evalNode(node.ChildAt(1))
		} else if falseBranch := node.ChildAt(2); falseBranch != nil {
			err = g.evalNode(falseBranch)
		} else {
			g.ax = false
		}

		if err != nil {
			return
		}

		return
	}

	if !strings.Contains("IDENTITY+-*/>=<=", fn) && !funs.Exists(fn) {
		println(fmt.Sprintf("Function not exists: %s", fn))
		err = errors.New(fmt.Sprintf("Function not exists: %s", fn))
		return
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
			break
		case f1F.NodeTypeLiteral:
			value := childNode.Value()
			g.push(value)
			break
		case f1F.NodeTypeFloat:
			value := childNode.Value()
			g.push(value)
			break
		case f1F.NodeTypeOperator:
			err = g.callFunc(childNode)
			if err != nil {
				return
			}
			g.push(g.ax)
			break
		case f1F.NodeTypeFunc:
			stackHeight := g.callstack.Len()
			err = g.callFunc(childNode)
			if stackHeight != g.callstack.Len() {
				logger.Printf("***Stack corruption: was %d, now %d***\n", stackHeight, g.callstack.Len())
			}

			if err != nil {
				return
			}
			g.push(g.ax)
			break
		}
	}
	g.runStack(invoke.(*Invoke))
	g.pop(&g.ax) // Leave the stack

	return
}

func (g *Engine) callDeref(node *f1F.Node) {
	cellIDString := node.Value().(string)
	activeSheet := g.activeSheet

	if strings.Contains(cellIDString, ":") {
		// Request for a range, even for single dimension ranges
		if cellRange, err := g.GetRange(cellIDString); err != nil {
			logger.Printf("Could not deref %s. Reason: %v", cellIDString, err)
			return
		} else {
			if cells, ok := cellRange.ToSlice(); ok {
				result := make([]interface{}, len(cells))
				for i := range result {
					if formulaString := cells[i].formula; formulaString != "" {
						logger.Printf("Evaluating cell[%d]: %s, f(x) %s\n", i, cellIDString, formulaString)
						formula := f1F.NewFormula(formulaString)
						g.EvalFormula(formula) // g.ax is updated
						result[i] = g.ax
					} else {
						result[i] = cells[i].value
					}
				}
				g.ax = result
			} else if cells, ok := cellRange.To2DSlice(); ok {
				result := make([][]interface{}, cellRange.rowCount)
				colCount := cellRange.colCount
				for i := 0; i < cellRange.rowCount; i++ {
					result[i] = make([]interface{}, colCount)
					for j := 0; j < cellRange.colCount; j++ {
						if formulaString := cells[i][j].formula; formulaString != "" {
							formula := f1F.NewFormula(formulaString)
							g.EvalFormula(formula) // g.ax is updated
							result[i][j] = g.ax
						} else {
							result[i][j] = cells[i][j].value
						}

					}
				}
				g.ax = result
			}
		}

	} else {
		if cell, err := g.GetCell(cellIDString); err != nil {
			logger.Printf("Could not deref %s. Reason: %v\n", cellIDString, err)
			g.activeSheet = activeSheet
			return
		} else if cell.formula != "" {
			logger.Printf("FORMULA: %s\n", cell.formula)
			formula := f1F.NewFormula(cell.formula)
			g.EvalFormula(formula) // g.ax is updated
		} else if cell.value != "" {
			g.ax = cell.value
		}
	}
	logger.Printf(">>>>>\n")
	if strings.Contains(cellIDString, "!") {
		logger.Printf("Deref'd cell(s): %s = %v\n", cellIDString, g.ax)
	} else if activeSheet != nil {
		logger.Printf("Deref'd cell(s): %s!%s = %v\n", activeSheet.Name, cellIDString, g.ax)
	} else {
		logger.Printf("Deref'd cell(s): default first sheet! %s = %v\n", cellIDString, g.ax)
	}
	logger.Printf("Stack height: %d\n", g.callstack.Len())
	logger.Printf("<<<<<\n")
	g.activeSheet = activeSheet
}

// callIf Evaluate a node to a bool in MS-EXCEL
// - False: nil, false, 0, error
// - True: anything else
func (g *Engine) callIf(node *f1F.Node) bool {
	if err := g.evalNode(node); err != nil {
		return false
	}

	switch g.ax.(type) {
	case error:
		return false
	case bool:
		return g.ax.(bool)
	case int:
		if g.ax.(int) == 0 {
			return false
		}
	case float64:
		if g.ax.(float64) == 0.0 {
			return false
		}
	case string:
		if g.ax.(string) == "FALSE" {
			return false
		}
	default:
		if g.ax == nil {
			return false
		}

		return true
	}
	return true
}
