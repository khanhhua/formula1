package funs

import "testing"

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
