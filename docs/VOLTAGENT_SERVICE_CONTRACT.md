# VoltAgent Service Internal Contract

این سند contract رسمی و canonical برای تماس داخلی `Go -> voltagent-service` در فاز اول migration است.

هدف این contract:

- یک payload پایدار و versioned بین backend Go و `voltagent-service` تعریف کند
- مسئولیت reasoning/planning را از execution جدا کند
- فقط use case فاز اول یعنی `deploy_website` را پوشش دهد
- taxonomy خطاها را برای service team و Go team یکسان کند

## scope

این contract فقط برای مسیر داخلی زیر معتبر است:

- `Go -> voltagent-service`

و صراحتا خارج از scope است:

- تماس مستقیم `BFF -> voltagent-service`
- migration کامل chat path
- migration کامل manifest/tool discovery
- execution مستقیم durable داخل `voltagent-service`

## اصول contract

- contract باید `contract-first` و `versioned` باشد
- `voltagent-service` فقط planning/reasoning برمی‌گرداند
- execution نهایی use caseهای durable در `Go + Temporal` انجام می‌شود
- payloadها باید JSON ساده، deterministic و بدون ambiguity باشند
- فیلدهای undeclared نباید برای behavior اصلی لازم باشند

## versioning

نسخه canonical فاز اول:

- `contract_version = "v1alpha1"`

قواعد versioning:

- client باید این مقدار را در `POST /plan` ارسال کند
- service باید همان نسخه را در response بازگرداند
- هر breaking change باید با تغییر نسخه contract انجام شود
- اضافه شدن فیلد optional در همان نسخه مجاز است

## transport conventions

Headerهای پیشنهادی:

- `Content-Type: application/json`
- `Accept: application/json`
- `X-Request-Id: <uuid>` از سمت Go
- `X-Correlation-Id: <uuid>` در صورت وجود trace upstream
- `X-Caller-Service: go-gateway`

قواعد tracing:

- اگر `X-Request-Id` ارسال شود، service باید همان مقدار را در response body بازتاب دهد
- اگر headerها وجود نداشته باشند، service مجاز است `request_id` جدید تولید کند

## endpoint: `GET /health`

هدف:

- health check سبک برای compose, readiness probe و client preflight

### response: `200 OK`

```json
{
  "status": "ok",
  "service": "voltagent-service",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "checks": {
    "api": "ok",
    "planner": "ok"
  }
}
```

### response schema

- `status`: مقدار ثابت `ok`
- `service`: مقدار ثابت `voltagent-service`
- `contract_version`: نسخه active contract
- `request_id`: شناسه request برای tracing
- `checks.api`: سلامت لایه HTTP
- `checks.planner`: سلامت planner/model wiring

### response: `503 Service Unavailable`

اگر service بالا باشد اما برای planning آماده نباشد:

```json
{
  "status": "error",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "error": {
    "code": "VOLTAGENT_UNAVAILABLE",
    "message": "VoltAgent planner is not ready",
    "retryable": true,
    "details": {
      "check": "planner"
    }
  }
}
```

## endpoint: `POST /plan`

هدف:

- دریافت intent نرمال‌شده از Go
- تولید plan ساخت‌یافته برای execution در `Go + Temporal`

### request body

```json
{
  "contract_version": "v1alpha1",
  "intent": "deploy_website",
  "input": {
    "project_name": "demo-site",
    "prompt": "A minimal landing page for a crypto product",
    "framework": "svelte",
    "theme": "minimal"
  },
  "context": {
    "user_id": "usr_123",
    "session_id": "sess_456",
    "source": "go-gateway"
  }
}
```

### request schema

- `contract_version`: `string`, required, فعلا فقط `v1alpha1`
- `intent`: `string`, required, فعلا فقط `deploy_website`
- `input`: `object`, required
- `context`: `object`, optional اما strongly recommended

### request schema for `intent = deploy_website`

فیلدهای `input`:

- `project_name`: `string`, required
- `prompt`: `string`, required
- `framework`: `string`, optional, enum: `svelte`, `react`, `vue`
- `theme`: `string`, optional, enum: `light`, `dark`, `modern`, `minimal`

فیلدهای `context`:

- `user_id`: `string`, optional
- `session_id`: `string`, optional
- `source`: `string`, optional, مثال: `go-gateway`

### validation rules

- `contract_version` باید برابر `v1alpha1` باشد
- `intent` باید در فاز اول فقط `deploy_website` باشد
- `project_name` نباید empty string باشد
- `prompt` نباید empty string باشد
- اگر `framework` ارسال شود باید یکی از enumهای تعریف‌شده باشد
- اگر `theme` ارسال شود باید یکی از enumهای تعریف‌شده باشد

## success response

### response: `200 OK`

```json
{
  "status": "ok",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "plan": {
    "intent": "deploy_website",
    "kind": "workflow",
    "execution_target": "go-temporal",
    "workflow": {
      "name": "website-deployment-v1",
      "action": "start_dynamic_pipeline",
      "input": {
        "project_name": "demo-site",
        "prompt": "A minimal landing page for a crypto product",
        "framework": "svelte",
        "theme": "minimal"
      }
    },
    "artifacts": {
      "project_name": "demo-site"
    },
    "warnings": []
  }
}
```

### success response schema

