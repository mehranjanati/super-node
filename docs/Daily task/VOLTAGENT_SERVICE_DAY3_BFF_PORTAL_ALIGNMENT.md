# VoltAgent Service Day 3: BFF And Portal Alignment

هدف روز:

- align کردن مصرف‌کننده‌های edge با مسیر جدید backend
- کم کردن اتکای frontend به routeهای legacy

## وضعیت تاییدشده در repo فعلی

- [x] `BFF/index.ts` برای `deploy_website` به route جدید Go یعنی `/internal/tools/deploy` وصل است
- [x] مسیر target فاز 1 در docs روشن است: `BFF -> Go -> voltagent-service -> Go -> Temporal`
- [x] `portal1/src/lib/services/supernode.ts` برای `system__deploy_website` به route جدید `/internal/tools/deploy` migrate شده است
- [x] تست unit برای migration سرویس `portal1` اضافه شده و سبز است
- [x] `portal1/src/lib/services/supernode.ts` برای ابزارهای non-deploy به `/internal/tools/execute` منتقل شد
- [x] `portal1/src/lib/services/supernode.ts` manifest را از `/internal/tools/manifest` می‌گیرد
- [x] dependency پنهان frontend به routeهای legacy در لایه service حذف شد

## کارهای امروز

1. بازبینی `BFF/index.ts` و تصمیم نهایی برای contract بین BFF و Go
2. هماهنگ‌کردن payload `deploy_website` با backend جدید
3. تعیین تکلیف `portal1`:
   - [x] به مسیر جدید migrate شود
4. مستندسازی واضح مسیر نهایی:
   - `BFF -> Go -> voltagent-service -> Go -> Temporal`

## تصمیم نهایی

- contract داخلی فاز 1 برای callerهای edge این است:
  - `GET /internal/tools/manifest`
  - `POST /internal/tools/execute`
  - `POST /internal/tools/deploy`
- routeهای `/voltagent/manifest` و `/voltagent/execute` فعلاً برای compatibility داخل gateway باقی می‌مانند
- `portal1` دیگر نباید مستقیماً به routeهای legacy وصل باشد
- `BFF` برای `deploy_website` روی `/internal/tools/deploy` باقی می‌ماند و نیاز به fallback جداگانه در frontend ندارد

## فایل‌های اصلی

- `BFF/index.ts`
- `portal1/src/lib/services/supernode.ts`
- `docs/Daily task/VOLTAGENT_SERVICE_DAILY_PLAN.md`

## ریسک‌های روز

- شکستن backward compatibility در frontend
- drift بین route جدید و callerهای فعلی
- باقی ماندن routeهای legacy در backend و امکان بازگشت ناخواسته callerها به آن‌ها

## معیار done

- BFF با payload سازگار با backend جدید کار کند
- تصمیم portal explicit ثبت شود
- هیچ ambiguity درباره route اصلی phase 1 باقی نماند
- frontend serviceها فقط از contract داخلی `/internal/tools/*` استفاده کنند

## carry-over به روز بعد

- health aggregation
- logging و metrics
