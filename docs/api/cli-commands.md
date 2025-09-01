# CLI Commands Reference

The `gos` command-line tool provides several commands for working with Go-Script programs.

## Command Overview

```
gos <command> [arguments]
```

## Commands

### `run`

Compile and execute a Go-Script program.

**Syntax:**
```bash
gos run <file.gos>
```

**Example:**
```bash
gos run hello.gos
```

**Description:**
- Compiles the `.gos` file to Go code
- Creates a temporary Go module
- Executes the program
- Cleans up temporary files

**Options:**
- Currently no additional options supported

### `build`

Compile a Go-Script program to Go code.

**Syntax:**
```bash
gos build <file.gos>
```

**Example:**
```bash
gos build hello.gos
# Creates hello.go
```

**Description:**
- Transpiles `.gos` file to `.go` file
- Adds necessary imports automatically
- Preserves original filename with `.go` extension

**Output:**
- Creates a `.go` file in the same directory
- Includes proper Go package declaration
- Adds required imports (fmt, math, etc.)

### `version`

Display the Go-Script version.

**Syntax:**
```bash
gos version
```

**Example:**
```bash
$ gos version
Go-Script v0.1.0
```

**Description:**
Shows the current version of the Go-Script compiler.

### `help`

Display help information.

**Syntax:**
```bash
gos help
```

**Example:**
```bash
$ gos help
Go-Script Compiler v0.1.0

Usage:
    gos <command> [arguments]

Commands:
    run <file>      Compile and run a .gos file
    build <file>    Compile a .gos file to Go code
    version         Show version information
    help            Show this help message
```

## Exit Codes

The `gos` command uses standard exit codes:

- `0` - Success
- `1` - General error (compilation, runtime, etc.)
- `2` - Invalid command line arguments

## Error Handling

### Compilation Errors

When compilation fails, `gos` provides detailed error messages:

```bash
$ gos build invalid.gos
Compilation error: Parsing errors:
  - expected next token to be COLON, got IDENT instead
  - no prefix parse function for COMMA found
```

### Runtime Errors

When execution fails, `gos` shows the Go runtime error:

```bash
$ gos run crash.gos
Runtime error: panic: runtime error: index out of range
```

### File Errors

When file operations fail:

```bash
$ gos run nonexistent.gos
Error: file 'nonexistent.gos' does not exist
```

## Environment Variables

Currently, Go-Script doesn't use specific environment variables, but it respects standard Go environment variables:

- `GOPATH` - Go workspace path
- `GOROOT` - Go installation path
- `GO111MODULE` - Go modules mode

## Temporary Files

The `run` command creates temporary files:

- Temporary directory: `/tmp/gos-*`
- Generated Go files: `main.go`
- Go module files: `go.mod`

These are automatically cleaned up after execution.

## Integration with Go Tools

Generated Go code is compatible with standard Go tools:

```bash
# After gos build hello.gos
go run hello.go
go build hello.go
go fmt hello.go
```

## Future Commands (Planned)

The following commands are planned for future releases:

### `init`
Initialize a new Go-Script project:
```bash
gos init myproject
```

### `test`
Run Go-Script tests:
```bash
gos test
```

### `fmt`
Format Go-Script code:
```bash
gos fmt file.gos
```

### `doc`
Generate documentation:
```bash
gos doc
```

### `install`
Install Go-Script packages:
```bash
gos install package
```

## Examples

### Basic Usage
```bash
# Create and run a simple program
echo 'func main(): print("Hello!")' > hello.gos
gos run hello.gos
```

### Build and Run Separately
```bash
# Compile to Go
gos build program.gos

# Run with Go
go run program.go
```

### Check Compilation
```bash
# Just check if it compiles
gos build program.gos && echo "Compilation successful"
```

### Batch Processing
```bash
# Compile multiple files
for file in *.gos; do
    echo "Compiling $file..."
    gos build "$file"
done
```

## Troubleshooting

### Command Not Found
```bash
$ gos version
bash: gos: command not found
```

**Solution:** Ensure Go-Script is installed and in your PATH.

### Permission Denied
```bash
$ gos run program.gos
Runtime error: permission denied
```

**Solution:** Check file permissions and execution rights.

### Module Errors
```bash
$ gos run program.gos
Error initializing Go module: go: cannot find main module
```

**Solution:** Ensure you have proper Go installation and module support.

For more help, see the [troubleshooting guide](../guides/troubleshooting.md) or [file an issue](https://github.com/GrandpaEJ/go-script/issues).
