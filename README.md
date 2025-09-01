# Go-Script

A high-level scripting language with intuitive syntax that seamlessly integrates with Go's ecosystem.

## Features

- **High-level syntax**: Clean, readable code that's easy to write and maintain
- **Seamless Go integration**: Import and use any Go module or package
- **Native Go packaging**: Full compatibility with Go's module system
- **Custom packaging**: Enhanced import system for .gos to .gos dependencies
- **Easy functions**: Simplified function definitions and calls
- **Minimal code**: Express complex logic with fewer lines
- **Professional codebase**: Production-ready architecture and tooling

## Quick Start

### Installation

```bash
go install github.com/GrandpaEJ/go-script/cmd/gos@latest
```

### Hello World

Create a file `hello.gos`:

```gos
func main():
    print("Hello, World!")
```

Run it:

```bash
gos run hello.gos
```

## Language Features

### Variables and Types

```gos
# Variables with type inference
name := "Go-Script"
count := 42
active := true

# Explicit typing
var message string = "Hello"
var numbers []int = [1, 2, 3, 4, 5]
```

### Functions

```gos
# Simple function
func greet(name string) string:
    return "Hello, " + name

# Function with multiple returns
func divide(a, b float64) (float64, error):
    if b == 0:
        return 0, error("division by zero")
    return a / b, nil
```

### Control Flow

```gos
# If statements
if age >= 18:
    print("Adult")
elif age >= 13:
    print("Teenager")
else:
    print("Child")

# For loops
for i in range(10):
    print(i)

for key, value in items:
    print(key, ":", value)
```

### Go Integration

```gos
# Import Go packages
import "fmt"
import "net/http"
import "encoding/json"

# Use Go functions directly
func main():
    resp, err := http.Get("https://api.github.com")
    if err != nil:
        fmt.Println("Error:", err)
        return
    
    defer resp.Body.Close()
    fmt.Println("Status:", resp.Status)
```

### Custom Types and Structs

```gos
# Define structs
struct Person:
    name string
    age int
    
    func greet(self) string:
        return "Hi, I'm " + self.name

# Create and use
person := Person{name: "Alice", age: 30}
print(person.greet())
```

## Project Structure

```
go-script/
├── cmd/gos/           # CLI tool
├── pkg/
│   ├── lexer/         # Lexical analysis
│   ├── parser/        # Syntax analysis
│   ├── ast/           # Abstract Syntax Tree
│   ├── codegen/       # Go code generation
│   ├── runtime/       # Runtime support
│   └── stdlib/        # Standard library
├── examples/          # Example programs
├── docs/             # Documentation
└── tests/            # Test suite
```

## Development Status

This project is under active development. Core features are being implemented in phases:

1. ✅ Project setup and architecture
2. 🔄 Lexer and parser implementation
3. ⏳ Code generation and Go integration
4. ⏳ Runtime and standard library
5. ⏳ CLI tools and build system
6. ⏳ Testing and documentation

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
