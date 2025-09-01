# Getting Started with Go-Script

Welcome to Go-Script! This tutorial will help you write and run your first Go-Script program.

## Prerequisites

- Go 1.22 or later installed on your system
- Basic familiarity with programming concepts
- A text editor or IDE

## Installation

### Option 1: Install from Source

```bash
# Clone the repository
git clone https://github.com/GrandpaEJ/go-script.git
cd go-script

# Build the compiler
go build -o gos ./cmd/gos

# Add to PATH (optional)
sudo mv gos /usr/local/bin/
```

### Option 2: Install via Go

```bash
go install github.com/GrandpaEJ/go-script/cmd/gos@latest
```

## Your First Program

Let's create a simple "Hello, World!" program.

### Step 1: Create a File

Create a new file called `hello.gos`:

```gos
# This is a comment
func main():
    print("Hello, World!")
    print("Welcome to Go-Script!")
```

### Step 2: Run the Program

```bash
gos run hello.gos
```

You should see:
```
Hello, World!
Welcome to Go-Script!
```

### Step 3: Compile to Go

You can also compile your Go-Script code to Go:

```bash
gos build hello.gos
```

This creates `hello.go`:

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("Welcome to Go-Script!")
}
```

## Understanding the Syntax

Let's break down the hello world program:

```gos
# This is a comment
func main():
    print("Hello, World!")
    print("Welcome to Go-Script!")
```

- `# This is a comment` - Comments start with `#`
- `func main():` - Function definition with colon
- Indented blocks - Use spaces/tabs for code blocks (like Python)
- `print()` - Built-in function for output

## Basic Program Structure

Every Go-Script program follows this structure:

```gos
# Optional comments and imports

func main():
    # Your code here
    print("Program starts here")
```

## Key Features

### 1. Clean Syntax
```gos
# Variables
name := "Go-Script"
version := 1

# Functions
func greet():
    print("Hello from Go-Script!")
```

### 2. Go Integration
```gos
import "fmt"
import "time"

func main():
    fmt.Println("Using Go's fmt package")
    time.Sleep(1000000000) # 1 second
```

### 3. Type Safety
Go-Script inherits Go's type system while providing a more convenient syntax.

## Common Commands

### Run a Program
```bash
gos run program.gos
```

### Compile to Go
```bash
gos build program.gos
```

### Check Version
```bash
gos version
```

### Get Help
```bash
gos help
```

## Next Steps

Now that you have Go-Script running, explore these topics:

1. [Basic Syntax](basic-syntax.md) - Learn the language fundamentals
2. [Functions and Variables](functions-variables.md) - Core programming concepts
3. [Go Integration](go-integration.md) - Using Go packages

## Troubleshooting

### Common Issues

**Problem**: `gos: command not found`
**Solution**: Make sure Go-Script is installed and in your PATH.

**Problem**: Compilation errors
**Solution**: Check your syntax, especially indentation and colons.

**Problem**: Import errors
**Solution**: Ensure Go modules are properly configured.

### Getting Help

- Check the [FAQ](../guides/faq.md)
- Review [error messages](../guides/error-handling.md)
- Ask questions in our [community](https://github.com/GrandpaEJ/go-script/discussions)

## Example Programs

Here are some simple programs to try:

### Variables and Math
```gos
func main():
    x := 10
    y := 5
    print("Sum:", x + y)
    print("Product:", x * y)
```

### Conditionals
```gos
func main():
    age := 25
    if age >= 18:
        print("Adult")
    else:
        print("Minor")
```

### Multiple Functions
```gos
func greet():
    print("Hello!")

func main():
    greet()
    print("Program finished")
```

Congratulations! You've written your first Go-Script programs. Continue with the [Basic Syntax](basic-syntax.md) tutorial to learn more.
