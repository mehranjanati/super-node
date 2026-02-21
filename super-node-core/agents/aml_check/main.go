package main

import (
	"strings"
	"unsafe"
)

// TinyGo export for the host to call
//export aml_check
func aml_check(ptr uint32, size uint32) uint64

// Mocking the "host" function to request GPU resources
//export request_gpu_acceleration
func request_gpu_acceleration()

// main is required for TinyGo, even if empty
func main() {}

//export perform_analysis
func perform_analysis(walletPtr uint32, walletSize uint32, txnPtr uint32, txnSize uint32) uint64 {
	// 1. Read input strings from host memory
	wallet := getString(walletPtr, walletSize)
	// txnID := getString(txnPtr, txnSize) // Not used in simple logic, but available

	// 2. Request Internal GPU (Simulation of resource_tag: "GPU_REQUIRED_INTERNAL")
	request_gpu_acceleration()

	// 3. Perform Heavy Crypto/AML Logic (Mocked)
	isSuspicious := false

	// Simple rule: wallets starting with "0xDEAD" are flagged
	if strings.HasPrefix(wallet, "0xDEAD") {
		isSuspicious = true
	}

	// 4. Return result (1 for safe, 0 for suspicious/rejected)
	if isSuspicious {
		return 0
	}
	return 1
}

// Helper to read string from WebAssembly linear memory
func getString(ptr uint32, size uint32) string {
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)
	return string(bytes)
}

// Helper to allocate memory for the host to write strings into
//export allocate
func allocate(size uint32) uint32 {
	buf := make([]byte, size)
	return uint32(uintptr(unsafe.Pointer(&buf[0])))
}
