# VoltAgent Service Day 4: Health And Observability

هدف روز:

- قابل مشاهده کردن سلامت و failure modeهای بین Go و `voltagent-service`

## وضعیت تاییدشده در repo فعلی

- [x] `voltagent-service` در `GET /health` پاسخ contract-aware برمی‌گرداند
- [x] `voltagent-service` در `POST /plan` لاگ‌های request-level با `request_id` و `correlation_id` دارد
- [x] route `GET /internal/health` در Go وجود دارد
- [ ] `GET /internal/health` هنوز health تجمیعی `voltagent-service` را منعکس نمی‌کند
- [ ] latency و error classification برای تماس داخلی در Go هنوز ثبت نشده‌اند
- [ ] logهای صریح برای `remote_voltagent` در برابر `embedded_fallback` هنوز در backend کامل نشده‌اند
- [ ] metric baseline برای request count/error count/timeout/fallback rate هنوز پیاده نشده است

## کارهای امروز

1. افزودن health تجمیعی در Go
2. مشخص کردن وضعیت `voltagent-service` در `GET /internal/health`
3. ثبت latency و error classification برای تماس داخلی
4. اضافه کردن logهای واضح برای:
   - `remote_voltagent`
   - `embedded_fallback`
5. تعریف metricهای پایه:
   - request count
   - error count
   - timeout count
   - fallback rate

## فایل‌های اصلی

- `internal/adapters/gateway/echo.go`
- `internal/adapters/voltagentclient/client.go`
- `internal/core/services/voltagent/voltagent.go`

## معیار done

- health backend وضعیت `voltagent-service` را منعکس کند
- failureهای remote و fallback از روی log قابل تشخیص باشند
- metric baseline یا حداقل instrumentation plan ثبت شده باشد

## carry-over به روز بعد

- hardening deployment
- runbook