- `status`: مقدار ثابت `ok`
- `contract_version`: نسخه contract
- `request_id`: شناسه tracing
- `plan.intent`: intent نهایی normalize‌شده
- `plan.kind`: در فاز اول مقدار ثابت `workflow`
- `plan.execution_target`: در فاز اول مقدار ثابت `go-temporal`
- `plan.workflow.name`: نام canonical workflow/pipeline
- `plan.workflow.action`: اکشن مورد انتظار در Go
- `plan.workflow.input`: payload قابل استفاده برای mapping به workflow
- `plan.artifacts`: داده‌های کمکی برای persistence/read model
- `plan.warnings`: آرایه warningهای non-fatal

## canonical response rules

- response موفق همیشه باید `plan` داشته باشد
- response خطا هرگز نباید `plan` داشته باشد
- `plan.workflow.input` باید self-contained باشد و به field hidden متکی نباشد
- `warnings` اگر وجود نداشت باید به صورت آرایه خالی برگردد
- client نباید روی ترتیب fieldها تکیه کند

## error envelope

همه خطاهای canonical باید این ساختار را داشته باشند:

```json
{
  "status": "error",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "error": {
    "code": "VOLTAGENT_INVALID_REQUEST",
    "message": "project_name is required",
    "retryable": false,
    "details": {
      "field": "input.project_name"
    }
  }
}
```

### error schema

- `status`: مقدار ثابت `error`
- `contract_version`: نسخه contract
- `request_id`: شناسه tracing
- `error.code`: کد canonical
- `error.message`: پیام human-readable
- `error.retryable`: قابل retry بودن error
- `error.details`: object اختیاری برای metadata

## taxonomy خطاها

### خطاهای validation و contract

- `VOLTAGENT_INVALID_REQUEST`
  - HTTP status: `400`
  - وقتی body معتبر نیست یا field required وجود ندارد
  - `retryable = false`

- `VOLTAGENT_UNSUPPORTED_CONTRACT_VERSION`
  - HTTP status: `400`
  - وقتی `contract_version` ناشناخته است
  - `retryable = false`

- `VOLTAGENT_UNSUPPORTED_INTENT`
  - HTTP status: `400`
  - وقتی intent خارج از scope فاز اول است
  - `retryable = false`

- `VOLTAGENT_PLAN_INVALID`
  - HTTP status: `422`
  - وقتی planner پاسخ ساخت‌یافته تولید کرده اما schema plan معتبر نیست
  - `retryable = false`

### خطاهای runtime داخل service

- `VOLTAGENT_PLAN_GENERATION_FAILED`
  - HTTP status: `500`
  - وقتی planner/model call شکست می‌خورد
  - `retryable = true`

- `VOLTAGENT_INTERNAL_ERROR`
  - HTTP status: `500`
  - وقتی failure داخلی غیرمنتظره رخ می‌دهد
  - `retryable = true`

- `VOLTAGENT_UNAVAILABLE`
  - HTTP status: `503`
  - وقتی سرویس یا planner برای پاسخ‌گویی آماده نیست
  - `retryable = true`

### خطاهای integration که معمولا توسط `VoltAgentClient` در Go synthesize می‌شوند

- `VOLTAGENT_TIMEOUT`
  - وقتی response در بازه timeout دریافت نمی‌شود
  - معمولا معادل transport timeout است
  - `retryable = true`

- `VOLTAGENT_BAD_RESPONSE`
  - وقتی HTTP status یا JSON response با contract سازگار نیست
  - مثال: JSON malformed یا body بدون `plan` در response موفق
  - `retryable = false`

## HTTP status mapping

- `200`: success
- `400`: request/contract validation error
- `422`: plan invalid
- `500`: planner/internal failure
- `503`: service unavailable

نکته:

- `VOLTAGENT_TIMEOUT` لزوما response HTTP ندارد و ممکن است فقط در client Go دیده شود
- `VOLTAGENT_BAD_RESPONSE` هم ممکن است فقط در لایه client classification تولید شود

## canonical examples

### example: unsupported intent

```json
{
  "status": "error",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "error": {
    "code": "VOLTAGENT_UNSUPPORTED_INTENT",
    "message": "intent `crypto_analysis` is not supported in phase 1",
    "retryable": false,
    "details": {
      "intent": "crypto_analysis",
      "supported_intents": [
        "deploy_website"
      ]
    }
  }
}
```

### example: invalid plan emitted by service

```json
{
  "status": "error",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "error": {
    "code": "VOLTAGENT_PLAN_INVALID",
    "message": "workflow.name is missing in planner output",
    "retryable": false,
    "details": {
      "field": "plan.workflow.name"
    }
  }
}
```

### example: bad response classified by Go client

این نمونه لزوما از خود service برنمی‌گردد و ممکن است توسط `VoltAgentClient` در Go ساخته شود:

```json
{
  "status": "error",
  "contract_version": "v1alpha1",
  "request_id": "0f7c95a0-0d8c-4d0f-a8f3-9e245bfa7a61",
  "error": {
    "code": "VOLTAGENT_BAD_RESPONSE",
    "message": "VoltAgent response body does not match the expected contract",
    "retryable": false,
    "details": {
      "reason": "missing plan field"
    }
  }
}
```

## contract of record for phase 1

برای فاز اول، طرفین روی این موارد توافق می‌کنند:

- intent پشتیبانی‌شده فقط `deploy_website` است
- `voltagent-service` فقط plan برمی‌گرداند و execution durable انجام نمی‌دهد
- Go باید response را validate کند و خطاهای transport/integration را classify کند
- هر response باید `contract_version` و `request_id` داشته باشد
- source of truth این contract همین سند است
