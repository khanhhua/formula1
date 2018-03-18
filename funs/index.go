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

var a2float64map = map[string]func(interface{}, interface{}) float64{
	"SUM":   SUM2,
	"POWER": POWER,
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
	}

	err = fmt.Errorf("Invalid fun %s", name)
	return
}
