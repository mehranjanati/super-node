# VoltAgent Service Day 7: Legacy Cleanup And Release

هدف روز:

- بستن ambiguityهای باقی‌مانده بین مسیر جدید و routeهای legacy
- آماده‌کردن فاز 1 برای release readiness و handoff روشن
- sync نهایی docs با وضعیت واقعی repo و runtime

## وضعیت شروع روز

موارد done قبل از شروع:

- مسیر target فاز 1 روشن است: `BFF -> Go -> voltagent-service -> Go -> Temporal`
- contract داخلی و مرز مسئولیت لایه‌ها ثبت شده‌اند
- wiring backend، alignment مصرف‌کننده‌های edge و baseline health/validation تعریف شده‌اند
- `deploy_website` به‌عنوان use case رسمی phase 1 تثبیت شده است

gapهای باز:

- routeهای legacy و وابستگی‌های پنهان هنوز باید explicit audit شوند
- facade embedded قدیمی هنوز در بخشی از backend حضور دارد
- وضعیت نهایی `BFF` و `portal1` نسبت به routeهای قدیمی باید بدون ابهام ثبت شود
- checklist نهایی release و rollback هنوز باید مشخص شود

## وضعیت تاییدشده در repo فعلی

- [x] وابستگی‌های legacy اصلی مشخص‌اند: `/voltagent/manifest` و `/voltagent/execute`
- [x] `BFF` برای deploy flow از route جدید Go استفاده می‌کند
- [x] routeهای legacy هنوز در `internal/adapters/gateway/echo.go` وجود دارند (حذف شدند)
- [x] `portal1` هنوز به routeهای legacy متکی است (اصلاح شد و به routeهای `internal/tools` تغییر یافت)
- [x] facade embedded هنوز از backend حذف یا deprecate نهایی نشده است (Deprecation notice در کد اضافه شد)
- [x] release checklist و rollback plan هنوز ثبت نشده‌اند

## کارهای امروز

1. audit کردن وابستگی‌های باقی‌مانده به routeهای legacy:
   - `/voltagent/manifest`
   - `/voltagent/execute`
2. مشخص کردن تصمیم نهایی برای callerهای edge:
   - `BFF`
   - `portal1`
3. تعیین تکلیف facadeهای embedded:
   - چه چیزی در phase 1 باقی می‌ماند
   - چه چیزی deprecate می‌شود
   - چه چیزی به phase بعد carry می‌شود
4. sync نهایی اسناد:
   - daily plan
   - architecture decision
   - contract
   - day files
5. ثبت release checklist کوتاه:
   - go/no-go
   - rollback expectation
   - owner و handoff

## فایل‌های اصلی

- `internal/adapters/gateway/echo.go`
- `internal/core/services/voltagent/voltagent.go`
- `internal/core/services/voltagent/types.go`
- `internal/core/services/voltagent/fx.go`
- `BFF/index.ts`
- `portal1/src/lib/services/supernode.ts`
- `docs/Daily task/VOLTAGENT_SERVICE_DAILY_PLAN.md`

## ریسک‌های روز

- شکستن backward compatibility در callerهایی که هنوز پنهانی به routeهای legacy متکی هستند
- حذف یا محدود کردن implementation embedded فراتر از scope فاز 1
- ناهماهنگی بین docs نهایی و رفتار واقعی runtime
- release بدون تعریف روشن rollback یا owner پاسخ‌گو

## معیار done

- dependencyهای باقی‌مانده به legacy routeها explicit و قابل‌تصمیم باشند
- تصمیم `BFF` و `portal1` نسبت به مسیر نهایی ثبت شده باشد
- docs اصلی migration با repo sync و بدون ambiguity باشند
- release readiness برای phase 1 با checklist و ownerهای روشن قابل ارزیابی باشد

## Release Checklist & Rollback Plan (Phase 1)

**Go/No-Go Criteria:**
- [x] مسیر `deploy_website` از طریق `BFF` به درستی `remote_voltagent` را فراخوانی می‌کند و workflow را در Temporal استارت می‌زند.
- [x] تست‌های `HealthAndFallbackValidation` صد درصد پاس می‌شوند و fallback به درستی کار می‌کند.
- [x] مسیرهای legacy (`/voltagent/execute` و `/voltagent/manifest`) حذف شده‌اند.
- [x] رابط کاربری و `portal1` با endpointهای جدید (`/internal/tools/execute` و `/internal/tools/manifest`) هماهنگ شده‌اند.

**Rollback Plan:**
در صورت بروز بحران در Production:
1. سوییچ `VOLTAGENT_ENABLED=false` و `VOLTAGENT_USE_EMBEDDED_FALLBACK=true` در کانفیگ `Go` تا مسیر ترافیک به سمت logic داخلی قبلی تغییر کند.
2. برگرداندن routeهای legacy در `echo.go` در صورت بروز قطعی در نسخه‌های قدیمی‌تر UI.

**Owner & Handoff:**
- مسئول نگهداری سرویس: تیم Core Backend / AI
- پایش Health: از طریق مانیتورینگ مسیر `/health` در Voltagent Service و `/internal/health` در Go Gateway.

## carry-over به فاز بعد

- migration کامل chat path
- migration کامل manifest/tool discovery
- حذف تدریجی fallback embedded بعد از پایداری production-like
- observability و hardening عمیق‌تر برای rolloutهای بعدی
