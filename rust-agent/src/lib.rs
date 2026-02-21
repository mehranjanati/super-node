use serde::{Deserialize, Serialize};
use std::ffi::{CStr, CString};
use std::mem;
use std::os::raw::{c_char, c_void};

// Redpanda Connect Wasm ABI (Simplified)
// We read input from memory, process it, and write output back.

#[derive(Serialize, Deserialize)]
struct MarketData {
    symbol: String,
    price: f64,
    volume_24h: f64,
}

#[derive(Serialize, Deserialize)]
struct AnalysisResult {
    strategy: String,
    top_pick: String,
    risk_level: String,
    reasoning: String,
}

#[no_mangle]
pub extern "C" fn process(ptr: *mut u8, len: usize) -> u64 {
    // 1. Read Input (JSON)
    let input_data = unsafe { std::slice::from_raw_parts(ptr, len) };
    let input_str = String::from_utf8_lossy(input_data);
    
    // Default Fallback
    let mut result = AnalysisResult {
        strategy: "HOLD".to_string(),
        top_pick: "UNKNOWN".to_string(),
        risk_level: "HIGH".to_string(),
        reasoning: "Invalid Input".to_string(),
    };

    if let Ok(data) = serde_json::from_str::<Vec<MarketData>>(&input_str) {
        // 2. Logic (The "AI" part - here rule-based, but could be ONNX inference)
        // In a real scenario, we would load a small ONNX model here.
        
        let mut best_coin = "BTC";
        let mut max_volume = 0.0;
        
        for coin in &data {
            if coin.volume_24h > max_volume {
                max_volume = coin.volume_24h;
                best_coin = &coin.symbol;
            }
        }
        
        // Simple Heuristic
        if max_volume > 1_000_000_000.0 {
            result.strategy = "BUY".to_string();
            result.risk_level = "MEDIUM".to_string();
            result.reasoning = format!("High volume detected on {}. Market is active.", best_coin);
        } else {
             result.strategy = "SELL".to_string();
             result.reasoning = "Volume is drying up.".to_string();
        }
        result.top_pick = best_coin.to_string();
    }

    // 3. Write Output
    let output_json = serde_json::to_string(&result).unwrap();
    let output_bytes = output_json.as_bytes();
    let output_len = output_bytes.len();
    
    // Allocate memory for output (Caller must free this if needed, 
    // but in Redpanda Connect Wasm ABI, usually we return a pointer/len pair packed in u64)
    // NOTE: This is a simplified ABI assumption. 
    // Real Redpanda Wasm modules use 'malloc' and specific exports.
    // For this demo, we assume the host handles memory via specific ABI calls not fully shown here.
    // We will stick to a basic JSON-in/JSON-out transformation compatible with standard Wasm processors.
    
    // To make it compile for wasm32-wasi/unknown:
    let ptr = output_bytes.as_ptr();
    
    // Return ptr << 32 | len
    ((ptr as u64) << 32) | (output_len as u64)
}

// Allocator for the host to write input data
#[no_mangle]
pub extern "C" fn allocate(size: usize) -> *mut c_void {
    let mut buffer = Vec::with_capacity(size);
    let ptr = buffer.as_mut_ptr();
    mem::forget(buffer);
    ptr as *mut c_void
}

#[no_mangle]
pub extern "C" fn deallocate(ptr: *mut c_void, capacity: usize) {
    unsafe {
        let _ = Vec::from_raw_parts(ptr, capacity, capacity);
    }
}
