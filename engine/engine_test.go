package engine

import (
	"math"
	"os"
	"testing"

	f1Formula "github.com/khanhhua/formula1/formula"
	"github.com/tealeg/xlsx"
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
	result, _ = engine.EvalFormula(formula)
	if result != 10.1 {
		t.Errorf("Expected: 10.1\tActual: %v", result)
	}
}

func TestAdditionOf2Literal(t *testing.T) {
	engine := NewEngine(xlFile)
	formula := f1Formula.NewFormula(`=1.1 + 2.2`)

	var result interface{}
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-3.3) > EPSILON {
		t.Errorf("Expected: 3.3\tActual: %v", result)
	}
}

func TestAdditionOf3Literals(t *testing.T) {
	engine := NewEngine(xlFile)
	formula := f1Formula.NewFormula(`=1.1 + 2.2 + 10`)

	var result interface{}
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-13.3) > EPSILON {
		t.Errorf("Expected: 13.3\tActual: %v", result)
	}
}

func TestArithOfLiterals(t *testing.T) {
	var engine *Engine
	var formula *f1Formula.Formula
	var result interface{}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 - 29`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-1) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=20 - 29 + 10`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-1) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 + 30 - 1`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-59) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 - 1 + 20 + 30`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-59) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 + 20 + 30 - 1 - 2`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-57) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=10 - 1 - 2 + 20 + 30`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-57) > EPSILON {
		t.Errorf("Expected: 0\tActual: %v", result)
	}
}

func TestSimpleCellRef(t *testing.T) {
	var engine *Engine
	var formula *f1Formula.Formula
	var result interface{}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=B2`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-10) > EPSILON {
		t.Errorf("Expected: 10\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=Discounts!E2`)
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(string); !ok || r != "Cheap" {
		t.Errorf("Expected: 10\tActual: %v", result)
	}
}

func TestArithWithCellRef(t *testing.T) {
	var engine *Engine
	var formula *f1Formula.Formula
	var result interface{}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=B2 + 1`) // 10 + 1
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-11) > EPSILON {
		t.Errorf("Expected: 11\tActual: %v", result)
	}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=B2 + C2 + D2`) // 10 + 11 + 13
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(float64); !ok || (r-34) > EPSILON {
		t.Errorf("Expected: 34\tActual: %v", result)
	}
}

func TestIndirectCellRef(t *testing.T) {
	var engine *Engine
	var formula *f1Formula.Formula
	var result interface{}

	engine = NewEngine(xlFile)
	formula = f1Formula.NewFormula(`=Input!B3`) // B3=Discounts!E2
	result, _ = engine.EvalFormula(formula)
	if r, ok := result.(string); !ok || r != "Cheap" {
		t.Errorf("Expected: Cheap\tActual: %v", result)
	}
}
