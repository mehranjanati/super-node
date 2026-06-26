# VoltAgent Service Day 6: Tests And Validation

هدف روز:

- بالا بردن confidence برای مسیر اصلی `BFF -> Go -> voltagent-service -> Go -> Temporal`
- تبدیل wiring و health baseline به validation قابل تکرار
- اطمینان از اینکه fallback و remote path هر دو قابل اتکا هستند

## وضعیت شروع روز

موارد done قبل از شروع:

- baseline معماری، contract و scope phase 1 روشن شده‌اند
- wiring backend برای `deploy_website` با remote planning و fallback embedded وجود دارد
- alignment لایه‌های edge در روز 3 و hardening پایه infra در روز 5 مشخص شده‌اند
- health و observability baseline برای تماس داخلی تعریف شده است

gapهای باز:

- coverage هدفمند برای `voltagentclient` و mapping runtime هنوز کافی نیست
- smoke test قابل تکرار برای compose stack هنوز formal نشده است
- validation دقیق remote success path در برابر embedded fallback هنوز کامل نیست
- تست end-to-end برای use case رسمی phase 1 هنوز باید explicit ثبت و اجرا شود

## وضعیت تاییدشده در repo فعلی

- [x] تست هدفمند Go برای `internal/adapters/gateway` سبز است
- [x] read model حالا با تست، نگه‌داشتن `planning_source` را بعد از persist/update تایید می‌کند
- [x] packageهای `internal/core/services/voltagent` و `internal/adapters/voltagentclient` حداقل build/test check را پاس می‌کنند
- [x] هنوز برای `voltagentclient` تست اختصاصی unit اضافه نشده است (انجام شد)
- [x] smoke test رسمی برای compose stack ثبت نشده است
- [x] validation رسمی end-to-end برای `deploy_website` در remote path و fallback path کامل نشده است (تست در `voltagent_test.go` اضافه شد)

## کارهای امروز

1. افزودن testهای هدفمند backend:
   - unit test برای `voltagentclient`
   - test برای mapping response `POST /plan`
   - test برای fallback policy
2. بازبینی validation مسیر service:
   - success response contract
   - error envelope
   - request metadata propagation
3. تعریف smoke test برای stack:
   - بالا آمدن compose
   - `GET /health`
   - `POST /plan`
4. اجرای validation use case `deploy_website` در دو حالت:
   - remote planning success
   - embedded fallback
5. ثبت outcome testها و failure modeهای شناخته‌شده برای release readiness

## فایل‌های اصلی

- `internal/adapters/voltagentclient/client.go`
- `internal/core/services/voltagent/voltagent.go`
- `internal/adapters/gateway/echo.go`
- `voltagent-service/src/index.ts`
- `docs/VOLTAGENT_SERVICE_CONTRACT.md`

## ریسک‌های روز

- flaky شدن تست‌ها به‌خاطر وابستگی به service readiness یا network timing
- قاطی شدن تست‌های contract با behaviorهای خارج از scope فاز اول
- پوشش‌ندادن fallback path با وجود سبز بودن remote path
- false confidence به‌خاطر تست‌های خیلی سطحی یا وابسته به mockهای غیرواقعی

## معیار done

- packageهای هدفمند backend test معنادار و سبز داشته باشند
- smoke test stack مسیر `health` و `plan` را پوشش دهد
- `deploy_website` هم در remote path و هم در fallback path validate شده باشد
- failure modeهای اصلی و نتیجه validation به‌صورت قابل‌ارجاع ثبت شده باشند

## carry-over به روز بعد

- cleanup مسیرهای legacy
- sync نهایی docs
- release checklist و go/no-go review
