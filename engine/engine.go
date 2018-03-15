package formula1

import (
	"github.com/golang-collections/collections/stack"
	"github.com/tealeg/xlsx"
)

// Engine the F1 Engine
type Engine struct {
	xlFile *xlsx.File
	// Execution stack
	stack *stack.Stack
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
