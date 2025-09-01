package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/GrandpaEJ/go-script/pkg/codegen"
	"github.com/GrandpaEJ/go-script/pkg/lexer"
	"github.com/GrandpaEJ/go-script/pkg/parser"
	"github.com/GrandpaEJ/go-script/pkg/stdlib"
)

const (
	version = "0.1.0"
	usage   = `Go-Script Compiler v%s

Usage:
    gos <command> [arguments]

Commands:
    run <file>              Compile and run a .gos file
    build <file>            Compile a .gos file to Go code
    build -o <file>         Compile and create binary executable
    build -go <file>        Compile to Go code (same as build)
    debug <file>            Compile and run with debug information

    # Package Management
    init                    Initialize a new Go-Script project
    mod init <name>         Initialize a new module
    mod tidy                Clean up module dependencies
    mod download            Download module dependencies

    # Module Commands
    install <module>        Install a Go-Script module
    uninstall <module>      Uninstall a Go-Script module
    list                    List installed modules
    search <query>          Search for available modules

    # Standard Library
    stdlib                  Show available import aliases

    version                 Show version information
    help                    Show this help message

Examples:
    gos run hello.gos
    gos build main.gos
    gos build -o myapp main.gos
    gos debug main.gos
    gos init
    gos mod init myproject
    gos install math-utils
    gos list
`
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			printError("run command requires a file argument")
			printUsage()
			os.Exit(1)
		}
		runFile(os.Args[2])
	case "build":
		if len(os.Args) < 3 {
			printError("build command requires a file argument")
			printUsage()
			os.Exit(1)
		}
		// Handle build flags
		if len(os.Args) >= 4 && os.Args[2] == "-o" {
			// gos build -o output file.gos
			if len(os.Args) < 5 {
				printError("build -o requires output name and file argument")
				printUsage()
				os.Exit(1)
			}
			buildBinary(os.Args[4], os.Args[3])
		} else if len(os.Args) >= 4 && os.Args[2] == "-go" {
			// gos build -go file.gos
			buildFile(os.Args[3])
		} else {
			// gos build file.gos
			buildFile(os.Args[2])
		}
	case "debug":
		if len(os.Args) < 3 {
			printError("debug command requires a file argument")
			printUsage()
			os.Exit(1)
		}
		debugFile(os.Args[2])
	case "init":
		initProject()
	case "mod":
		if len(os.Args) < 3 {
			printError("mod command requires a subcommand")
			printUsage()
			os.Exit(1)
		}
		handleModCommand(os.Args[2:])
	case "install":
		if len(os.Args) < 3 {
			printError("install command requires a module name")
			printUsage()
			os.Exit(1)
		}
		installModule(os.Args[2])
	case "uninstall":
		if len(os.Args) < 3 {
			printError("uninstall command requires a module name")
			printUsage()
			os.Exit(1)
		}
		uninstallModule(os.Args[2])
	case "list":
		listModules()
	case "search":
		if len(os.Args) < 3 {
			printError("search command requires a query")
			printUsage()
			os.Exit(1)
		}
		searchModules(os.Args[2])
	case "stdlib":
		showStdlibAliases()
	case "version":
		fmt.Printf("%sGo-Script v%s%s\n", ColorBold, version, ColorReset)
	case "help":
		printUsage()
	default:
		printError(fmt.Sprintf("unknown command '%s'", command))
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(usage, version)
}

// Enhanced error formatting functions
func printError(message string) {
	fmt.Printf("%s%sError:%s %s\n", ColorBold, ColorRed, ColorReset, message)
}

func printSuccess(message string) {
	fmt.Printf("%s%sSuccess:%s %s\n", ColorBold, ColorGreen, ColorReset, message)
}

func printWarning(message string) {
	fmt.Printf("%s%sWarning:%s %s\n", ColorBold, ColorYellow, ColorReset, message)
}

func printInfo(message string) {
	fmt.Printf("%s%sInfo:%s %s\n", ColorBold, ColorBlue, ColorReset, message)
}

func printCompilationError(filename string, errors []string) {
	fmt.Printf("%s%sCompilation failed:%s %s%s%s\n", ColorBold, ColorRed, ColorReset, ColorCyan, filename, ColorReset)
	fmt.Println()

	for i, err := range errors {
		fmt.Printf("%s%d.%s %s\n", ColorYellow, i+1, ColorReset, err)
	}

	fmt.Println()
	fmt.Printf("%sHint:%s Check your syntax, especially indentation and colons after function definitions.\n", ColorBlue, ColorReset)
}

func measureExecutionTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

func runFile(filename string) {
	// Check if file exists and has .gos extension
	if !strings.HasSuffix(filename, ".gos") {
		printError("file must have .gos extension")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		printError(fmt.Sprintf("file '%s' does not exist", filename))
		os.Exit(1)
	}

	// Compile to Go code with timing
	var goCode string
	var err error
	compileTime := measureExecutionTime(func() {
		goCode, err = compileFile(filename)
	})

	if err != nil {
		if strings.Contains(err.Error(), "Parsing errors:") {
			// Extract error list from error message
			errorMsg := err.Error()
			lines := strings.Split(errorMsg, "\n")
			var errors []string
			for _, line := range lines[1:] { // Skip first line "Parsing errors:"
				if strings.TrimSpace(line) != "" {
					errors = append(errors, strings.TrimPrefix(strings.TrimSpace(line), "- "))
				}
			}
			printCompilationError(filename, errors)
		} else {
			printError(fmt.Sprintf("compilation failed: %v", err))
		}
		os.Exit(1)
	}

	// Create temporary directory for generated Go code
	tempDir, err := os.MkdirTemp("", "gos-*")
	if err != nil {
		printError(fmt.Sprintf("creating temp directory: %v", err))
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Write Go code to temporary file
	goFile := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(goFile, []byte(goCode), 0644)
	if err != nil {
		printError(fmt.Sprintf("writing Go code: %v", err))
		os.Exit(1)
	}

	// Initialize Go module in temp directory
	cmd := exec.Command("go", "mod", "init", "temp")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		printError(fmt.Sprintf("initializing Go module: %v", err))
		os.Exit(1)
	}

	// Run the Go code with timing
	fmt.Printf("%sCompiled in:%s %v\n", ColorGreen, ColorReset, compileTime)
	fmt.Printf("%sRunning:%s %s%s%s\n", ColorBlue, ColorReset, ColorCyan, filename, ColorReset)
	fmt.Println()

	var execTime time.Duration
	execTime = measureExecutionTime(func() {
		cmd = exec.Command("go", "run", "main.go")
		cmd.Dir = tempDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			printError(fmt.Sprintf("runtime error: %v", err))
			os.Exit(1)
		}
	})

	fmt.Println()
	fmt.Printf("%sExecution completed in:%s %v\n", ColorGreen, ColorReset, execTime)
}

func buildFile(filename string) {
	// Check if file exists and has .gos extension
	if !strings.HasSuffix(filename, ".gos") {
		printError("file must have .gos extension")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		printError(fmt.Sprintf("file '%s' does not exist", filename))
		os.Exit(1)
	}

	// Compile to Go code with timing
	var goCode string
	var err error
	compileTime := measureExecutionTime(func() {
		goCode, err = compileFile(filename)
	})

	if err != nil {
		if strings.Contains(err.Error(), "Parsing errors:") {
			// Extract error list from error message
			errorMsg := err.Error()
			lines := strings.Split(errorMsg, "\n")
			var errors []string
			for _, line := range lines[1:] { // Skip first line "Parsing errors:"
				if strings.TrimSpace(line) != "" {
					errors = append(errors, strings.TrimPrefix(strings.TrimSpace(line), "- "))
				}
			}
			printCompilationError(filename, errors)
		} else {
			printError(fmt.Sprintf("compilation failed: %v", err))
		}
		os.Exit(1)
	}

	// Generate output filename
	baseName := strings.TrimSuffix(filename, ".gos")
	outputFile := baseName + ".go"

	// Write Go code to file
	err = os.WriteFile(outputFile, []byte(goCode), 0644)
	if err != nil {
		printError(fmt.Sprintf("writing output file: %v", err))
		os.Exit(1)
	}

	printSuccess(fmt.Sprintf("compiled '%s' to '%s' in %v", filename, outputFile, compileTime))
}

