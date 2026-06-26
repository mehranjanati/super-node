# MVP Daily Task 3 - Day 1: Agent Draft Foundation

هدف روز:

- خروج `Foundry` و `Projects` از حالت نمایشی
- ساختن اولین flow واقعی برای draft agent در فرانت

## وضعیت شروع روز

- `Daily task 2` نشان داد که MVP فعلی بیشتر روی مسیر deploy متمرکز است
- `Foundry` فرم و preview نمایشی داشت اما draft واقعی ذخیره نمی‌کرد
- `Projects` کارت ثابت و placeholder داشت و به state واقعی agent وصل نبود
- store مربوط به agentها هم هنوز با mock و delay مصنوعی کار می‌کرد

## کارهای امروز

1. حذف رفتار mock از draft flow فرانت:
   - `portal1/src/lib/stores/agents.ts`
2. تبدیل persistence از mock به local-first draft storage:
   - ذخیره در `localStorage`
   - امکان load مجدد draftها
3. بازنویسی `Foundry` به‌عنوان draft studio واقعی:
   - تعریف فیلدهای agent
   - ذخیره draft
   - بازکردن draft قبلی
   - حذف draft
   - preview واقعی JSON
4. بازنویسی `Projects` برای نمایش draftهای واقعی
5. اجرای check نهایی

## وضعیت فعلی

- `completed`

## کارهای انجام‌شده

- `agentsStore` از mock data و timeout مصنوعی خارج شد و حالا draftهای agent را در `localStorage` می‌خواند و می‌نویسد
- `Foundry` بازنویسی شد تا کاربر بتواند:
  - نام agent وارد کند
  - prompt/behavior تعریف کند
  - type, language, framework و runtime انتخاب کند
  - draft را ذخیره و دوباره باز کند
  - draft را حذف کند
  - config JSON همان draft را ببیند
- `Projects` از کارت ثابت خارج شد و حالا draftهای واقعی agent را list می‌کند
- مسیر `Projects -> Foundry` برای بازکردن draft انتخاب‌شده برقرار شد
- `portal1` بعد از این تغییرات با `npm run check` سبز ماند

## فایل‌های اصلی

- `portal1/src/lib/stores/agents.ts`
- `portal1/src/lib/components/foundry/Foundry.svelte`
- `portal1/src/lib/components/projects/Projects.svelte`

## validation

```bash
cd /Users/elbaan/Documents/super\ node\ 1/portal1
npm run check
```

خروجی مورد انتظار:

- `svelte-check found 0 errors and 0 warnings`
- draft ذخیره‌شده بعد از refresh از بین نرود
- `Projects` draftهای ساخته‌شده در `Foundry` را نشان بدهد

## limitationهای فعلی

- draft agent هنوز به execution contract واقعی backend وصل نشده است
- نتیجه این روز "design is real" را می‌بندد، نه هنوز "agent execution is real"
- persistence فعلاً local-first است و multi-user/backend synced نیست

## carry-over به روز بعد

- انتخاب use case واقعی برای agent اول
- تعریف contract بین draft agent و runtime path
- تصمیم روی این‌که agent اول از چه surface و چه toolهای read-only استفاده کند
