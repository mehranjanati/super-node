# MVP Daily Task 2 - Day 2: Remove Mocks And Failure Surface

هدف روز:

- حذف fallbackهای خاموش باقی‌مانده
- واقعی‌کردن تجربه خطا در `portal1`
- جلوگیری از این‌که failureهای backend شبیه success یا demo دیده شوند

## وضعیت شروع روز

- Day 1 contract چت و محدودسازی demo mode شروع شده است
- هنوز بعضی surfaceها ممکن است داده empty یا demo-like نشان دهند

## کارهای امروز

1. audit کردن تمام fallbackهای باقی‌مانده در فرانت
2. تبدیل هر mock خاموش به یکی از این دو حالت:
   - `demo_mode`
   - `error surface`
3. یکسان‌سازی پیام خطا در:
   - `Builder`
   - `Workflows`
   - `Logs`
   - `GlobalChat`
4. مشخص‌کردن mode dev/demo در env و docs

## وضعیت فعلی

- `completed`

کارهای انجام‌شده:
- مابقی fallbackهای خاموش در فایل `supernode.ts` (مانند ایجاد `mock-wf-id` در زمان نبود `workflowId`) پیدا و با `DEMO_MODE_ENABLED` محدود شدند.
- نمایش خطا در کامپوننت‌های `Builder`, `Workflows`, `Logs`, و `GlobalChat` یکسان‌سازی شد و همگی از استایل‌های `danger/red` استفاده می‌کنند.
- مستندات و `.env` برای توضیح `VITE_ENABLE_DEMO_MODE` به‌روزرسانی شدند.
