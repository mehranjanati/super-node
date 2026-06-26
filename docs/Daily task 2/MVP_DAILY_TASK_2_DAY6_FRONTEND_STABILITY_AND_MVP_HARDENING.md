# MVP Daily Task 2 - Day 6: Frontend Stability And MVP Hardening

هدف روز:

- کاهش خطاهای فرانت که مستقیماً به اعتماد MVP آسیب می‌زنند
- تثبیت SPA برای use caseهای اصلی و جلوگیری از شکست‌های قابل مشاهده

## چرا این روز لازم است

حتی اگر مسیر اصلی deploy کار کند، MVP تا وقتی در فرانت:

- import شکسته داشته باشد
- route ناقص داشته باشد
- type/compile errorهای تاثیرگذار داشته باشد
- یا stateهای صفحه بین routeها inconsistent باشند

هنوز product-grade محسوب نمی‌شود.

## وضعیت شروع روز

- use case اصلی از نظر contract و smoke path بسته شده است
- بعضی componentها ممکن است هنوز legacy یا non-MVP باشند
- باید بین "workspace cleanliness" و "MVP risk" تمایز گذاشته شود

## اصل تصمیم

در این روز فقط خطاهایی رفع می‌شوند که یکی از این اثرها را دارند:

- مانع build/dev server شوند
- routeهای اصلی SPA را بشکنند
- باعث خطای visible در preview/browser شوند
- توسعه use case اصلی را پرهزینه‌تر و پرریسک‌تر کنند

خطاهای صرفاً تزئینی یا unrelated فقط در صورت کم‌هزینه بودن انجام می‌شوند.

## کارهای امروز

1. دسته‌بندی خطاهای فرانت به سه گروه:
   - `blocking`
   - `mvp-relevant`
   - `defer`
2. رفع خطاهای `blocking` و `mvp-relevant`
3. یکسان‌سازی routeهای اصلی:
   - `#/`
   - `#/builder`
   - `#/workflows`
   - `#/logs`
   - `GlobalChat`
4. ثبت لیست خطاهای defer شده با دلیل تعویق
5. اجرای check نهایی برای اطمینان از پایداری مسیرهای اصلی

## وضعیت فعلی

- `completed`

کارهای انجام‌شده:

- hash routeهای اصلی در `portal1` یکسان‌سازی شدند تا `#/`، `#/dashboard` و routeهای MVP مثل `#/builder`, `#/workflows`, `#/logs` رفتار یکدست داشته باشند
- یک helper مشترک در `portal1/src/lib/utils.ts` اضافه شد تا نرمال‌سازی route، page title و active state از یک منبع واحد استفاده کنند
- `GlobalChat` از state واقعی hash route استفاده می‌کند و حالا context صفحه را به‌صورت `currentPath` و `currentRoute` به `BFF` می‌فرستد؛ این باگ باعث می‌شد BFF عملاً همیشه `"/"` ببیند
- header چت و `BFF/index.ts` با contract جدید route-context هم‌راستا شدند تا system prompt واقعاً از صفحه فعال کاربر آگاه باشد
- checkهای فرانت برای مسیرهای اصلی دوباره اجرا شدند و `npm run check`, `npm run build` و diagnostics فایل‌های edit شده سبز ماندند

## دسته‌بندی خطاها

### blocking

- مورد blocking در `portal1` برای build/check پیدا نشد

### mvp-relevant

- drift بین hash routing و context ارسالی از `GlobalChat` به `BFF`
- ناهماهنگی در نرمال‌سازی `#/` و routeهای اصلی بین shell، topbar و sidebar

### defer

- فایل legacy `portal1/src/lib/components/global/main-router.svelte` هنوز در repo وجود دارد اما در مسیر MVP فعلی استفاده نمی‌شود؛ حذف آن به cleanup جداگانه موکول شد
- routeها و surfaceهای non-MVP مثل marketplace/foundry همچنان به polish بیشتر نیاز دارند، اما blocker use case اصلی نیستند

## validation

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check
npm run build
```

خروجی مورد انتظار:

- `svelte-check found 0 errors and 0 warnings`
- build production بدون خطای route/type انجام شود
- مسیرهای `#/`, `#/builder`, `#/workflows`, `#/logs` در SPA از یک نرمال‌سازی مشترک استفاده کنند

## معیار done

- routeهای اصلی SPA بدون خطای visible بالا بیایند
- typecheck/build برای بخش‌های MVP در وضعیت پایدار قرار بگیرد
- تیم بداند کدام errorها عمداً defer شده‌اند و چرا

## خروجی مورد انتظار

در پایان روز 6، فرانت دیگر برای use case اصلی یک منبع بی‌اعتمادی نباشد و اجرای MVP در preview/dev environment قابل اتکا بماند.
