package engine

import (
	"os"
	"math"
	"testing"
	"github.com/tealeg/xlsx"
	f1Formula "github.com/khanhhua/formula1/formula"
)

var EPSILON = math.Pow(10, -9)
var xlFile *xlsx.File

func TestMain(m *testing.M) {
	// setup
	var err error
	xlFile, err = xlsx.OpenFile("/Users/khanhhua/Downloads/formula1-x1.xlsx")
	if err != nil {
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestNumberLiteral(t *testing.T) {
	engine := NewEngine(xlFile)
	formula := f1Formula.NewFormula(`=10.1`)

	var result interface{}
	result = engine.EvalFormula(formula)
	if result != 10.1 {
		t.Errorf("Expected: 10.1\tActual: %v", result)
	}
}

func TestAdditionOf2Literal(t *testing.T) {
	engine := NewEngine(xlFile)
	formula := f1Formula.NewFormula(`=1.1 + 2.2`)

	var result interface{}
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 3.3) > EPSILON {
		t.Errorf("Expected: 3.3\tActual: %v", result)
	}
}

func TestAdditionOf3Literals(t *testing.T) {
	engine := NewEngine(xlFile)
	formula := f1Formula.NewFormula(`=1.1 + 2.2 + 10`)

	var result interface{}
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 13.3) > EPSILON {
		t.Errorf("Expected: 13.3\tActual: %v", result)
	}
}

func TestArithOfLiterals(t *testing.T) {
	var engine *Engine
	var formula *f1Formula.Formula
	var result interface{}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 - 29`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 1) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=20 - 29 + 10`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 1) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 + 30 - 1`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 59) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 - 1 + 20 + 30`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 59) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 + 30 - 1 - 2`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 57) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 - 1 - 2 + 20 + 30`)
	result = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r - 57) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}
}
