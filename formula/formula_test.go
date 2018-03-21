package formula

import "testing"

func TestLiteral(t *testing.T) {
	formulaText := `=10`

	formula := NewFormula(formulaText)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild() == nil {
		t.Errorf("First child must not be nil")
	}
}

func TestInfixAtRoot(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=10 + 20`)
	if formula.root.FirstChild().ChildCount() != 2 {
		t.Errorf("Formula root has 2 children")
	}

	formula = NewFormula(`=10.0 + A1`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.FirstChild().children) != 2 {
		t.Errorf("Formula root has two children")
	}

	formula = NewFormula(`=A1 + 10.0`)
	if len(formula.root.FirstChild().children) != 2 {
		t.Errorf("Formula root has two children")
	}
}

func TestMultiOperandInfixAtRoot(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=10 + 20 + 30`)
	if formula.root.FirstChild().ChildCount() != 3 {
		t.Errorf("Formula root has 3 children")
	}
}

func TestMultiInfixOperatorAtRoot(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=10 + 20 - 30`)
	if result := formula.root.FirstChild().ChildCount(); result != 2 {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 2\tActual: %v", result)
	}
	if result := formula.root.FirstChild().Value().(string); result != "-" {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 2\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().ChildCount(); result != 2 {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 2\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().Value().(string); result != "+" {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 10\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().ChildAt(0).Value().(float64); result != 10 {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 10\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().ChildAt(1).Value().(float64); result != 20 {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 20\tActual: %v", result)
	}
	if result := formula.root.FirstChild().ChildAt(1).Value().(float64); result != 30 {
		t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 1\tActual: %v", result)
	}

	formula = NewFormula(`=10 + 20 + 30 - 1 - 2`)
	if result := formula.root.FirstChild().ChildCount(); result != 3 {
		t.Errorf("POSTFIX: ((10 20 30)+ 1 2)-. Expect: 3\tActual: %v", result)
	}
	if result := formula.root.FirstChild().Value().(string); result != "-" {
		t.Errorf("POSTFIX: ((10 20 30)+ 1 2)-. Expect: -\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().ChildCount(); result != 3 {
		t.Errorf("POSTFIX: ((10 20 30)+ 1 2)-. Expect: 3\tActual: %v", result)
	}
	if result := formula.root.FirstChild().FirstChild().Value().(string); result != "+" {
		t.Errorf("POSTFIX: ((10 20 30)+ 1 2)-. Expect: 10\tActual: %v", result)
	}
	// if result := formula.root.FirstChild().FirstChild().ChildAt(0).Value().(float64); result != 10 {
	// 	t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 10\tActual: %v", result)
	// }
	// if result := formula.root.FirstChild().FirstChild().ChildAt(1).Value().(float64); result != 20 {
	// 	t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 20\tActual: %v", result)
	// }
	// if result := formula.root.FirstChild().ChildAt(1).Value().(float64); result != 30 {
	// 	t.Errorf("POSTFIX: ((10 20)+ 30)-. Expect: 1\tActual: %v", result)
	// }
}

func TestSingleParentheses(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=2 * (5 - 1)`)

	if result := formula.root.ChildCount(); result != 1 {
		t.Errorf("POSTFIX: (2 (5 1)-)*. Expect: 1\tActual: %v", result)
	}
	if result := formula.root.FirstChild().ChildCount(); result != 2 {
		t.Errorf("POSTFIX: (2 (5 1)-)*. Expect: 2\tActual: %v", result)
	}
	if result := formula.root.FirstChild().Value().(string); result != "*" {
		t.Errorf("POSTFIX: (2 (5 1)-)*. Expect: 2\tActual: %v", result)
	}

	formula = NewFormula(`=(5 - 1) * 2`)

	if result := formula.root.ChildCount(); result != 1 {
		t.Errorf("POSTFIX: ((5 1)- 2)*. Expect: 1\tActual: %v", result)
	}
	if result := formula.root.FirstChild().ChildCount(); result != 2 {
		t.Errorf("POSTFIX: ((5 1)- 2)*. Expect: 2\tActual: %v", result)
	}
	if result := formula.root.FirstChild().Value().(string); result != "*" {
		t.Errorf("POSTFIX: ((5 1)- 2)*. Expect: 2\tActual: %v", result)
	}
}

func TestSingleCellRef(t *testing.T) {
	formulaText := `=A1`

	formula := NewFormula(formulaText)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if result := formula.root.FirstChild().NodeType(); result != NodeTypeRef {
		t.Errorf("Expect: NodeTypeRef\tActual: %v", result)
	}
}

