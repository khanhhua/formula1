package funs

import "testing"

func TestCall1(t *testing.T) {
	if result, err := Call1("OR", 1.1); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(bool) != true {
		t.Errorf("Expected: 1\tActual: %v", result)
	}

	if result, err := Call1("AND", 1.1); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(bool) != true {
		t.Errorf("Expected: 1.1\tActual: %v", result)
	}

	if result, err := Call1("FLOOR", 1.1); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(float64) != 1 {
		t.Errorf("Expected: 1\tActual: %v", result)
	}

	if result, err := Call1("SUM", 1.1); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(float64) != 1.1 {
		t.Errorf("Expected: 1.1\tActual: %v", result)
	}
}

func TestCall2(t *testing.T) {
	if result, err := Call2("IFERROR", 1.1, 2.2); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(float64) != 1.1 {
		t.Errorf("Expected: 1\tActual: %v", result)
	}

	if result, err := Call2("OR", 1.1, 2.2); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(bool) != true {
		t.Errorf("Expected: 1\tActual: %v", result)
	}

	if result, err := Call2("AND", 1.1, 2.2); err != nil {
		t.Errorf("Call1 error. %v", err)
	} else if result.(bool) != true {
		t.Errorf("Expected: 1.1\tActual: %v", result)
	}

	if result, err := Call2("SUM", 10, 2); err != nil {
		t.Errorf("Call2 error. %v", err)
	} else if result.(float64) != 12 {
		t.Errorf("Expected: 12\tActual: %v", result)
	}

	if result, err := Call2("POWER", 10, 2); err != nil {
		t.Errorf("Call2 error. %v", err)
	} else if result.(float64) != 100 {
		t.Errorf("Expected: 100\tActual: %v", result)
	}

	var lookupRange = make([]interface{}, 2)

	lookupRange[0] = "hello"
	lookupRange[1] = "World"
	if result, err := Call2("COUNTIF", lookupRange, "hello"); err != nil {
		t.Errorf("Call2 error. %v", err)
	} else if result.(float64) != 1 {
		t.Errorf("Expected: 1\tActual: %v", result)
	}
}
