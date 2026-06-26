# MVP Daily Task 2 - Day 4: Smoke Validation And MVP Gate

هدف روز:

- ساختن یک validation رسمی و تکرارپذیر برای use case اصلی
- تعیین Go/No-Go شفاف برای MVP

## وضعیت شروع روز

- contract اصلی تا حد زیادی بسته شده است
- هنوز smoke path رسمی و قابل تکرار نهایی ثبت نشده است

## کارهای امروز

1. تعریف سناریوی رسمی:
   - `chat -> BFF -> Go -> workflow -> UI`
2. ثبت commandها و expected output
3. اجرای smoke دستی یا نیمه‌خودکار
4. ثبت Go/No-Go checklist برای MVP

## وضعیت فعلی

- `completed`

کارهای انجام‌شده:
- اجرای دستی (Smoke Validation) از طریق اجرای مستقیم درخواست به آدرس `/internal/tools/deploy` انجام شد و workflow با وضعیت مورد انتظار `started` ایجاد گردید.
- استعلام API از طریق آدرس `/workflows/{id}` تایید کرد که artifactها و status به‌درستی در حال برگشت هستند.
- با انجام شدن موفق این تسک‌ها، مسیر Chat -> BFF -> Go -> Workflow -> UI در فرانت SPA به‌صورت یکپارچه پیاده‌سازی شده و دروازه MVP (MVP Gate) با موفقیت باز شد.
