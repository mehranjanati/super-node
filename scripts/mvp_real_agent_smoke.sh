#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PORTAL_DIR="${ROOT_DIR}/portal1"
BFF_DIR="${ROOT_DIR}/BFF"
BFF_URL="${BFF_URL:-http://localhost:3001}"
GO_URL="${GO_URL:-http://localhost:3000}"
SKIP_NETWORK="${SKIP_NETWORK:-0}"
TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

print_section() {
  echo
  echo "== $1 =="
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

check_contains() {
  local file="$1"
  local needle="$2"
  local description="$3"

  if grep -q "${needle}" "${file}"; then
    echo "PASS: ${description}"
  else
    echo "FAIL: ${description}" >&2
    echo "Expected to find: ${needle}" >&2
    echo "Captured output:" >&2
    cat "${file}" >&2
    exit 1
  fi
}

healthcheck() {
  local url="$1"
  curl -fsS --max-time 5 "${url}" >/dev/null 2>&1
}

require_cmd curl
require_cmd python3
require_cmd npm
require_cmd bunx

print_section "Portal Check"
(cd "${PORTAL_DIR}" && npm run check)

print_section "BFF Type Check"
(cd "${BFF_DIR}" && bunx tsc --noEmit index.ts --lib es2022,dom --module esnext --moduleResolution bundler --target es2022 --types bun)

if [[ "${SKIP_NETWORK}" == "1" ]]; then
  print_section "Network Smoke"
  echo "SKIP_NETWORK=1 -> skipping live BFF/Go validation."
  exit 0
fi

print_section "Service Health"
if ! healthcheck "${BFF_URL}/api/health"; then
  echo "BFF healthcheck failed at ${BFF_URL}/api/health" >&2
  echo "Tip: start BFF first, e.g. (cd ${BFF_DIR} && GO_SERVER_URL=${GO_URL} bun run index.ts)" >&2
  exit 1
fi
if ! healthcheck "${GO_URL}/health"; then
  echo "Go gateway healthcheck failed at ${GO_URL}/health" >&2
  echo "Tip: start the Go gateway before running this smoke script." >&2
  exit 1
fi
curl -fsS "${BFF_URL}/api/health"
echo
curl -fsS "${GO_URL}/health"
echo

print_section "Main Smoke Request"
MAIN_OUTPUT="${TMP_DIR}/workflow_insight_main.txt"
curl -sS -N -X POST "${BFF_URL}/api/chat" \
  -H 'Content-Type: application/json' \
  -d '{
    "messages":[
      {
        "id":"msg-1",
        "role":"user",
        "parts":[
          {
            "type":"text",
            "text":"آخرین workflowها را خلاصه کن"
          }
        ]
      }
    ],
    "data":{
      "currentPath":"global_chat",
      "currentRoute":"#/dashboard",
      "selectedAgent":{
        "id":"agent-smoke",
        "name":"workflow-insight-agent",
        "type":"analytics",
        "config":{
          "capability":"workflow_insight",
          "executionMode":"read_only_workflow_insight",
          "resultSurface":"global_chat",
          "tools":["workflow_insight"],
          "systemPrompt":"Summarize workflow state for the user"
        }
      }
    }
  }' | tee "${MAIN_OUTPUT}"
echo
check_contains "${MAIN_OUTPUT}" "workflow_insight" "tool invocation should target workflow_insight"
check_contains "${MAIN_OUTPUT}" "summary" "response should include summary"
check_contains "${MAIN_OUTPUT}" "selected_workflows" "response should include selected_workflows"

print_section "Failure Sanity: Missing Workflow"
MISSING_WORKFLOW_OUTPUT="${TMP_DIR}/workflow_not_found.txt"
curl -sS -N -X POST "${BFF_URL}/api/chat" \
  -H 'Content-Type: application/json' \
  -d '{
    "messages":[
      {
        "id":"msg-2",
        "role":"user",
        "parts":[
          {
            "type":"text",
            "text":"وضعیت workflow workflow-does-not-exist را بگو"
          }
        ]
      }
    ],
    "data":{
      "currentPath":"global_chat",
      "currentRoute":"#/dashboard",
      "selectedAgent":{
        "id":"agent-smoke",
        "name":"workflow-insight-agent",
        "type":"analytics",
        "config":{
          "capability":"workflow_insight",
          "executionMode":"read_only_workflow_insight",
          "resultSurface":"global_chat",
          "tools":["workflow_insight"],
          "systemPrompt":"Summarize workflow state for the user"
        }
      }
    }
  }' | tee "${MISSING_WORKFLOW_OUTPUT}"
echo
check_contains "${MISSING_WORKFLOW_OUTPUT}" "workflow_not_found" "missing workflow should produce workflow_not_found"

print_section "Summary"
echo "Smoke validation passed."
echo "Manual UI checks still recommended:"
echo "- Foundry -> Save Draft"
echo "- Projects -> Reopen Draft"
echo "- GlobalChat -> Active Draft visibility"
echo "- GlobalChat -> Result card readability"
