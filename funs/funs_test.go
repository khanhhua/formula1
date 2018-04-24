package funs

import (
	"errors"
	"testing"
)

func TestMATCH(t *testing.T) {
	lookupRange := []interface{}{2, 4, 6, 8, 10}
	if result := MATCH(2, lookupRange, 0); result != 1 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}
	if result := MATCH(3, lookupRange, 1); result != 1 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}
	if result := MATCH(3, lookupRange, -1); result != 2 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}
}

func TestIFERROR(t *testing.T) {
	if result := IFERROR(1.1, 2.2); result != 1.1 {
		t.Errorf("Expected: 1.1\tActual:%v", result)
	}

	if result := IFERROR(errors.New("#ERROR"), 2.2); result != 2.2 {
		t.Errorf("Expected: 2.2\tActual:%v", result)
	}
}

func TestFLOOR(t *testing.T) {
	if result := FLOOR(1.1); result != 1 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}
}

func TestSUM(t *testing.T) {
	if result := SUM(2.5); result != 2.5 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	arr1D := []float64{1, 2, 3, 4, 5}
	if result := SUM(arr1D); result != 15 {
		t.Errorf("Expected: 15\tActual:%v", result)
	}

	arr2D := [][]float64{{1, 2}, {3, 4}, {5, 6}}
	if result := SUM(arr2D); result != 21 {
		t.Errorf("Expected: 21\tActual:%v", result)
	}
}

func TestPOWER(t *testing.T) {
	if result := POWER(2.0, 3.0); result != 8 {
		t.Errorf("Expected: 8\tActual:%v", result)
	}
}

func TestVLOOKUPfloat64(t *testing.T) {
	var lookupRange = make([][]interface{}, 2)
	lookupRange[0] = make([]interface{}, 1)
	lookupRange[0][0] = 1.0

	lookupRange[1] = make([]interface{}, 1)
	lookupRange[1][0] = 2.0

	if result := VLOOKUP(1.0, lookupRange, 1, false); result != 1.0 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	lookupRange = make([][]interface{}, 2)
	lookupRange[0] = make([]interface{}, 3)
	lookupRange[0][0] = 1.0
	lookupRange[0][1] = 10.0
	lookupRange[0][2] = 100.0

	lookupRange[1] = make([]interface{}, 3)
	lookupRange[1][0] = 21.0
	lookupRange[1][1] = 22.0
	lookupRange[1][2] = 23.0

	if result := VLOOKUP(21.0, lookupRange, 3, false); result != 23.0 {
		t.Errorf("Expected: 23.0\tActual:%v", result)
	}

	if result := VLOOKUP(99.0, lookupRange, 3, false); result.(error).Error() != "N/A" {
		t.Errorf("Expected: N/A\tActual:%v", result)
	}
}

func TestVLOOKUPfloat64Approx(t *testing.T) {
	var lookupRange = make([][]interface{}, 3)
	lookupRange[0] = make([]interface{}, 1)
	lookupRange[0][0] = 1.0

	lookupRange[1] = make([]interface{}, 1)
	lookupRange[1][0] = 2.0

	lookupRange[2] = make([]interface{}, 1)
	lookupRange[2][0] = 3.0

	if result := VLOOKUP(1.5, lookupRange, 1, true); result != 1.0 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	lookupRange = make([][]interface{}, 3)
	lookupRange[0] = make([]interface{}, 3)
	lookupRange[0][0] = 0
	lookupRange[0][1] = "0-17"
	lookupRange[0][2] = 72

	lookupRange[1] = make([]interface{}, 3)
	lookupRange[1][0] = 18
	lookupRange[1][1] = "18-30"
	lookupRange[1][2] = 72

	lookupRange[2] = make([]interface{}, 3)
	lookupRange[2][0] = 31
	lookupRange[2][1] = "31-40"
	lookupRange[2][2] = 72

	if result := VLOOKUP(0, lookupRange, 3, true); result != 72 {
		t.Errorf("Expected: 72\tActual: %v", result)
	}

	if result := VLOOKUP(18, lookupRange, 3, true); result != 72 {
		t.Errorf("Expected: 72\tActual: %v", result)
	}

	if result := VLOOKUP(19, lookupRange, 3, true); result != 72 {
		t.Errorf("Expected: 72\tActual: %v", result)
	}
}

