package mcp

import "fmt"

// NewSampleToolBelt creates a new sample tool belt.
func NewSampleToolBelt() ToolBelt {
	return ToolBelt{
		Name: "sample",
		Tools: []Tool{
			{
				Name:        "hello",
				Description: "Prints a hello message.",
				Action: func(args ...interface{}) (interface{}, error) {
					if len(args) > 0 {
						return fmt.Sprintf("Hello, %s!", args[0]), nil
					}
					return "Hello, world!", nil
				},
			},
		},
	}
}
