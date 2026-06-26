# MVP Daily Task 2 - Day 3: Workflow Artifacts And Result Surface

هدف روز:

- تبدیل `workflow_id` از تنها خروجی visible به یک نتیجه محصولی‌تر
- نشان‌دادن artifactهای معنی‌دار در UI

## وضعیت شروع روز

- workflow و status در UI دیده می‌شوند
- artifactهای واقعی هنوز کامل و product-grade نیستند

## کارهای امروز

1. audit کردن artifactهای موجود در backend read model
2. نرمال‌سازی fieldهای خروجی مثل:
   - `repo_url`
   - `pr_url`
   - `liveUrl`
   - `previewUrl`
3. اضافه کردن surface واضح برای artifact در:
   - `GlobalChat`
   - `Builder`
   - `Workflows`
4. ثبت این‌که کدام activityها هنوز placeholder هستند

## وضعیت فعلی

- `completed`

کارهای انجام‌شده:
- فیلدهای artifacts شامل `liveUrl`، `previewUrl` و `repoUrl` در backend و SPA (فایل‌های `supernode.ts`، `GlobalChat.svelte` و `Builder.svelte`) هماهنگ و نرمال‌سازی شدند.
- در UI برای هر کدام از artifactها، لینک‌های معتبر (Open Preview URL، View Source Code و غیره) اضافه شدند که تنها در صورت وجود هر فیلد، در UI رندر می‌شوند.
- وضعیت Activityها در فایل مستندسازی شد.
