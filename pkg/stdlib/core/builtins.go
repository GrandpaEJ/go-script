package core

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// BuiltinFunction represents a built-in function
type BuiltinFunction struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

// Builtins contains all built-in functions available in Go-Script
var Builtins = map[string]*BuiltinFunction{
	"print":       {Name: "print", Fn: Print},
	"println":     {Name: "println", Fn: Println},
	"printf":      {Name: "printf", Fn: Printf},
	"input":       {Name: "input", Fn: Input},
	"len":         {Name: "len", Fn: Len},
	"range":       {Name: "range", Fn: Range},
	"str":         {Name: "str", Fn: Str},
	"int":         {Name: "int", Fn: Int},
	"float":       {Name: "float", Fn: Float},
	"bool":        {Name: "bool", Fn: Bool},
	"type":        {Name: "type", Fn: Type},
	"append":      {Name: "append", Fn: Append},
	"make":        {Name: "make", Fn: Make},
	"new":         {Name: "new", Fn: New},
	"now":         {Name: "now", Fn: Now},
	"format_time": {Name: "format_time", Fn: FormatTime},
}

// Print prints values separated by spaces
func Print(args ...interface{}) interface{} {
	fmt.Print(args...)
	return nil
}

// Println prints values separated by spaces with a newline
func Println(args ...interface{}) interface{} {
	fmt.Println(args...)
	return nil
}

// Printf prints formatted output (like Python's print with % formatting)
func Printf(args ...interface{}) interface{} {
	if len(args) == 0 {
		return nil
	}

	format, ok := args[0].(string)
	if !ok {
		// If first arg is not string, just print all args
		fmt.Print(args...)
		return nil
	}

	if len(args) == 1 {
		fmt.Print(format)
	} else {
		fmt.Printf(format, args[1:]...)
	}
	return nil
}

// Input reads a line from standard input (like Python's input())
func Input(args ...interface{}) interface{} {
	// Print prompt if provided
	if len(args) > 0 {
		fmt.Print(args[0])
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	// Remove trailing newline
	return strings.TrimSuffix(line, "\n")
}

// Len returns the length of a collection
func Len(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("len() takes exactly one argument")
	}

	arg := args[0]
	v := reflect.ValueOf(arg)

	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return v.Len()
	default:
		panic(fmt.Sprintf("object of type '%T' has no len()", arg))
	}
}

// Range generates a range of numbers
func Range(args ...interface{}) interface{} {
	switch len(args) {
	case 1:
		// range(n) -> 0 to n-1
		n := toInt(args[0])
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i
		}
		return result
	case 2:
		// range(start, stop) -> start to stop-1
		start := toInt(args[0])
		stop := toInt(args[1])
		result := make([]int, 0, stop-start)
		for i := start; i < stop; i++ {
			result = append(result, i)
		}
		return result
	case 3:
		// range(start, stop, step)
		start := toInt(args[0])
		stop := toInt(args[1])
		step := toInt(args[2])
		if step == 0 {
			panic("range() step argument must not be zero")
		}
		result := make([]int, 0)
		if step > 0 {
			for i := start; i < stop; i += step {
				result = append(result, i)
			}
		} else {
			for i := start; i > stop; i += step {
				result = append(result, i)
			}
		}
		return result
	default:
		panic("range() takes 1 to 3 arguments")
	}
}

// Str converts a value to string
func Str(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("str() takes exactly one argument")
	}
	return fmt.Sprintf("%v", args[0])
}

// Int converts a value to integer
func Int(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("int() takes exactly one argument")
	}
	return toInt(args[0])
}

// Float converts a value to float
func Float(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("float() takes exactly one argument")
	}
	return toFloat(args[0])
}

// Bool converts a value to boolean
func Bool(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("bool() takes exactly one argument")
	}
	return toBool(args[0])
}

// Type returns the type of a value
func Type(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("type() takes exactly one argument")
	}
	return reflect.TypeOf(args[0]).String()
}