func TestSimpleFunctionWithLiteral(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=FLOOR(10.1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeFloat {
		t.Errorf("First child is a literal")
	}
}

func TestSumSingleRange(t *testing.T) {
	var formula *Formula
	formula = NewFormula(`=SUM(A1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.FirstChild().children) != 1 {
		t.Errorf("Formula root has 1 child")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeRef {
		t.Errorf("SUM has one ref parameter")
	}

	formula = NewFormula(`=SUM(A1:B2)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.FirstChild().children) != 1 {
		t.Errorf("Formula root has 1 child")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeRef {
		t.Errorf("SUM has one ref parameter")
	}
}

func TestSimpleFunctionWithInfixExpression(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=FLOOR(10.1 + 20.2)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeOperator {
		t.Errorf("First child is an operator")
	}
	if formula.root.FirstChild().FirstChild().FirstChild().nodeType != NodeTypeFloat {
		t.Errorf("1st infix operand is a float literal")
	}
	if formula.root.FirstChild().FirstChild().children[1].nodeType != NodeTypeFloat {
		t.Errorf("2nd infix operand is a float literal")
	}

	formula = NewFormula(`=FLOOR(10.1 + A1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeOperator {
		t.Errorf("First child is an operator")
	}
	if formula.root.FirstChild().FirstChild().FirstChild().nodeType != NodeTypeFloat {
		t.Errorf("1st infix operand is a float literal")
	}
	if formula.root.FirstChild().FirstChild().children[1].nodeType != NodeTypeRef {
		t.Errorf("2nd infix operand is a ref")
	}
}

func TestMultiParamFunctions(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(TRUE, 10, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().ChildCount() != 3 {
		t.Errorf("Formula root has three children")
	}
	if formula.root.FirstChild().children[1].nodeType != NodeTypeFloat {
		t.Errorf("True branch is a literal")
	}
	if formula.root.FirstChild().children[2].nodeType != NodeTypeFloat {
		t.Errorf("False branch is a literal")
	}

	formula = NewFormula(`=IF(TRUE, A1, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().ChildCount() != 3 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.FirstChild().children[1].nodeType != NodeTypeRef {
		t.Errorf("True branch is a range")
	}
	if formula.root.FirstChild().children[2].nodeType != NodeTypeFloat {
		t.Errorf("False branch is a float literal")
	}

	formula = NewFormula(`=OR(TRUE, FALSE)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().ChildCount() != 2 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a range")
	}
	if formula.root.FirstChild().children[1].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}
}

func TestNestedFunctions(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(OR(TRUE,FALSE), 10, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.FirstChild().ChildCount() != 3 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeFunc {
		t.Errorf("Condition is a function")
	}
	if formula.root.FirstChild().children[1].nodeType != NodeTypeFloat {
		t.Errorf("True branch is a float literal")
	}
	if formula.root.FirstChild().children[2].nodeType != NodeTypeFloat {
		t.Errorf("False branch is a float literal")
	}

	condition := formula.root.FirstChild().FirstChild()
	if condition.value != "OR" {
		t.Errorf("Condition is an OR")
	}
	if condition.ChildCount() != 2 {
		t.Errorf("Formula root has two children")
	}
	if condition.FirstChild().nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a range")
	}
	if condition.children[1].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}
}

func TestNestedFunctionsWithInfixes(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(OR(A1 > 1,FALSE), 10, 20)`)
	if formula.root.FirstChild().FirstChild().nodeType != NodeTypeFunc {
		t.Errorf("Condition is a function")
	}

	condition := formula.root.FirstChild().FirstChild()
	if condition.value != "OR" {
		t.Errorf("Condition is an OR")
	}
	if condition.ChildCount() != 2 {
		t.Errorf("Formula root has two children. Actual: %d", condition.ChildCount())
	}
	if condition.FirstChild().nodeType != NodeTypeOperator {
		t.Errorf("True branch is an operator. Actual: %d", condition.FirstChild().nodeType)
	}
	if condition.children[1].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal. Actual: %d", condition.children[1].nodeType)
	}
}

//
// func TestIfOrFormula(t *testing.T) {
// 	formulaText := `=IF(OR(CalculatorNB!$B$12="Decline",CalculatorNB!$B$12="Refer"),CalculatorNB!$B$12,CalculatorNB!E48)`
//
// 	formula := NewFormula(formulaText)
// 	if formula.root.value != "root" {
// 		t.Errorf("Formula root's value must be 'root'")
// 	}
// }
