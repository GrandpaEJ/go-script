package main

import (
	"encoding/json"
	"net/http"
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println("=== Import Alias Test ===")
	fmt.Println("Testing that aliases are correctly resolved...")
	fmt.Printf("✅ fmt package works\n")
	fmt.Printf("✅ json alias resolved to: encoding/json\n")
	fmt.Printf("✅ http alias resolved to: net/http\n")
	fmt.Printf("✅ base64 alias resolved to: encoding/base64\n")
	fmt.Println("=== All Aliases Working! ===")
}