// Append appends elements to a slice
func Append(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("append() takes at least 2 arguments")
	}

	slice := args[0]
	elements := args[1:]

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("first argument to append must be slice")
	}

	for _, elem := range elements {
		v = reflect.Append(v, reflect.ValueOf(elem))
	}

	return v.Interface()
}

// Make creates slices, maps, and channels
func Make(args ...interface{}) interface{} {
	if len(args) < 1 {
		panic("make() takes at least 1 argument")
	}

	// This is a simplified version - in a real implementation,
	// you'd need to handle type information properly
	typeStr := args[0].(string)

	switch typeStr {
	case "[]int":
		if len(args) > 1 {
			size := toInt(args[1])
			return make([]int, size)
		}
		return make([]int, 0)
	case "[]string":
		if len(args) > 1 {
			size := toInt(args[1])
			return make([]string, size)
		}
		return make([]string, 0)
	case "map[string]int":
		return make(map[string]int)
	default:
		panic(fmt.Sprintf("make: unsupported type %s", typeStr))
	}
}

// New allocates memory for a type
func New(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("new() takes exactly one argument")
	}

	// Simplified implementation
	typeStr := args[0].(string)
	switch typeStr {
	case "int":
		return new(int)
	case "string":
		return new(string)
	case "bool":
		return new(bool)
	default:
		panic(fmt.Sprintf("new: unsupported type %s", typeStr))
	}
}

// Helper functions

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		panic(fmt.Sprintf("invalid literal for int(): %s", val))
	default:
		panic(fmt.Sprintf("cannot convert %T to int", v))
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		panic(fmt.Sprintf("invalid literal for float(): %s", val))
	default:
		panic(fmt.Sprintf("cannot convert %T to float", v))
	}
}

func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case int64:
		return val != 0
	case float64:
		return val != 0.0
	case string:
		return val != ""
	default:
		return v != nil
	}
}

// Now returns the current time
func Now(args ...interface{}) interface{} {
	return time.Now()
}

// FormatTime formats time with human-readable format strings
func FormatTime(args ...interface{}) interface{} {
	if len(args) < 2 {
		panic("format_time() takes at least 2 arguments: time and format")
	}

	timeVal, ok := args[0].(time.Time)
	if !ok {
		panic("first argument to format_time must be a time value")
	}

	formatStr, ok := args[1].(string)
	if !ok {
		panic("second argument to format_time must be a format string")
	}

	// Convert human-readable format to Go's time format
	goFormat := convertTimeFormat(formatStr)
	return timeVal.Format(goFormat)
}

// convertTimeFormat converts human-readable time formats to Go's format
func convertTimeFormat(format string) string {
	// Map of human-readable formats to Go's reference time format
	replacements := map[string]string{
		"YYYY":    "2006",    // 4-digit year
		"YY":      "06",      // 2-digit year
		"MM":      "01",      // month with zero padding
		"M":       "1",       // month without zero padding
		"DD":      "02",      // day with zero padding
		"D":       "2",       // day without zero padding
		"HH":      "15",      // hour (24-hour) with zero padding
		"H":       "15",      // hour (24-hour) without zero padding
		"hh":      "03",      // hour (12-hour) with zero padding
		"h":       "3",       // hour (12-hour) without zero padding
		"mm":      "04",      // minute with zero padding
		"m":       "4",       // minute without zero padding
		"ss":      "05",      // second with zero padding
		"s":       "5",       // second without zero padding
		"AM":      "PM",      // AM/PM
		"am":      "pm",      // am/pm
		"Mon":     "Mon",     // abbreviated weekday
		"Monday":  "Monday",  // full weekday
		"Jan":     "Jan",     // abbreviated month
		"January": "January", // full month
	}

	result := format
	for human, goFmt := range replacements {
		result = strings.ReplaceAll(result, human, goFmt)
	}

	return result
}
