# Go-Script Language Specification

## Overview

Go-Script (.gos) is a high-level scripting language that transpiles to Go code, providing an intuitive syntax while maintaining full compatibility with Go's ecosystem.

## Syntax Design Principles

1. **Indentation-based blocks**: Like Python/GDScript, use indentation instead of braces
2. **Type inference**: Automatic type detection where possible
3. **Simplified syntax**: Reduce boilerplate while maintaining Go's power
4. **Seamless Go integration**: Direct access to Go packages and modules

## Lexical Elements

### Comments

```gos
# Single line comment
# Multi-line comments use multiple # lines
```

### Identifiers

- Start with letter or underscore
- Followed by letters, digits, or underscores
- Case-sensitive

### Keywords

```
and, or, not, if, elif, else, for, while, func, return, import, from
struct, interface, var, const, true, false, nil, in, range, break, continue
defer, go, chan, select, case, default, switch, type, package
```

### Operators

```
Arithmetic: +, -, *, /, %, **
Comparison: ==, !=, <, <=, >, >=
Logical: and, or, not
Assignment: =, :=, +=, -=, *=, /=, %=
Bitwise: &, |, ^, <<, >>, &^
Other: ., ->, <-, ++, --
```

## Data Types

### Basic Types

```gos
# Integers
var age int = 25
count := 42

# Floats
var price float64 = 19.99
ratio := 0.75

# Strings
var name string = "Go-Script"
message := "Hello, World!"

# Booleans
var active bool = true
ready := false
```

### Collections

```gos
# Arrays
numbers := [1, 2, 3, 4, 5]
var scores [10]int

# Slices
items := []string{"apple", "banana", "cherry"}

# Maps
ages := {"Alice": 30, "Bob": 25, "Charlie": 35}
var config map[string]interface{}
```

### Pointers

```gos
var ptr *int
value := 42
ptr = &value
print(*ptr)  # Dereference
```

## Control Structures

### Conditional Statements

```gos
# If-elif-else
if score >= 90:
    grade = "A"
elif score >= 80:
    grade = "B"
elif score >= 70:
    grade = "C"
else:
    grade = "F"

# Short if
if err := doSomething(); err != nil:
    return err
```

### Loops

```gos
# For loop with range
for i in range(10):
    print(i)

# For loop with slice/array
for index, value in items:
    print(index, value)

# For loop with map
for key, value in ages:
    print(key, "is", value, "years old")

# While loop
while condition:
    # do something
    if break_condition:
        break
    if continue_condition:
        continue

# Infinite loop
for:
    # do something
```

### Switch Statements

```gos
switch value:
    case 1:
        print("One")
    case 2, 3:
        print("Two or Three")
    default:
        print("Other")

# Type switch
switch v := interface{}(x):
    case int:
        print("Integer:", v)
    case string:
        print("String:", v)
    default:
        print("Unknown type")
```

## Functions

### Function Definition

```gos
# Simple function
func greet(name string) string:
    return "Hello, " + name

# Multiple parameters and returns
func divide(a, b float64) (float64, error):
    if b == 0:
        return 0, error("division by zero")
    return a / b, nil

# Named return values
func calculate(x, y int) (sum, product int):
    sum = x + y
    product = x * y
    return  # naked return

# Variadic function
func sum(numbers ...int) int:
    total := 0
    for num in numbers:
        total += num
    return total
```

### Function Types and Closures

```gos
# Function as variable
var operation func(int, int) int = func(a, b int) int:
    return a + b

# Closure
func makeCounter() func() int:
    count := 0
    return func() int:
        count++
        return count

counter := makeCounter()
print(counter())  # 1
print(counter())  # 2
```

## Structs and Methods

### Struct Definition

```gos
struct Person:
    name string
    age int
    email string

# Constructor-like function
func NewPerson(name string, age int) Person:
    return Person{
        name: name,
        age: age,
        email: name + "@example.com"
    }
```

### Methods

```gos
# Method with receiver
func greet(self Person) string:
    return "Hi, I'm " + self.name

# Method with pointer receiver
func birthday(self *Person):
    self.age++

# Usage
person := NewPerson("Alice", 30)
print(person.greet())
person.birthday()
```

## Interfaces

```gos
interface Writer:
    write(data []byte) (int, error)

interface Reader:
    read(data []byte) (int, error)

# Embedding interfaces
interface ReadWriter:
    Reader
    Writer
```

## Error Handling

```gos
# Traditional Go-style
result, err := riskyOperation()
if err != nil:
    return err

# Simplified error handling (syntactic sugar)
result := riskyOperation() or return err
```

## Concurrency

```gos
# Goroutines
go doSomething()
go func():
    print("Anonymous goroutine")

# Channels
ch := make(chan int)
go func():
    ch <- 42

value := <-ch

# Select statement
select:
    case msg := <-ch1:
        print("Received from ch1:", msg)
    case ch2 <- value:
        print("Sent to ch2")
    default:
        print("No communication")
```

## Package and Import System

### Package Declaration

```gos
package main

import "fmt"
import "net/http"
from "encoding/json" import Marshal, Unmarshal
```

### Go-Script to Go-Script Imports

```gos
# Import from another .gos file
from "./utils.gos" import helper_function
from "./models/user.gos" import User

# Import entire module
import "./math.gos" as math_utils
```

## Transpilation Rules

1. **Indentation to Braces**: Convert indented blocks to Go's brace syntax
2. **Type Inference**: Add explicit types where Go requires them
3. **Import Resolution**: Convert .gos imports to generated Go packages
4. **Syntax Sugar**: Expand simplified syntax to full Go equivalents
5. **Error Handling**: Convert simplified error handling to Go idioms

## File Structure

```
project/
├── main.gos           # Entry point
├── utils.gos          # Utility functions
├── models/
│   ├── user.gos       # User model
│   └── product.gos    # Product model
└── generated/         # Generated Go code (auto-created)
    ├── main.go
    ├── utils.go
    └── models/
        ├── user.go
        └── product.go
```