func buildBinary(filename, outputName string) {
	// Check if file exists and has .gos extension
	if !strings.HasSuffix(filename, ".gos") {
		printError("file must have .gos extension")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		printError(fmt.Sprintf("file '%s' does not exist", filename))
		os.Exit(1)
	}

	// Compile to Go code with timing
	var goCode string
	var err error
	compileTime := measureExecutionTime(func() {
		goCode, err = compileFile(filename)
	})

	if err != nil {
		if strings.Contains(err.Error(), "Parsing errors:") {
			// Extract error list from error message
			errorMsg := err.Error()
			lines := strings.Split(errorMsg, "\n")
			var errors []string
			for _, line := range lines[1:] { // Skip first line "Parsing errors:"
				if strings.TrimSpace(line) != "" {
					errors = append(errors, strings.TrimPrefix(strings.TrimSpace(line), "- "))
				}
			}
			printCompilationError(filename, errors)
		} else {
			printError(fmt.Sprintf("compilation failed: %v", err))
		}
		os.Exit(1)
	}

	// Create temporary directory for generated Go code
	tempDir, err := os.MkdirTemp("", "gos-*")
	if err != nil {
		printError(fmt.Sprintf("creating temp directory: %v", err))
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	// Write Go code to temporary file
	goFile := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(goFile, []byte(goCode), 0644)
	if err != nil {
		printError(fmt.Sprintf("writing Go code: %v", err))
		os.Exit(1)
	}

	// Initialize Go module in temp directory
	cmd := exec.Command("go", "mod", "init", "temp")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		printError(fmt.Sprintf("initializing Go module: %v", err))
		os.Exit(1)
	}

	// Build binary
	fmt.Printf("%sCompiled in:%s %v\n", ColorGreen, ColorReset, compileTime)
	fmt.Printf("%sBuilding binary:%s %s%s%s\n", ColorBlue, ColorReset, ColorCyan, outputName, ColorReset)

	// Get absolute path for output
	currentDir, _ := os.Getwd()
	outputPath := filepath.Join(currentDir, outputName)

	var buildTime time.Duration
	buildTime = measureExecutionTime(func() {
		cmd = exec.Command("go", "build", "-o", outputPath, "main.go")
		cmd.Dir = tempDir
		if err := cmd.Run(); err != nil {
			printError(fmt.Sprintf("building binary: %v", err))
			os.Exit(1)
		}
	})

	printSuccess(fmt.Sprintf("built binary '%s' in %v (total: %v)", outputName, buildTime, compileTime+buildTime))
}

func debugFile(filename string) {
	// Check if file exists and has .gos extension
	if !strings.HasSuffix(filename, ".gos") {
		printError("file must have .gos extension")
		os.Exit(1)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		printError(fmt.Sprintf("file '%s' does not exist", filename))
		os.Exit(1)
	}

	fmt.Printf("%sDebug Mode:%s %s%s%s\n", ColorYellow, ColorReset, ColorCyan, filename, ColorReset)
	fmt.Println()

	// Show lexer tokens
	content, err := os.ReadFile(filename)
	if err != nil {
		printError(fmt.Sprintf("reading file: %v", err))
		os.Exit(1)
	}

	fmt.Printf("%sLexer Tokens:%s\n", ColorBlue, ColorReset)
	l := lexer.New(string(content))
	tokenCount := 0
	for {
		tok := l.NextToken()
		if tok.Type == lexer.EOF {
			break
		}
		if tok.Type != lexer.COMMENT && tok.Type != lexer.NEWLINE {
			fmt.Printf("  %s%d.%s %s%s%s: %q\n", ColorYellow, tokenCount+1, ColorReset,
				ColorGreen, lexer.TokenTypeString(tok.Type), ColorReset, tok.Literal)
			tokenCount++
		}
	}
	fmt.Printf("Total tokens: %d\n\n", tokenCount)

	// Compile and run with debug info
	runFile(filename)
}

func compileFile(filename string) (string, error) {
	// Read the source file
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Create lexer
	l := lexer.New(string(content))

	// Create parser
	p := parser.New(l)

	// Parse the program
	program := p.ParseProgram()

	// Check for parsing errors
	if errors := p.Errors(); len(errors) > 0 {
		var errorMsg strings.Builder
		errorMsg.WriteString("Parsing errors:\n")
		for _, err := range errors {
			errorMsg.WriteString(fmt.Sprintf("  - %s\n", err))
		}
		return "", fmt.Errorf(errorMsg.String())
	}

	// Generate Go code
	generator := codegen.New()
	goCode := generator.Generate(program)

	// Add necessary imports if they're used
	goCode = addRequiredImports(goCode)

	return goCode, nil
}

