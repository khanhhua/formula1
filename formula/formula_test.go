package engine

import "testing"

func TestLiteral(t *testing.T) {
	formulaText := `=10`

	formula := NewFormula(formulaText)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
}

func TestInfixAtRoot(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=10 + 20`)
	if formula.root.firstChild().childCount() != 2 {
		t.Errorf("Formula root has 2 children")
	}

	formula = NewFormula(`=10.0 + A1`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.firstChild().children) != 2 {
		t.Errorf("Formula root has two children")
	}

	formula = NewFormula(`=A1 + 10.0`)
	if len(formula.root.firstChild().children) != 2 {
		t.Errorf("Formula root has two children")
	}
}

func TestMultiInfixAtRoot(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=10 + 20 + 30`)
	if formula.root.firstChild().childCount() != 3 {
		t.Errorf("Formula root has 2 children")
	}
}

func TestSingleCellRef(t *testing.T) {
	formulaText := `=A1`

	formula := NewFormula(formulaText)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
}

func TestSimpleFunctionWithLiteral(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=FLOOR(10.1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeLiteral {
		t.Errorf("First child is a literal")
	}
}

func TestSumSingleRange(t *testing.T) {
	var formula *Formula
	formula = NewFormula(`=SUM(A1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.firstChild().children) != 1 {
		t.Errorf("Formula root has 1 child")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeRef {
		t.Errorf("SUM has one ref parameter")
	}

	formula = NewFormula(`=SUM(A1:B2)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}

	if len(formula.root.firstChild().children) != 1 {
		t.Errorf("Formula root has 1 child")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeRef {
		t.Errorf("SUM has one ref parameter")
	}
}

func TestSimpleFunctionWithInfixExpression(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=FLOOR(10.1 + 20.2)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeOperator {
		t.Errorf("First child is an operator")
	}
	if formula.root.firstChild().firstChild().firstChild().nodeType != NodeTypeLiteral {
		t.Errorf("1st infix operand is a literal")
	}
	if formula.root.firstChild().firstChild().children[1].nodeType != NodeTypeLiteral {
		t.Errorf("2nd infix operand is a literal")
	}

	formula = NewFormula(`=FLOOR(10.1 + A1)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().nodeType != NodeTypeFunc {
		t.Errorf("First child is a function")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeOperator {
		t.Errorf("First child is an operator")
	}
	if formula.root.firstChild().firstChild().firstChild().nodeType != NodeTypeLiteral {
		t.Errorf("1st infix operand is a literal")
	}
	if formula.root.firstChild().firstChild().children[1].nodeType != NodeTypeRef {
		t.Errorf("2nd infix operand is a ref")
	}
}

func TestMultiParamFunctions(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(TRUE, 10, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().childCount() != 3 {
		t.Errorf("Formula root has three children")
	}
	if formula.root.firstChild().children[1].nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a literal")
	}
	if formula.root.firstChild().children[2].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}

	formula = NewFormula(`=IF(TRUE, A1, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().childCount() != 3 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.firstChild().children[1].nodeType != NodeTypeRef {
		t.Errorf("True branch is a range")
	}
	if formula.root.firstChild().children[2].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}

	formula = NewFormula(`=OR(TRUE, FALSE)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().childCount() != 2 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a range")
	}
	if formula.root.firstChild().children[1].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}
}

func TestNestedFunctions(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(OR(TRUE,FALSE), 10, 20)`)
	if formula.root.value != "root" {
		t.Errorf("Formula root's value must be 'root'")
	}
	if formula.root.firstChild().childCount() != 3 {
		t.Errorf("Formula root has two children")
	}
	if formula.root.firstChild().firstChild().nodeType != NodeTypeFunc {
		t.Errorf("Condition is a function")
	}
	if formula.root.firstChild().children[1].nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a literal")
	}
	if formula.root.firstChild().children[2].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}

	condition := formula.root.firstChild().firstChild()
	if condition.value != "OR" {
		t.Errorf("Condition is an OR")
	}
	if condition.childCount() != 2 {
		t.Errorf("Formula root has two children")
	}
	if condition.firstChild().nodeType != NodeTypeLiteral {
		t.Errorf("True branch is a range")
	}
	if condition.children[1].nodeType != NodeTypeLiteral {
		t.Errorf("False branch is a literal")
	}
}

func TestNestedFunctionsWithInfixes(t *testing.T) {
	var formula *Formula

	formula = NewFormula(`=IF(OR(A1 > 1,FALSE), 10, 20)`)
	if formula.root.firstChild().firstChild().nodeType != NodeTypeFunc {
		t.Errorf("Condition is a function")
	}

	condition := formula.root.firstChild().firstChild()
	if condition.value != "OR" {
		t.Errorf("Condition is an OR")
	}
	if condition.childCount() != 2 {
		t.Errorf("Formula root has two children. Actual: %d", condition.childCount())
	}
	if condition.firstChild().nodeType != NodeTypeOperator {
		t.Errorf("True branch is an operator. Actual: %d", condition.firstChild().nodeType)
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
