package math

import (
	"math"
)

// MathFunctions contains mathematical functions available in Go-Script
var MathFunctions = map[string]func(args ...interface{}) interface{}{
	"abs":   Abs,
	"ceil":  Ceil,
	"floor": Floor,
	"round": Round,
	"sqrt":  Sqrt,
	"pow":   Pow,
	"sin":   Sin,
	"cos":   Cos,
	"tan":   Tan,
	"log":   Log,
	"log10": Log10,
	"exp":   Exp,
	"min":   Min,
	"max":   Max,
}

// Constants
const (
	Pi = math.Pi
	E  = math.E
)

// Abs returns the absolute value
func Abs(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("abs() takes exactly one argument")
	}
	
	switch v := args[0].(type) {
	case int:
		if v < 0 {
			return -v
		}
		return v
	case int64:
		return math.Abs(float64(v))
	case float64:
		return math.Abs(v)
	default:
		panic("abs() argument must be a number")
	}
}

// Ceil returns the ceiling of a number
func Ceil(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("ceil() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Ceil(f)
}

// Floor returns the floor of a number
func Floor(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("floor() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Floor(f)
}

// Round rounds a number to the nearest integer
func Round(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("round() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Round(f)
}

// Sqrt returns the square root
func Sqrt(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("sqrt() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Sqrt(f)
}

// Pow returns x raised to the power of y
func Pow(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("pow() takes exactly two arguments")
	}
	
	x := toFloat64(args[0])
	y := toFloat64(args[1])
	return math.Pow(x, y)
}

// Sin returns the sine of x
func Sin(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("sin() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Sin(f)
}

// Cos returns the cosine of x
func Cos(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("cos() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Cos(f)
}

// Tan returns the tangent of x
func Tan(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("tan() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Tan(f)
}

// Log returns the natural logarithm of x
func Log(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("log() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Log(f)
}

// Log10 returns the base-10 logarithm of x
func Log10(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("log10() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Log10(f)
}

// Exp returns e raised to the power of x
func Exp(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("exp() takes exactly one argument")
	}
	
	f := toFloat64(args[0])
	return math.Exp(f)
}

// Min returns the minimum of the arguments
func Min(args ...interface{}) interface{} {
	if len(args) == 0 {
		panic("min() takes at least one argument")
	}
	
	min := toFloat64(args[0])
	for i := 1; i < len(args); i++ {
		val := toFloat64(args[i])
		if val < min {
			min = val
		}
	}
	return min
}

// Max returns the maximum of the arguments
func Max(args ...interface{}) interface{} {
	if len(args) == 0 {
		panic("max() takes at least one argument")
	}
	
	max := toFloat64(args[0])
	for i := 1; i < len(args); i++ {
		val := toFloat64(args[i])
		if val > max {
			max = val
		}
	}
	return max
}

// Helper function to convert to float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	default:
		panic("argument must be a number")
	}
}
