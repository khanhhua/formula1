package funs

import (
	"fmt"
)

var a1boolmap = map[string]func(interface{}) bool{
	"OR":  OR,
	"AND": AND,
}

var a1float64map = map[string]func(interface{}) float64{
	"FLOOR": FLOOR,
	"SUM":   SUM,
}

var a2boolmap = map[string]func(interface{}, interface{}) bool{
	"OR":  OR2,
	"AND": AND2,
}

var a3boolmap = map[string]func(interface{}, interface{}, interface{}) bool{
	"OR":  OR3,
	"AND": AND3,
}

var a2float64map = map[string]func(interface{}, interface{}) float64{
	"SUM":   SUM2,
	"POWER": POWER,
	"ROUND": func(p1 interface{}, precision interface{}) float64 {
		return ROUND(p1.(float64), precision.(float64))
	},
	"COUNTIF": COUNTIF,
}

var a2int64map = map[string]func(interface{}, interface{}) int{
	"MATCH": func(p1, p2 interface{}) int {
		return MATCH(p1, p2, 1)
	},
}

var a3int64map = map[string]func(interface{}, interface{}, interface{}) int{
	"MATCH": func(p1, p2 interface{}, p3 interface{}) int {
		return MATCH(p1, p2, int(p3.(float64)))
	},
}

var a2inter = map[string]func(interface{}, interface{}) interface{}{
	"IFERROR": IFERROR,
}

var a4inter = map[string]func(interface{}, interface{}, interface{}, interface{}) interface{}{
	"VLOOKUP": func(p1 interface{}, p2 interface{}, p3 interface{}, p4 interface{}) interface{} {
		var index int
		var approx bool
		if result, ok := p3.(int); ok {
			index = result
		} else if result, ok := p3.(float64); ok {
			index = int(result)
		} else {
			return "N/A"
		}

		if result, ok := p4.(int); ok {
			approx = result == 1
		} else if result, ok := p4.(float64); ok {
			approx = result == 1.0
		} else {
			approx = false
		}

		return VLOOKUP(p1, p2, index, approx)
	},
}

// Exists Checks if a function has been implemented
func Exists(name string) bool {
	if _, ok := a1boolmap[name]; ok {
		return true
	} else if _, ok := a1float64map[name]; ok {
		return true
	} else if _, ok := a2boolmap[name]; ok {
		return true
	} else if _, ok := a2int64map[name]; ok {
		return true
	} else if _, ok := a3boolmap[name]; ok {
		return true
	} else if _, ok := a3int64map[name]; ok {
		return true
	} else if _, ok := a2float64map[name]; ok {
		return true
	} else if _, ok := a2inter[name]; ok {
		return true
	} else if _, ok := a4inter[name]; ok {
		return true
	} else {
		return false
	}
}

// Call1 Invoke arity-1 functions
func Call1(name string, input interface{}) (ret interface{}, err error) {
	if fn, ok := a1boolmap[name]; ok {
		return fn(input), nil
	} else if fn, ok := a1float64map[name]; ok {
		return fn(input), nil
	}

	err = fmt.Errorf("Invalid fun %s", name)
	return
}

// Call2 Invoke arity-2 functions
func Call2(name string, input1 interface{}, input2 interface{}) (ret interface{}, err error) {
	if fn, ok := a2boolmap[name]; ok {
		return fn(input1, input2), nil
	} else if fn, ok := a2float64map[name]; ok {
		return fn(input1, input2), nil
	} else if fn, ok := a2int64map[name]; ok {
		return fn(input1, input2), nil
	} else if fn, ok := a2inter[name]; ok {
		return fn(input1, input2), nil
	}

	err = fmt.Errorf("Invalid fun %s", name)
	return
}

// Call2 Invoke arity-2 functions
func Call3(name string, input1 interface{}, input2 interface{}, input3 interface{}) (ret interface{}, err error) {
	if fn, ok := a3boolmap[name]; ok {
		return fn(input1, input2, input3), nil
	} else if fn, ok := a3int64map[name]; ok {
		return fn(input1, input2, input3), nil
	}

	err = fmt.Errorf("Invalid fun %s", name)
	return
}

func Call4(name string, input1 interface{}, input2 interface{}, input3 interface{}, input4 interface{}) (ret interface{}, err error) {
	if fn, ok := a4inter[name]; ok {
		return fn(input1, input2, input3, input4), nil
	}

	err = fmt.Errorf("Invalid fun %s", name)
	return
}