func addRequiredImports(code string) string {
	var imports []string

	// Always add fmt for built-in functions (print, printf, etc.)
	imports = append(imports, `"fmt"`)

	// Add other imports based on usage
	if strings.Contains(code, "bufio.") {
		imports = append(imports, `"bufio"`)
	}
	if strings.Contains(code, "os.") {
		imports = append(imports, `"os"`)
	}
	if strings.Contains(code, "time.") {
		imports = append(imports, `"time"`)
	}
	if strings.Contains(code, "strings.") {
		imports = append(imports, `"strings"`)
	}
	if strings.Contains(code, "strconv.") {
		imports = append(imports, `"strconv"`)
	}
	if strings.Contains(code, "reflect.") {
		imports = append(imports, `"reflect"`)
	}

	// Check if imports already exist to avoid duplicates
	hasExistingImports := strings.Contains(code, "import (") || strings.Contains(code, `import "`)

	// Only add automatic imports if there are no existing imports
	if !hasExistingImports {
		// Check for math usage (power operator)
		if strings.Contains(code, "math.Pow") {
			imports = append(imports, `"math"`)
		}
	}

	// If we have imports to add, insert them
	if len(imports) > 0 {
		lines := strings.Split(code, "\n")
		var result []string

		// Find package declaration
		packageFound := false
		importSectionAdded := false

		for _, line := range lines {
			result = append(result, line)

			if strings.HasPrefix(line, "package ") {
				packageFound = true
			}

			// Add imports after package declaration and empty line
			if packageFound && !importSectionAdded && line == "" {
				result = append(result, "import (")
				for _, imp := range imports {
					result = append(result, "\t"+imp)
				}
				result = append(result, ")")
				result = append(result, "")
				importSectionAdded = true
			}
		}

		// If no empty line was found after package, add imports at the end
		if packageFound && !importSectionAdded {
			result = append(result, "")
			result = append(result, "import (")
			for _, imp := range imports {
				result = append(result, "\t"+imp)
			}
			result = append(result, ")")
			result = append(result, "")
		}

		return strings.Join(result, "\n")
	}

	return code
}

// Package management functions

func initProject() {
	printInfo("Initializing new Go-Script project...")

	// Create gos.mod file
	modContent := `module main

go 1.21

# Go-Script module configuration
gos_version "1.0.0"

# Dependencies
require (
    # Standard Go modules work automatically
)

# Go-Script specific configuration
config {
    default_package "main"
    output_dir "./generated"
    module_paths ["./modules", "./lib"]
}
`

	err := os.WriteFile("gos.mod", []byte(modContent), 0644)
	if err != nil {
		printError(fmt.Sprintf("Failed to create gos.mod: %v", err))
		return
	}

	// Create basic directory structure
	dirs := []string{"modules", "lib", "examples", "tests"}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			printWarning(fmt.Sprintf("Failed to create directory %s: %v", dir, err))
		}
	}

	// Create example main.gos
	mainContent := `# Welcome to Go-Script!

func main():
    print("Hello, Go-Script!")
    print("Edit this file to get started.")
`

	err = os.WriteFile("main.gos", []byte(mainContent), 0644)
	if err != nil {
		printWarning(fmt.Sprintf("Failed to create main.gos: %v", err))
	}

	printSuccess("Project initialized successfully!")
	printInfo("Created: gos.mod, main.gos, and project directories")
	printInfo("Run 'gos run main.gos' to test your setup")
}

func handleModCommand(args []string) {
	if len(args) == 0 {
		printError("mod command requires a subcommand")
		return
	}

	switch args[0] {
	case "init":
		if len(args) < 2 {
			printError("mod init requires a module name")
			return
		}
		initModule(args[1])
	case "tidy":
		tidyModule()
	case "download":
		downloadDependencies()
	default:
		printError(fmt.Sprintf("unknown mod subcommand '%s'", args[0]))
	}
}

func initModule(name string) {
	printInfo(fmt.Sprintf("Initializing module '%s'...", name))

	modContent := fmt.Sprintf(`module %s

go 1.21

gos_version "1.0.0"

config {
    default_package "main"
    output_dir "./generated"
}
`, name)

	err := os.WriteFile("gos.mod", []byte(modContent), 0644)
	if err != nil {
		printError(fmt.Sprintf("Failed to create gos.mod: %v", err))
		return
	}

	printSuccess(fmt.Sprintf("Module '%s' initialized successfully!", name))
}

func tidyModule() {
	printInfo("Tidying module dependencies...")
	// TODO: Implement dependency cleanup
	printSuccess("Module dependencies tidied")
}

