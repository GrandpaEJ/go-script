# Go-Script VS Code Extension

This extension provides language support for Go-Script (.gos) files in Visual Studio Code.

## Features

- **Syntax Highlighting**: Full syntax highlighting for Go-Script language constructs
- **Code Snippets**: Predefined snippets for common Go-Script patterns
- **Auto-completion**: IntelliSense support for Go-Script keywords and constructs
- **Indentation**: Smart indentation based on Go-Script syntax rules
- **Commands**: Integrated commands to run, build, and debug Go-Script files
- **Bracket Matching**: Automatic bracket and quote pairing

## Installation

1. Install the Go-Script compiler (`gos`) on your system
2. Install this extension from the VS Code marketplace
3. Open any `.gos` file to activate the extension

## Commands

The extension provides the following commands:

- **Go-Script: Run** (`Ctrl+F5`): Compile and run the current Go-Script file
- **Go-Script: Build** (`Ctrl+Shift+B`): Compile the current Go-Script file to Go code
- **Go-Script: Debug**: Run the file with debug information

## Snippets

Type the following prefixes and press `Tab` to expand:

- `func` - Create a function
- `main` - Create main function
- `if` - Create if statement
- `ifelse` - Create if-else statement
- `ifelif` - Create if-elif-else statement
- `for` - Create for loop
- `while` - Create while loop
- `struct` - Create struct
- `method` - Create method
- `var` - Create variable assignment
- `print` - Print statement
- `import` - Import statement
- `from` - From import statement
- `err` - Error handling
- `switch` - Switch statement

## Configuration

The extension can be configured through VS Code settings:

- `go-script.gosPath`: Path to the gos executable (default: "gos")
- `go-script.enableAutoCompletion`: Enable auto-completion (default: true)
- `go-script.enableSyntaxHighlighting`: Enable syntax highlighting (default: true)

## Example

```gos
# Hello World in Go-Script
func main():
    print("Hello, World!")
    print("Welcome to Go-Script!")
```

## Requirements

- Go-Script compiler (`gos`) must be installed and available in PATH
- VS Code version 1.74.0 or higher

## Known Issues

- Complex multi-function files may have parsing limitations
- Some advanced Go features may not be fully supported yet

## Release Notes

### 0.1.0

- Initial release
- Basic syntax highlighting
- Code snippets
- Command integration
- Auto-completion support

## Contributing

Contributions are welcome! Please see the [main repository](https://github.com/GrandpaEJ/go-script) for contribution guidelines.

## License

MIT License - see the main repository for details.
