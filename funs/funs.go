package funs

import "math"

func boolean(input interface{}) bool {
	switch input.(type) {
	case bool:
		return input.(bool)
	case float64:
		return input.(float64) > 0 || input.(float64) < 0
	default:
		return false
	}
}

// OR Evaluate to a boolean
func OR(input interface{}) bool {
	return boolean(input)
}

// OR2 Evaluate to a boolean
func OR2(input1 interface{}, input2 interface{}) bool {
	return OR(input1) || OR(input2)
}

// AND Evaluate to a boolean
func AND(input interface{}) bool {
	return boolean(input)
}

// AND2 Evaluate to a boolean
func AND2(input1 interface{}, input2 interface{}) bool {
	return boolean(input1) && boolean(input2)
}

// FLOOR Floor function
func FLOOR(input interface{}) float64 {
	return math.Floor(input.(float64))
}

// SUM Sum of single range
// - Single number
// - Single range
func SUM(input interface{}) float64 {
	switch input.(type) {
	case int:
		return float64(input.(int))
	case float64:
		return input.(float64)
	case []float64:
		sum := 0.0
		for _, item := range input.([]float64) {
			sum += item
		}
		return sum
	case [][]float64:
		sum := 0.0
		for _, items := range input.([][]float64) {
			for _, item := range items {
				sum += item
			}
		}
		return sum
	default:
		return 0.0
	}
}

func SUM2(input1 interface{}, input2 interface{}) float64 {
	return SUM(input1) + SUM(input2)
}

func SUMn(inputs ...interface{}) float64 {
	sum := 0.0
	for _, input := range inputs {
		sum += SUM(input)
	}
	return sum
}

func POWER(num interface{}, power interface{}) float64 {
	switch num.(type) {
	case float64:
		switch power.(type) {
		case float64:
			return math.Pow(num.(float64), power.(float64))
		case int:
			return math.Pow(num.(float64), float64(power.(int)))
		}
	case int:
		switch power.(type) {
		case float64:
			return math.Pow(float64(num.(int)), power.(float64))
		case int:
			return math.Pow(float64(num.(int)), float64(power.(int)))
		}
	}
	return 0.0
}
