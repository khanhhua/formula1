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
