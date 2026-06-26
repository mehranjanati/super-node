#!/usr/bin/env bash

set -euo pipefail

BFF_URL="${BFF_URL:-http://localhost:3001}"
GO_URL="${GO_URL:-http://localhost:3000}"
SUFFIX="$(date +%s)"
DIRECT_PROJECT="mvp-release-rc-${SUFFIX}"
CHAT_PROJECT="mvp-release-chat-${SUFFIX}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required command: $1" >&2
    exit 1
  fi
}

read_json_field() {
  local field="$1"
  python3 -c 'import json,sys; data=json.load(sys.stdin); value=data.get(sys.argv[1], ""); print(value if value is not None else "")' "$field"
}

print_section() {
  echo
  echo "== $1 =="
}

require_cmd curl
require_cmd python3

print_section "Health"
curl -fsS "${BFF_URL}/api/health"
echo
curl -fsS "${GO_URL}/internal/health"
echo

print_section "Direct Deploy"
DIRECT_DEPLOY_RESPONSE="$(
  curl -fsS -X POST "${GO_URL}/internal/tools/deploy" \
    -H 'Content-Type: application/json' \
    -d "{
      \"project_name\":\"${DIRECT_PROJECT}\",
      \"prompt\":\"Create a simple landing page for an AI studio with hero, pricing, and CTA sections.\",
      \"framework\":\"svelte\",
      \"theme\":\"modern\",
      \"template\":\"default\"
    }"
)"
echo "${DIRECT_DEPLOY_RESPONSE}"
DIRECT_WORKFLOW_ID="$(printf '%s' "${DIRECT_DEPLOY_RESPONSE}" | read_json_field workflow_id)"

print_section "Chat Trigger"
CHAT_RESPONSE="$(
  curl -fsS -N -X POST "${BFF_URL}/api/chat" \
    -H 'Content-Type: application/json' \
    -d "{
      \"messages\":[
        {
          \"id\":\"msg-1\",
          \"role\":\"user\",
          \"parts\":[
            {
              \"type\":\"text\",
              \"text\":\"برای من یک سایت ساده برای استودیو هوش مصنوعی بساز. اسم پروژه ${CHAT_PROJECT} باشد و با svelte اجرا شود.\"
            }
          ]
        }
      ],
      \"data\":{
        \"currentPath\":\"builder\",
        \"currentRoute\":\"#/builder\"
      }
    }"
)"
echo "${CHAT_RESPONSE}"
CHAT_WORKFLOW_ID="$(printf '%s' "${CHAT_RESPONSE}" | grep -o 'deploy-site-[^"]*' | tail -n 1 || true)"

print_section "Workflow Records"
if [[ -n "${DIRECT_WORKFLOW_ID}" ]]; then
  curl -fsS "${GO_URL}/workflows/${DIRECT_WORKFLOW_ID}"
  echo
fi

if [[ -n "${CHAT_WORKFLOW_ID}" ]]; then
  curl -fsS "${GO_URL}/workflows/${CHAT_WORKFLOW_ID}"
  echo
fi

print_section "Failure Sanity"
curl -sS -o /tmp/mvp_release_candidate_missing_fields.json -w 'invalid deploy => HTTP %{http_code}\n' \
  -X POST "${GO_URL}/internal/tools/deploy" \
  -H 'Content-Type: application/json' \
  -d '{}'
cat /tmp/mvp_release_candidate_missing_fields.json
echo

curl -sS -o /tmp/mvp_release_candidate_missing_workflow.json -w 'missing workflow => HTTP %{http_code}\n' \
  "${GO_URL}/workflows/workflow-does-not-exist"
cat /tmp/mvp_release_candidate_missing_workflow.json
echo

print_section "Summary"
echo "direct workflow_id: ${DIRECT_WORKFLOW_ID:-missing}"
echo "chat workflow_id: ${CHAT_WORKFLOW_ID:-missing}"
echo "note: release is GO only if a workflow is created, status is visible, and a meaningful artifact URL is returned."
