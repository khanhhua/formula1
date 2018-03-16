package engine

import f1Formula "github.com/khanhhua/formula1/formula"
import "testing"

func TestNumberLiteral(t *testing.T) {
	engine := NewEngine(nil)
	formula := f1Formula.NewFormula(`=10.1`)

	var result interface{}
	result = engine.EvalFormula(formula)
	if result != 10.1 {
		t.Errorf("Expected: 10.1\tActual: %f", result)
	}
}
