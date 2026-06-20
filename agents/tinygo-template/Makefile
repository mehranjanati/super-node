.PHONY: build clean

# Build the WebAssembly binary using TinyGo targeting WASI
build:
	tinygo build -o agent.wasm -target=wasi main.go

# Clean the generated binaries
clean:
	rm -f agent.wasm