func downloadDependencies() {
	printInfo("Downloading module dependencies...")
	// TODO: Implement dependency download
	printSuccess("Dependencies downloaded")
}

func installModule(moduleName string) {
	printInfo(fmt.Sprintf("Installing module '%s'...", moduleName))

	// Create modules directory if it doesn't exist
	err := os.MkdirAll("modules", 0755)
	if err != nil {
		printError(fmt.Sprintf("Failed to create modules directory: %v", err))
		return
	}

	// TODO: Implement actual module installation from registry
	// For now, create a placeholder
	moduleDir := filepath.Join("modules", moduleName)
	err = os.MkdirAll(moduleDir, 0755)
	if err != nil {
		printError(fmt.Sprintf("Failed to create module directory: %v", err))
		return
	}

	// Create a basic module file
	moduleContent := fmt.Sprintf(`# %s module
# This is a placeholder module

func %s_function():
    print("Function from %s module")
`, moduleName, moduleName, moduleName)

	moduleFile := filepath.Join(moduleDir, moduleName+".gos")
	err = os.WriteFile(moduleFile, []byte(moduleContent), 0644)
	if err != nil {
		printError(fmt.Sprintf("Failed to create module file: %v", err))
		return
	}

	printSuccess(fmt.Sprintf("Module '%s' installed successfully!", moduleName))
}

func uninstallModule(moduleName string) {
	printInfo(fmt.Sprintf("Uninstalling module '%s'...", moduleName))

	moduleDir := filepath.Join("modules", moduleName)
	err := os.RemoveAll(moduleDir)
	if err != nil {
		printError(fmt.Sprintf("Failed to uninstall module: %v", err))
		return
	}

	printSuccess(fmt.Sprintf("Module '%s' uninstalled successfully!", moduleName))
}

func listModules() {
	printInfo("Installed modules:")

	modulesDir := "modules"
	if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
		printInfo("No modules directory found. Run 'gos init' to create project structure.")
		return
	}

	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		printError(fmt.Sprintf("Failed to read modules directory: %v", err))
		return
	}

	if len(entries) == 0 {
		printInfo("No modules installed.")
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  %s%s%s\n", ColorGreen, entry.Name(), ColorReset)
		}
	}
}

func searchModules(query string) {
	printInfo(fmt.Sprintf("Searching for modules matching '%s'...", query))

	// TODO: Implement actual module registry search
	// For now, show some example modules
	exampleModules := []string{
		"math-utils - Mathematical utility functions",
		"string-helpers - String manipulation helpers",
		"file-utils - File system utilities",
		"http-client - HTTP client library",
		"json-parser - JSON parsing utilities",
	}

	printInfo("Available modules:")
	for _, module := range exampleModules {
		if strings.Contains(strings.ToLower(module), strings.ToLower(query)) {
			fmt.Printf("  %s%s%s\n", ColorCyan, module, ColorReset)
		}
	}

	printInfo("Use 'gos install <module-name>' to install a module")
}

func showStdlibAliases() {
	printInfo("Go-Script Standard Library Import Aliases")
	fmt.Println()

	categories := stdlib.GetCommonAliases()

	for category, aliases := range categories {
		fmt.Printf("%s%s:%s\n", ColorBold+ColorBlue, category, ColorReset)
		for _, alias := range aliases {
			parts := strings.Split(alias, " - ")
			if len(parts) == 2 {
				aliasName := parts[0]
				description := parts[1]
				realPath := stdlib.GetRealPackagePath(aliasName)
				fmt.Printf("  %s%-12s%s -> %s%-20s%s %s%s%s\n",
					ColorGreen, aliasName, ColorReset,
					ColorCyan, realPath, ColorReset,
					ColorYellow, description, ColorReset)
			}
		}
		fmt.Println()
	}

	fmt.Printf("%sUsage Examples:%s\n", ColorBold+ColorPurple, ColorReset)
	fmt.Printf("  %simport \"json\"%s     # Imports encoding/json\n", ColorGreen, ColorReset)
	fmt.Printf("  %simport \"http\"%s     # Imports net/http\n", ColorGreen, ColorReset)
	fmt.Printf("  %simport \"fs\"%s       # Imports io/fs\n", ColorGreen, ColorReset)
	fmt.Printf("  %simport \"crypto\"%s   # Imports crypto\n", ColorGreen, ColorReset)
	fmt.Println()

	printInfo("Total aliases available: " + fmt.Sprintf("%d", len(stdlib.ImportAliases)))
}
