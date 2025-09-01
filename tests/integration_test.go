package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHelloWorldIntegration(t *testing.T) {
	// Create a temporary .gos file
	content := `func main():
    print("Hello, World!")
    print("Go-Script is working!")`

	tempFile := createTempGosFile(t, "hello_test.gos", content)
	defer os.Remove(tempFile)

	// Build the gos binary
	buildGos(t)

	// Test compilation
	goFile := strings.TrimSuffix(tempFile, ".gos") + ".go"
	defer os.Remove(goFile)

	cmd := exec.Command("./gos", "build", tempFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile .gos file: %v\nOutput: %s", err, output)
	}

	// Check if Go file was created
	if _, err := os.Stat(goFile); os.IsNotExist(err) {
		t.Fatalf("Generated Go file does not exist: %s", goFile)
	}

	// Test running
	cmd = exec.Command("./gos", "run", tempFile)
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run .gos file: %v\nOutput: %s", err, output)
	}

	expectedOutput := "Hello, World!\nGo-Script is working!\n"
	if string(output) != expectedOutput {
		t.Fatalf("Unexpected output. Expected:\n%s\nGot:\n%s", expectedOutput, string(output))
	}
}

func TestCalculatorIntegration(t *testing.T) {
	content := `func add(a int, b int) int:
    return a + b

func multiply(a int, b int) int:
    return a * b

func main():
    x := 10
    y := 5
    
    sum := add(x, y)
    product := multiply(x, y)
    
    print("Sum:", sum)
    print("Product:", product)`

	tempFile := createTempGosFile(t, "calc_test.gos", content)
	defer os.Remove(tempFile)

	buildGos(t)

	cmd := exec.Command("./gos", "run", tempFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run calculator: %v\nOutput: %s", err, output)
	}

	expectedLines := []string{
		"Sum: 15",
		"Product: 50",
	}

	outputStr := string(output)
	for _, line := range expectedLines {
		if !strings.Contains(outputStr, line) {
			t.Fatalf("Expected output to contain '%s', but got:\n%s", line, outputStr)
		}
	}
}

func TestVariablesAndTypesIntegration(t *testing.T) {
	content := `func main():
    # Test different variable types
    name := "Go-Script"
    age := 25
    active := true
    score := 98.5
    
    print("Name:", name)
    print("Age:", age)
    print("Active:", active)
    print("Score:", score)
    
    # Test arrays
    numbers := [1, 2, 3, 4, 5]
    print("Numbers:", numbers[0], numbers[4])
    
    # Test maps
    person := {"name": "Alice", "age": 30}
    print("Person name:", person["name"])`

	tempFile := createTempGosFile(t, "vars_test.gos", content)
	defer os.Remove(tempFile)

	buildGos(t)

	cmd := exec.Command("./gos", "run", tempFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run variables test: %v\nOutput: %s", err, output)
	}

	expectedLines := []string{
		"Name: Go-Script",
		"Age: 25",
		"Active: true",
		"Score: 98.5",
	}

	outputStr := string(output)
	for _, line := range expectedLines {
		if !strings.Contains(outputStr, line) {
			t.Fatalf("Expected output to contain '%s', but got:\n%s", line, outputStr)
		}
	}
}

func TestControlFlowIntegration(t *testing.T) {
	content := `func main():
    # Test if-elif-else
    score := 85
    
    if score >= 90:
        print("Grade: A")
    elif score >= 80:
        print("Grade: B")
    elif score >= 70:
        print("Grade: C")
    else:
        print("Grade: F")
    
    # Test for loop with range
    print("Counting:")
    for i in range(5):
        print("Count:", i)
    
    # Test while loop
    count := 0
    while count < 3:
        print("While:", count)
        count = count + 1`

	tempFile := createTempGosFile(t, "control_test.gos", content)
	defer os.Remove(tempFile)

	buildGos(t)

	cmd := exec.Command("./gos", "run", tempFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run control flow test: %v\nOutput: %s", err, output)
	}

	expectedLines := []string{
		"Grade: B",
		"Counting:",
		"Count: 0",
		"Count: 4",
		"While: 0",
		"While: 2",
	}

	outputStr := string(output)
	for _, line := range expectedLines {
		if !strings.Contains(outputStr, line) {
			t.Fatalf("Expected output to contain '%s', but got:\n%s", line, outputStr)
		}
	}
}

func TestErrorHandlingIntegration(t *testing.T) {
	// Test syntax error
	content := `func main(
    print("Missing closing parenthesis")`

	tempFile := createTempGosFile(t, "error_test.gos", content)
	defer os.Remove(tempFile)

	buildGos(t)

	cmd := exec.Command("./gos", "build", tempFile)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected compilation to fail, but it succeeded. Output: %s", output)
	}

	if !strings.Contains(string(output), "Compilation error") {
		t.Fatalf("Expected compilation error message, got: %s", output)
	}
}

func TestComplexProgramIntegration(t *testing.T) {
	content := `struct Person:
    name string
    age int
    
    func greet(self) string:
        return "Hello, I'm " + self.name

func fibonacci(n int) int:
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

func main():
    # Test struct
    person := Person{name: "Alice", age: 30}
    print(person.greet())
    
    # Test recursive function
    print("Fibonacci of 7:", fibonacci(7))
    
    # Test array operations
    numbers := [1, 2, 3, 4, 5]
    sum := 0
    for i in range(len(numbers)):
        sum = sum + numbers[i]
    print("Sum of array:", sum)`

	tempFile := createTempGosFile(t, "complex_test.gos", content)
	defer os.Remove(tempFile)

	buildGos(t)

	cmd := exec.Command("./gos", "run", tempFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run complex program: %v\nOutput: %s", err, output)
	}

	expectedLines := []string{
		"Hello, I'm Alice",
		"Fibonacci of 7: 13",
	}

	outputStr := string(output)
	for _, line := range expectedLines {
		if !strings.Contains(outputStr, line) {
			t.Fatalf("Expected output to contain '%s', but got:\n%s", line, outputStr)
		}
	}
}

func TestCLICommands(t *testing.T) {
	buildGos(t)

	// Test version command
	cmd := exec.Command("./gos", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	if !strings.Contains(string(output), "Go-Script v") {
		t.Fatalf("Version command output unexpected: %s", output)
	}

	// Test help command
	cmd = exec.Command("./gos", "help")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	if !strings.Contains(string(output), "Usage:") {
		t.Fatalf("Help command output unexpected: %s", output)
	}

	// Test invalid command
	cmd = exec.Command("./gos", "invalid")
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Expected invalid command to fail")
	}

	if !strings.Contains(string(output), "unknown command") {
		t.Fatalf("Invalid command output unexpected: %s", output)
	}
}

// Helper functions

func createTempGosFile(t *testing.T, filename, content string) string {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, filename)

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return tempFile
}

func buildGos(t *testing.T) {
	// Check if gos binary already exists
	if _, err := os.Stat("./gos"); err == nil {
		return
	}

	cmd := exec.Command("go", "build", "-o", "gos", "../cmd/gos")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build gos binary: %v\nOutput: %s", err, output)
	}
}
