# VoltAgent Service Day 2: Backend Runtime Wiring

هدف روز:

- تکمیل wiring backend برای `Go -> voltagent-service`
- پایدار کردن fallback برای `deploy_website`
- ثبت دقیق منبع planning در داده‌های اجرایی

## وضعیت شروع روز

موارد done قبل از شروع:

- `VoltAgentClient` وجود دارد
- `POST /plan` و `GET /health` در `voltagent-service` فعال هستند
- wiring اولیه `deploy_website` با fallback در backend اضافه شده است

gapهای باز:

- facade runtime هنوز به‌صورت کامل از bridge implementation جدا نشده است
- همه call siteهای runtime هنوز migrated نیستند
- coverage اختصاصی برای `voltagentclient` هنوز اضافه نشده است

## وضعیت تاییدشده در repo فعلی

- [x] `VoltAgentClient` در `internal/adapters/voltagentclient/client.go` وجود دارد
- [x] `POST /plan` و `GET /health` در `voltagent-service/src/index.ts` فعال هستند
- [x] wiring مسیر `deploy_website` از Go به `voltagent-service` با fallback embedded در `internal/core/services/voltagent/voltagent.go` وجود دارد
- [x] metadataهای `request_id`, `correlation_id`, `source` از Go به remote plan request عبور داده می‌شوند
- [x] policyهای runtime یعنی `Enabled` و `UseEmbeddedFallback` در config و service wiring فعال‌اند
- [x] `planning_source` در response وجود دارد و در read model/persistence هم ذخیره می‌شود
- [x] تست هدفمند backend برای gateway و مسیرهای مرتبط سبز است
- [ ] facade runtime هنوز به‌صورت کامل از bridge implementation جدا نشده است
- [ ] همه call siteهای runtime هنوز migrated نیستند

## کارهای امروز

1. ثبت `planning_source` در response, persistence و read model
2. عبور دادن metadata درخواست:
   - `request_id`
   - `correlation_id`
   - `source`
3. تمیز کردن facade داخلی برای remote planning در برابر embedded fallback
4. روشن کردن policyهای runtime:
   - `Enabled`
   - `UseEmbeddedFallback`
5. بازبینی خروجی routeهای درگیر برای سازگاری با callerها

## فایل‌های اصلی

- `internal/core/services/voltagent/types.go`
- `internal/core/services/voltagent/voltagent.go`
- `internal/adapters/gateway/echo.go`
- `internal/adapters/gateway/workflow_read_model.go`
- `internal/core/domain/workflow_execution.go`

## معیار done

- `deploy_website` از remote plan به execution واقعی برسد
- در fallback هم execution از کار نیفتد
- `planning_source` در response یا داده اجرایی قابل مشاهده باشد
- build/test هدفمند packageهای backend سبز باشند

## carry-over به روز بعد

- migrate کردن BFF و portal
- health aggregation
- observability کامل
