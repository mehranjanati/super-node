package workflow

import (
	"strings"
)

// SvelteTemplate is a simple template for SvelteKit + Tailwind frontend
const SvelteTemplate = `
<script>
	import { onMount } from 'svelte';
	let orderStatus = '';
	let loading = false;

	async function placeOrder() {
		loading = true;
		orderStatus = 'Processing...';
		
		const orderData = {
			order_id: Math.floor(Math.random() * 1000),
			item: "Super Node Laptop",
			timestamp: new Date().toISOString()
		};

		// In a real app, this goes to an API Gateway -> Redpanda
		// Here we simulate the structure needed for the Benthos/WASM pipeline
		try {
			const response = await fetch('/api/order', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(orderData)
			});
			const result = await response.json();
			orderStatus = 'Order Confirmed: ' + result.message;
		} catch (e) {
			orderStatus = 'Error: ' + e.message;
		} finally {
			loading = false;
		}
	}
</script>

<main class="container mx-auto p-4">
	<h1 class="text-3xl font-bold mb-4">Nexus Store</h1>
	<div class="card bg-base-100 shadow-xl">
		<div class="card-body">
			<h2 class="card-title">Super Laptop</h2>
			<p>Powered by WASM Backend Logic</p>
			<div class="card-actions justify-end">
				<button class="btn btn-primary" on:click={placeOrder} disabled={loading}>
					{loading ? 'Ordering...' : 'Buy Now'}
				</button>
			</div>
			{#if orderStatus}
				<div class="alert alert-info mt-4">
					{orderStatus}
				</div>
			{/if}
		</div>
	</div>
</main>
`

// WasmGoTemplate is the TinyGo-compatible backend logic
// It follows the Benthos WASM processor signature
const WasmGoTemplate = `package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// Required for TinyGo/WASM memory management with host
func main() {}

//export process
func process(ptr int32, size int32) int64 {
	// 1. Read Input (Order Data)
	inputBytes := readBuffer(ptr, size)
	
	var order map[string]interface{}
	if err := json.Unmarshal(inputBytes, &order); err != nil {
		return writeError("Invalid JSON")
	}

	// 2. Business Logic (The "Brain")
	// Example: Check inventory, validate user, etc.
	orderID := order["order_id"]
	item := order["item"]
	
	fmt.Printf("Processing order %v for %v\n", orderID, item)

	// 3. Create Response
	response := map[string]interface{}{
		"status": "success",
		"message": fmt.Sprintf("Processed order %v for %v via WASM", orderID, item),
		"processed_at_node": "SuperNode-Core-1",
	}

	outputBytes, _ := json.Marshal(response)
	
	// 4. Return Result to Benthos
	return writeBuffer(outputBytes)
}

// --- Helper Functions for WebAssembly Memory ---

//export allocate
func allocate(size uint32) uint32 {
	buf := make([]byte, size)
	return uint32(uintptr(unsafe.Pointer(&buf[0])))
}

func readBuffer(ptr int32, size int32) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)
}

func writeBuffer(buf []byte) int64 {
	ptr := uintptr(unsafe.Pointer(&buf[0]))
	size := len(buf)
	// Pack ptr and size into a single int64 for return
	return (int64(ptr) << 32) | int64(size)
}

func writeError(msg string) int64 {
	errBytes := []byte("{\"error\": \"" + msg + "\"}")
	return writeBuffer(errBytes)
}
`

const GithubCIWorkflow = `name: Build and Deploy

on:
  push:
    branches: [ "main", "feature/*" ]
  pull_request:
    branches: [ "main" ]

jobs:
  build-wasm:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install TinyGo
      uses: acifani/setup-tinygo@v1
      with:
        tinygo-version: '0.30.0'

    - name: Build WASM Backend
      run: |
        cd backend
        tinygo build -o ../artifact.wasm -target=wasi main.go

    - name: Upload WASM Artifact
      uses: actions/upload-artifact@v4
      with:
        name: wasm-bundle
        path: artifact.wasm

  build-frontend:
    runs-on: ubuntu-latest
    needs: build-wasm
    steps:
    - uses: actions/checkout@v4

    - name: Set up Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '20'

    - name: Build SvelteKit Frontend
      run: |
        cd frontend
        npm install
        npm run build

    - name: Upload Frontend Artifact
      uses: actions/upload-artifact@v4
      with:
        name: frontend-build
        path: frontend/.svelte-kit/output
`

// GenerateProjectStructure returns the file map for a full stack project
func GenerateProjectStructure(projectName string, schema string) map[string]string {
	files := make(map[string]string)

	// 0. CI/CD Pipeline
	files[".github/workflows/main.yml"] = GithubCIWorkflow

	// 1. Frontend (SvelteKit)
	files["frontend/package.json"] = `{
		"name": "` + projectName + `",
		"version": "0.0.1",
		"scripts": {
			"dev": "vite dev",
			"build": "vite build",
			"check": "svelte-kit sync && svelte-check --tsconfig ./tsconfig.json"
		},
		"devDependencies": {
			"@sveltejs/adapter-auto": "^3.0.0",
			"@sveltejs/kit": "^2.0.0",
			"@sveltejs/vite-plugin-svelte": "^3.0.0",
			"svelte": "^4.0.0",
			"vite": "^5.0.0",
			"tailwindcss": "^3.0.0",
			"daisyui": "^4.0.0"
		}
	}`

	files["frontend/svelte.config.js"] = `
		import adapter from '@sveltejs/adapter-auto';
		import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
		export default {
			preprocess: vitePreprocess(),
			kit: { adapter: adapter() }
		};
	`

	files["frontend/src/routes/+page.svelte"] = strings.Replace(SvelteTemplate, "Nexus Store", "Nexus Store - "+projectName, 1)
	files["frontend/src/app.html"] = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width" />
		<title>` + projectName + `</title>
	</head>
	<body data-theme="dark">
		<div style="display: contents">%sveltekit.body%</div>
	</body>
</html>`

	// 2. Backend (Go/WASM)
	files["backend/go.mod"] = `module ` + projectName + `_backend

go 1.21`

	files["backend/main.go"] = WasmGoTemplate

	// 3. Build Configuration (Makefile for CI)
	files["Makefile"] = `
build-wasm:
	cd backend && tinygo build -o ../artifact.wasm -target=wasi main.go

build-frontend:
	cd frontend && npm install && npm run build
`

	return files
}