func TestVLOOKUPstring(t *testing.T) {
	var lookupRange = make([][]interface{}, 2)
	lookupRange[0] = make([]interface{}, 1)
	lookupRange[0][0] = "Tom"

	lookupRange[1] = make([]interface{}, 1)
	lookupRange[1][0] = "Jerry"

	if result := VLOOKUP("Tom", lookupRange, 1, false); result != "Tom" {
		t.Errorf("Expected: Tom\tActual:%v", result)
	}

	lookupRange = make([][]interface{}, 2)
	lookupRange[0] = make([]interface{}, 3)
	lookupRange[0][0] = "Tom"
	lookupRange[0][1] = "Cat"
	lookupRange[0][2] = "Blue"

	lookupRange[1] = make([]interface{}, 3)
	lookupRange[1][0] = "Jerry"
	lookupRange[1][1] = "Brown"
	lookupRange[1][2] = "Rat"

	if result := VLOOKUP("Jerry", lookupRange, 3, false); result != "Rat" {
		t.Errorf("Expected: Rat\tActual:%v", result)
	}
}

func TestCOUNTIFfloat64(t *testing.T) {
	var lookupRange = make([]interface{}, 2)

	lookupRange[0] = 1.0
	lookupRange[1] = 2.0

	if result := COUNTIF(lookupRange, 1.0); result != 1.0 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	var lookupRange2D = make([][]interface{}, 2)
	lookupRange2D[0] = make([]interface{}, 1)
	lookupRange2D[0][0] = 1.0

	lookupRange2D[1] = make([]interface{}, 1)
	lookupRange2D[1][0] = 2.0

	if result := COUNTIF(lookupRange2D, 1.0); result != 1.0 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	lookupRange2D = make([][]interface{}, 2)
	lookupRange2D[0] = make([]interface{}, 3)
	lookupRange2D[0][0] = 1.0
	lookupRange2D[0][1] = 10.0
	lookupRange2D[0][2] = 100.0

	lookupRange2D[1] = make([]interface{}, 3)
	lookupRange2D[1][0] = 21.0
	lookupRange2D[1][1] = 22.0
	lookupRange2D[1][2] = 23.0

	if result := COUNTIF(lookupRange2D, 1.0); result != 1.0 {
		t.Errorf("Expected: 1.0\tActual:%v", result)
	}

	if result := COUNTIF(lookupRange2D, 2.0); result != 0.0 {
		t.Errorf("Expected: 0.0\tActual:%v", result)
	}
}

func TestCOUNTIFstring(t *testing.T) {
	var lookupRange = make([]interface{}, 2)

	lookupRange[0] = "hello"
	lookupRange[1] = "World"

	if result := COUNTIF(lookupRange, "hello"); result != 1.0 {
		t.Errorf("Expected: 1\tActual:%v", result)
	}

	var lookupRange2D = make([][]interface{}, 2)
	lookupRange2D[0] = make([]interface{}, 1)
	lookupRange2D[0][0] = "hello"

	lookupRange2D[1] = make([]interface{}, 1)
	lookupRange2D[1][0] = "World"

	if result := COUNTIF(lookupRange2D, "hello"); result != 1.0 {
		t.Errorf("Expected: hello\tActual:%v", result)
	}

	lookupRange2D = make([][]interface{}, 2)
	lookupRange2D[0] = make([]interface{}, 3)
	lookupRange2D[0][0] = "Tom"
	lookupRange2D[0][1] = "Cartoon"
	lookupRange2D[0][2] = "Blue"

	lookupRange2D[1] = make([]interface{}, 3)
	lookupRange2D[1][0] = "Jerry"
	lookupRange2D[1][1] = "Cartoon"
	lookupRange2D[1][2] = "Brown"

	if result := COUNTIF(lookupRange2D, "Cartoon"); result != 2.0 {
		t.Errorf("Expected: 2.0\tActual:%v", result)
	}

	if result := COUNTIF(lookupRange2D, 2.0); result != 0.0 {
		t.Errorf("Expected: 0.0\tActual:%v", result)
	}
}
