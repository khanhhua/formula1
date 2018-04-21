package funs

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	CompareEqual   int = 0
	CompareGreater int = 1
	CompareLesser  int = -1
)

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

func IFERROR(input1 interface{}, input2 interface{}) interface{} {
	switch input1.(type) {
	case error:
		return input2
	default:
		return input1
	}
}

// OR2 Evaluate to a boolean
func OR3(input1 interface{}, input2 interface{}, input3 interface{}) bool {
	return OR(input1) || OR(input2) || OR(input3)
}

// AND Evaluate to a boolean
func AND(input interface{}) bool {
	return boolean(input)
}

// AND2 Evaluate to a boolean
func AND2(input1 interface{}, input2 interface{}) bool {
	return boolean(input1) && boolean(input2)
}

func AND3(input1 interface{}, input2 interface{}, input3 interface{}) bool {
	return boolean(input1) && boolean(input2) && boolean(input3)
}

// FLOOR Floor function
func FLOOR(input interface{}) float64 {
	return math.Floor(input.(float64))
}

func ROUND(input interface{}, precision float64) float64 {
	pow := math.Pow(10, precision)
	return math.Round(input.(float64)*pow) / pow
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
	case []interface{}:
		sum := 0.0
		for _, item := range input.([]interface{}) {
			switch item.(type) {
			case float64:
				sum += item.(float64)
				break
			default:
				break
			}
		}
		return sum
	case [][]interface{}:
		sum := 0.0
		var outer [][]interface{}
		var inner []interface{}

		outer = input.([][]interface{})
		for i := 0; i < len(outer); i++ {
			inner = outer[i]
			for j := 0; j < len(inner); j++ {
				sum += inner[j].(float64)
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

// MATCH
// MATCH is an Excel function used to locate the position of a lookup value in a
// row, column, or table. MATCH supports approximate and exact matching, and
// wildcards (* ?) for partial matches. Often, the INDEX function is combined
// with MATCH to retrieve the value at the position returned by MATCH.
func MATCH(value interface{}, lookupRange interface{}, matchType int) int {
	switch lookupRange.(type) {
	case []interface{}:
		outer := lookupRange.([]interface{})
		if matchType == 1 {
			for i := 0; i < len(outer); i++ {
				if compare(value, outer[i]) == CompareGreater {
					return i + 1
				}
			}
		} else if matchType == -1 {
			for i := 0; i < len(outer); i++ {
				if compare(value, outer[i]) == CompareLesser {
					return i + 1
				}
			}
		} else {
			for i := 0; i < len(outer); i++ {
				if compare(value, outer[i]) == CompareEqual {
					return i + 1
				}
			}
		}
	}

	return 0
}

func compare(a, b interface{}) int {
	var sA, sB string

	switch a.(type) {
	case string:
		sA = a.(string)
		break
	case int:
		sA = strconv.FormatInt(int64(a.(int)), 10)
		break
	case float64:
		sA = strconv.FormatFloat(a.(float64), 4, 9, 64)
		break
	}

	switch b.(type) {
	case string:
		sB = b.(string)
		break
	case int:
		sB = strconv.FormatInt(int64(b.(int)), 10)
		break
	case float64:
		sB = strconv.FormatFloat(b.(float64), 4, 9, 64)
		break
	}

	return strings.Compare(sA, sB)
}

// VLOOKUP
// Use VLOOKUP, one of the lookup and reference functions, when you need to find
// things in a table or a range by row. For example, look up a price of an
// automotive part by the part number.
//
// In its simplest form, the VLOOKUP function says:
//
// =VLOOKUP(Value you want to look up,
// 					range where you want to lookup the value,
//					the column number in the range containing the return value,
// 					Exact Match or Approximate Match – indicated as 0/FALSE or 1/TRUE).
// @see https://support.office.com/en-us/article/VLOOKUP-function-0bbc8083-26fe-4963-8ab8-93a18ad188a1
func VLOOKUP(value interface{}, lookupRange interface{}, index int, approx bool) interface{} {
	if index < 1 {
		return errors.New("Index cannot be less than 1")
	}
	nativeIndex := index - 1

	switch lookupRange.(type) {
	case [][]interface{}:
		var outer [][]interface{}
		var inner []interface{}
		outer = lookupRange.([][]interface{})
		for i := len(outer) - 1; i >= 0; i-- {
			inner = outer[i]

			var referenceValue interface{}

			referenceValue = inner[0]
			switch referenceValue.(type) {
			case float64:
				if approx {
					if result, ok := value.(float64); ok && result >= referenceValue.(float64) {
						return inner[nativeIndex]
					}
				} else {
					if result, ok := value.(float64); ok && result == referenceValue.(float64) {
						return inner[nativeIndex]
					}
				}
			case string:
				if approx {
					if result, ok := value.(string); ok && result >= referenceValue.(string) {
						return inner[nativeIndex]
					}
				} else {
					if result, ok := value.(string); ok && result == referenceValue.(string) {
						return inner[nativeIndex]
					}
				}
			}
		}

		return errors.New("N/A")
	default:
		return errors.New("Lookup range must be 2D Slice")
	}
}

func COUNTIF(lookupRange interface{}, referenceValue interface{}) (count float64) {
	count = 0

	switch lookupRange.(type) {
	case [][]interface{}:
		var outer [][]interface{}
		var inner []interface{}
		outer = lookupRange.([][]interface{})
		for i := 0; i < len(outer); i++ {
			inner = outer[i]

			for j := 0; j < len(inner); j++ {
				value := inner[j]
				switch referenceValue.(type) {
				case float64:
					if result, ok := value.(float64); ok && result == referenceValue.(float64) {
						count++
					}
				case string:
					if result, ok := value.(string); ok && result == referenceValue.(string) {
						count++
					}
				}
			}
		}

		return count
	case []interface{}:
		inner := lookupRange.([]interface{})

		for j := 0; j < len(inner); j++ {
			value := inner[j]
			switch referenceValue.(type) {
			case float64:
				if result, ok := value.(float64); ok && result == referenceValue.(float64) {
					count++
				}
			case string:
				if result, ok := value.(string); ok && result == referenceValue.(string) {
					count++
				}
			}
		}

		return count
	default:
		return 0.0
	}
}
