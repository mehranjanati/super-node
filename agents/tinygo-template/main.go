package main

import (
	"fmt"
	"os"
)

// main is required for TinyGo to compile to WASI
func main() {
	// You can read arguments passed by Wazero
	args := os.Args
	if len(args) > 1 {
		fmt.Printf("Hello from TinyGo Agent! Argument: %s\n", args[1])
	} else {
		fmt.Println("Hello from TinyGo Agent! No arguments provided.")
	}

	// Complex agent logic goes here...
}

//export processData
func processData() int32 {
	// Example of an exported function that Wazero can call directly
	return 42
}
