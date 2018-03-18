package funs

import (
	"errors"
	"fmt"
)

var arity1map = map[string]func(input interface{}) float64{
	"FLOOR": FLOOR,
	"SUM":   SUM,
}

var arity2map = map[string]func(input1 interface{}, input2 interface{}) float64{
	"SUM":   SUM2,
	"POWER": POWER,
}

func Call1(name string, input interface{}) (ret interface{}, err error) {
	fn := arity1map[name]
	if fn == nil {
		err = errors.New(fmt.Sprintf("Invalid fun %s", name))
		return
	}

	return fn(input), nil
}

func Call2(name string, input1 interface{}, input2 interface{}) (ret interface{}, err error) {
	fn := arity2map[name]
	if fn == nil {
		err = errors.New(fmt.Sprintf("Invalid fun %s", name))
		return
	}

	return fn(input1, input2), nil
}
