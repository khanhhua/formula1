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
	// Execution stack
	callstack *stack.Stack
	// Execute and remember stuff here
	cache map[string]interface{}
}

// NewEngine Create a new engine to execute formula suitable for xlFile
func NewEngine(xlFile *xlsx.File) *Engine {
	return &Engine{
		xlFile: xlFile,
	}
}

// Execute the engine
func (engine *Engine) Execute(inputs map[string]string) (names []string, outputs map[string]string) {
	// formula := NewFormula(formulaText)

	return
}

func (engine *Engine) EvalFormula(f *f1F.Formula) interface{} {
	// The registers
	var ax interface{}
	var cci int = 0
	var maxCCI int
	var currentNode *f1F.Node

	currentNode = f.GetEntryNode()
	codes := []*f1F.Node{currentNode}
	if currentNode.HasChildren() {
		codes = append(codes, currentNode.Children()...)
	}
	maxCCI = len(codes) - 1
	fmt.Printf("Evaluating formula...\n")
	fmt.Printf("- Entry: %s\n", currentNode.Value())
	fmt.Printf("- MaxCCI: %d\n", maxCCI)

	// The mighty loop :D... LOOP LOOP LOOP until cci reaches maxCCI
	for cci <= maxCCI {
		currentNode = codes[cci]
		nodeType := currentNode.NodeType()

		if nodeType == f1F.NodeTypeFunc {

		} else if nodeType == f1F.NodeTypeRef {

		} else if nodeType == f1F.NodeTypeLiteral {
			ax = currentNode.Value()
		} else if nodeType == f1F.NodeTypeOperator {

		}

		cci++
	}

	if ax == nil {
		return "#NA"
	} else {
		switch ax.(type) {
		case string:
			return ax.(string)
		case int32:
			return ax.(int32)
		case float32:
			return ax.(float32)
		case float64:
			return ax.(float64)
		}
	}

	return "#NA"
}
