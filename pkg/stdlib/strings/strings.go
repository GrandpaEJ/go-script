package strings

import (
	"fmt"
	"strings"
)

// StringFunctions contains string manipulation functions available in Go-Script
var StringFunctions = map[string]func(args ...interface{}) interface{}{
	"upper":     Upper,
	"lower":     Lower,
	"title":     Title,
	"strip":     Strip,
	"lstrip":    LStrip,
	"rstrip":    RStrip,
	"split":     Split,
	"join":      Join,
	"replace":   Replace,
	"contains":  Contains,
	"startswith": StartsWith,
	"endswith":  EndsWith,
	"find":      Find,
	"count":     Count,
}

// Upper converts string to uppercase
func Upper(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("upper() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.ToUpper(s)
}

// Lower converts string to lowercase
func Lower(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("lower() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.ToLower(s)
}

// Title converts string to title case
func Title(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("title() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.Title(s)
}

// Strip removes whitespace from both ends
func Strip(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("strip() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.TrimSpace(s)
}

// LStrip removes whitespace from the left end
func LStrip(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("lstrip() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.TrimLeft(s, " \t\n\r")
}

// RStrip removes whitespace from the right end
func RStrip(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("rstrip() takes exactly one argument")
	}
	
	s := toString(args[0])
	return strings.TrimRight(s, " \t\n\r")
}

// Split splits a string by separator
func Split(args ...interface{}) interface{} {
	if len(args) < 1 || len(args) > 2 {
		panic("split() takes 1 or 2 arguments")
	}
	
	s := toString(args[0])
	
	if len(args) == 1 {
		// Split by whitespace
		return strings.Fields(s)
	}
	
	sep := toString(args[1])
	return strings.Split(s, sep)
}

// Join joins strings with separator
func Join(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("join() takes exactly two arguments")
	}
	
	sep := toString(args[0])
	
	// Convert slice to string slice
	switch slice := args[1].(type) {
	case []string:
		return strings.Join(slice, sep)
	case []interface{}:
		strSlice := make([]string, len(slice))
		for i, v := range slice {
			strSlice[i] = toString(v)
		}
		return strings.Join(strSlice, sep)
	default:
		panic("join() second argument must be a slice")
	}
}

// Replace replaces occurrences of old with new
func Replace(args ...interface{}) interface{} {
	if len(args) < 3 || len(args) > 4 {
		panic("replace() takes 3 or 4 arguments")
	}
	
	s := toString(args[0])
	old := toString(args[1])
	new := toString(args[2])
	
	n := -1 // replace all by default
	if len(args) == 4 {
		n = toInt(args[3])
	}
	
	return strings.Replace(s, old, new, n)
}

// Contains checks if string contains substring
func Contains(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("contains() takes exactly two arguments")
	}
	
	s := toString(args[0])
	substr := toString(args[1])
	return strings.Contains(s, substr)
}

// StartsWith checks if string starts with prefix
func StartsWith(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("startswith() takes exactly two arguments")
	}
	
	s := toString(args[0])
	prefix := toString(args[1])
	return strings.HasPrefix(s, prefix)
}

// EndsWith checks if string ends with suffix
func EndsWith(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("endswith() takes exactly two arguments")
	}
	
	s := toString(args[0])
	suffix := toString(args[1])
	return strings.HasSuffix(s, suffix)
}

// Find finds the index of substring
func Find(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("find() takes exactly two arguments")
	}
	
	s := toString(args[0])
	substr := toString(args[1])
	return strings.Index(s, substr)
}

// Count counts occurrences of substring
func Count(args ...interface{}) interface{} {
	if len(args) != 2 {
		panic("count() takes exactly two arguments")
	}
	
	s := toString(args[0])
	substr := toString(args[1])
	return strings.Count(s, substr)
}

// Helper functions

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		panic("argument must be a number")
	}
}
