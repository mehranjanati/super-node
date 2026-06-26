# MVP Daily Task 3 - Day 5: Result Surface And UX Hardening

هدف روز:

- نمایش output معنادار agent به‌صورت product-usable
- بستن failure surfaceها و ابهام‌های UX مرتبط با agent MVP

## چرا این روز حیاتی است

تا پایان Day 4، agent اول واقعاً کاری انجام می‌دهد.

اما اگر result surface آن هنوز خام، ناپایدار یا مبهم باشد:

- کاربر تفاوت بین "tool payload" و "agent result" را حس نمی‌کند
- تجربه دوباره شبیه debug panel می‌شود نه product behavior
- MVP کاربردی از نظر UX بسته نمی‌شود، حتی اگر backend technically کار کند

Day 5 باید اطمینان بدهد که خروجی agent:

- خوانا است
- failure آن واضح است
- در context draft انتخاب‌شده معنی دارد
- به جای developer-only view، به user-facing result نزدیک است

## وضعیت شروع روز

- `workflow_insight` در `BFF` و `GlobalChat` اجرا می‌شود
- output ساخت‌یافته اولیه وجود دارد
- اما UX هنوز کامل نیست:
  - empty stateها باید واضح‌تر شوند
  - failure stateها باید کمتر raw باشند
  - `Projects` و `Foundry` هنوز contract/result relationship را کامل نشان نمی‌دهند

## اصل تصمیم

Day 5 نباید use case جدید اضافه کند.

این روز فقط مجاز است:

- surfaceهای موجود را readableتر کند
- ambiguity را کم کند
- failureهای قابل‌انتظار را بهتر نشان دهد
- consistency بین `Foundry`, `Projects`, `GlobalChat` را بالا ببرد

این روز مجاز نیست:

- agent دوم اضافه کند
- capability جدید تعریف کند
- backend orchestration جدید بسازد
- execution را multi-step یا stateful‌تر کند

## کارهای امروز

1. hardening `GlobalChat` result surface:
   - تفکیک روشن بین message عادی و result card
   - بهترکردن copy برای loading, empty, failure
2. بهبود context visibility:
   - نمایش روشن draft انتخاب‌شده
   - نمایش capability و result surface در جای مناسب
3. hardening `Projects`:
   - نمایش capability هر draft روی کارت
   - نمایش execution mode یا result surface به‌صورت خلاصه
4. hardening `Foundry`:
   - copy واضح‌تر درباره این‌که draft دقیقاً چه کاری انجام خواهد داد
   - کاهش ambiguity بین `deploy` و `workflow insight`
5. بستن failure surfaceهای اصلی:
   - `workflow_not_found`
   - `no_workflows_available`
   - `workflow_data_unavailable`
   - draft انتخاب‌نشده

## خروجی مورد انتظار از این روز

در پایان Day 5 باید این موارد سبز باشند:

1. user بتواند تشخیص بدهد کدام draft agent فعال است
2. user بفهمد این agent چه capabilityای دارد
3. result agent به‌صورت readable و نه فقط raw payload دیده شود
4. اگر data وجود نداشت یا workflow پیدا نشد، failure واضح و non-misleading باشد
5. `Projects` و `Foundry` به contract/runtime path فعلی اشاره روشن‌تری داشته باشند

## فایل‌های محتمل

- `portal1/src/lib/components/chat/GlobalChat.svelte`
- `portal1/src/lib/components/foundry/Foundry.svelte`
- `portal1/src/lib/components/projects/Projects.svelte`
- `portal1/src/lib/types/index.ts`

## validation روز

### validation دستی

1. یک draft با capability `workflow_insight` بساز
2. draft را ذخیره کن
3. همان draft را از `Projects` باز کن
4. در `GlobalChat` سوالی مثل این بپرس:
   - `آخرین workflowها را خلاصه کن`
   - `کدام workflow fail شده؟`
5. بررسی کن:
   - draft فعال مشخص باشد
   - result card readable باشد
   - در صورت نبود data، copy واضح برگردد

### validation فنی

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check
```

## معیار done

- result surface فعلی دیگر شبیه debug-only UI نباشد
- capability فعال agent در فرانت پنهان نماند
- failure و empty stateها واضح و کم‌ابهام باشند
- هیچ use case جدیدی خارج از scope `workflow_insight` وارد نشده باشد

## carry-over

- اگر UX هنوز مبهم بود، فقط polishهای لازم به Day 6 منتقل می‌شوند
- هر نوع گسترش scope به agent دوم یا tool جدید، failure این روز محسوب می‌شود
