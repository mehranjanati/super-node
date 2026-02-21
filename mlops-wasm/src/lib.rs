use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct Event {
    #[serde(rename = "type")]
    event_type: String,
    content: String,
    feedback: Option<Feedback>,
    tool_call: Option<ToolCall>,
    agent_response: Option<String>,
    correction: Option<String>,
}

#[derive(Serialize, Deserialize, Debug)]
struct Feedback {
    score: Option<f64>,
}

#[derive(Serialize, Deserialize, Debug)]
struct ToolCall {
    valid: bool,
}

#[derive(Serialize, Deserialize, Debug)]
struct TrainingExample {
    instruction: String,
    input: String,
    output: String,
    metadata: Metadata,
}

#[derive(Serialize, Deserialize, Debug)]
struct Metadata {
    reward_score: f64,
    source: String,
}

#[no_mangle]
pub extern "C" fn process_event(ptr: *mut u8, len: usize) -> *mut u8 {
    // 1. Read input JSON from memory
    let input_slice = unsafe { std::slice::from_raw_parts(ptr, len) };
    let event: Event = match serde_json::from_slice(input_slice) {
        Ok(e) => e,
        Err(_) => return std::ptr::null_mut(), // Return null on error
    };

    // 2. Calculate Reward Score (The "Brain" logic)
    let base_score = 0.5;
    
    // Factor 1: User Feedback
    let feedback_score = match &event.feedback {
        Some(f) => match f.score {
            Some(s) if s >= 4.0 => 0.5,
            Some(s) if s <= 2.0 => -0.5,
            _ => 0.0,
        },
        None => 0.0,
    };

    // Factor 2: Syntax Validity
    let syntax_score = match (&event.event_type.as_str(), &event.tool_call) {
        (&"tool_call", Some(t)) if t.valid => 0.3,
        (&"tool_call", _) => -0.5,
        _ => 0.1,
    };

    // Factor 3: Reasoning Depth
    let reasoning_score = if event.content.contains("<thought>") && event.content.len() > 200 {
        0.2
    } else {
        0.0
    };

    let total_score = base_score + feedback_score + syntax_score + reasoning_score;

    // 3. Filter Low Quality
    if total_score < 0.7 {
        return std::ptr::null_mut(); // Drop this event
    }

    // 4. Transform to Training Format
    let instruction = match event.event_type.as_str() {
        "user_query" => "Answer the following user query with reasoning.",
        "tool_call" => "Generate a valid JSON tool call for the given context.",
        _ => "Process the following event.",
    }.to_string();

    let output = event.correction.clone().or(event.agent_response.clone()).unwrap_or_default();

    let example = TrainingExample {
        instruction,
        input: event.content.clone(),
        output,
        metadata: Metadata {
            reward_score: total_score,
            source: "wasm_processor".to_string(),
        },
    };

    // 5. Serialize output to JSON
    let output_json = serde_json::to_vec(&example).unwrap();
    
    // Return pointer to new memory (simplified for example, needs memory management in real Wasm)
    // In production, we'd use wit-bindgen or similar for robust ABI
    let mut boxed = output_json.into_boxed_slice();
    let ptr = boxed.as_mut_ptr();
    std::mem::forget(boxed);
    ptr
}
