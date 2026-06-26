# MVP Daily Task 2 - Day 5: Second Tool Or Targeted Polish

هدف روز:

- تثبیت MVP بعد از سبز شدن مسیر اصلی `deploy_website`
- تصمیم‌گیری کنترل‌شده بین:
  - `targeted polish` برای surfaceهای فعلی
  - یا اضافه‌کردن یک tool دوم کم‌ریسک و read-only

## اصل تصمیم روز

روز 5 فقط وقتی به سمت tool دوم می‌رود که این موارد از روزهای 1 تا 4 واقعاً سبز مانده باشند:

- chat path در runtime واقعی deploy را trigger کند
- workflow و log و artifact در SPA قابل مشاهده باشند
- failure surface در فرانت واضح و بدون fallback خاموش باشد
- smoke validation بدون ambiguity قابل تکرار باشد

اگر هر کدام از این موارد unstable بود، روز 5 باید صرف `targeted polish` شود و نه گسترش scope.

## وضعیت شروع روز

- use case اصلی `deploy_website` بسته شده اما هنوز ممکن است بعضی surfaceها product-grade کامل نباشند
- فرانت SPA هنوز ممکن است خطاهای workspace یا routing غیر MVP داشته باشد
- اضافه‌کردن tool دوم فقط در صورت حفظ اعتماد به MVP اصلی مجاز است

## مسیر مجاز A: Targeted Polish

این مسیر پیش‌فرض است مگر این‌که runtime کاملاً پایدار باشد.

کارهای امروز:

1. حذف frictionهای مهم در `GlobalChat`, `Builder`, `Workflows`, `Logs`
2. کاهش noise در UI:
   - متن‌های نامشخص
   - statusهای مبهم
   - CTAهای غیرضروری
3. بستن خطاهای فرانت که روی اعتماد MVP اثر می‌گذارند:
   - importهای شکسته
   - routingهای ناقص
   - خطاهای compile/typecheck در مسیرهای قابل مشاهده
4. ثبت یک لیست کوتاه از "known limitations" برای MVP

## مسیر مجاز B: Second Tool

این مسیر فقط وقتی مجاز است که مسیر A blocker نداشته باشد.

ویژگی‌های tool دوم:

- read-only یا کم‌ریسک باشد
- نیاز به orchestration پیچیده نداشته باشد
- از همان surfaceهای موجود قابل نمایش باشد
- failure آن به MVP اصلی ضربه نزند

گزینه‌های مناسب:

- یک tool ساده از manifest
- یک query/read-only tool برای مشاهده اطلاعات
- نسخه visibility-only از `crypto_analysis` بدون اقدام حساس

کارهای امروز در این مسیر:

1. انتخاب یک tool دوم با contract ساده
2. عبور آن از `BFF` و `Go` بدون drift
3. نمایش result آن در `GlobalChat`
4. ثبت smoke validation مختصر برای آن

## وضعیت فعلی

- `completed`

مسیر انتخاب‌شده:

- `Targeted Polish`

دلیل انتخاب:

- هرچند مسیر اصلی deploy سبز است، اما اضافه‌کردن tool دوم در این مرحله هنوز ریسک ambiguity و drift در surfaceهای MVP را بالا می‌برد
- ارزش بیشتر برای MVP فعلی از کاهش noise و حذف behaviorهای misleading در فرانت به دست می‌آید

کارهای انجام‌شده:

- fallback نمایشی manifest در `Builder` حذف شد تا UI در صورت failure به mock برنگردد
- متن‌ها و CTAهای گمراه‌کننده در `GlobalChat`, `Builder`, `Workflows`, `Logs` ساده‌تر و MVP-oriented شدند
- بخش raw payload در `GlobalChat` به حالت collapsible منتقل شد تا surface اصلی product-gradeتر بماند
- buttonها و labelهای غیرواقعی مثل `Export CSV`, `Live Tail`, `Run New`, `Reset View` حذف یا به رفتار واقعی‌تر تغییر داده شدند
- check نهایی فرانت با `npm run check` سبز شد

## Known Limitations

- فقط مسیر `deploy_website` در حال حاضر product-grade محسوب می‌شود
- tool دوم عمداً defer شد تا به use case اصلی regression وارد نشود
- بعضی routeها و surfaceهای غیر MVP هنوز نیاز به polish یا productization بیشتر دارند، اما blocker مسیر اصلی نیستند

## معیار done

- polishهای ضروری MVP بسته شده باشند
- use case اصلی `deploy_website` regress نکرده باشد
- UI ambiguity کمتری نسبت به قبل داشته باشد

## خروجی نهایی

- MVP اصلی الآن پایدارتر، تمیزتر و قابل‌اعتمادتر از قبل است
