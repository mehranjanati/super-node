# ادغام CMS داخلی برای اپ‌های کاربران در Nexus

## مقدمه
این سند معماری پیشنهادی برای ادغام یک CMS سبک و داخلی در پلتفرم Nexus را توصیف می‌کند. هدف: اجازه دادن به کاربران برای مدیریت محتوای اپ‌های خود (مانند صفحات لندینگ، توضیحات محصول، بلاگ) بدون نیاز به regenerate کردن UI توسط AI، با تمرکز روی قابلیت آپدیت SEO توسط ایجنت‌های Wasm.

## چرا CMS داخلی؟
- **داخلی بودن**: باید بخشی از استک Nexus باشد تا ایجنت‌ها بتوانند مستقیماً محتوا را بخوانند/بنویسند (مثلاً برای SEO dynamic).
- **ادغام با ایجنت‌ها**: ایجنت‌های Wasm می‌توانند محتوا را generate/optimize کنند (مثل meta tags، keywords).
- **سبک و Headless**: فقط API بدهد، UI در پورتال Nexus ادغام شود.
- **پشتیبانی SEO**: فیلدهایی مثل meta description، slugs، sitemaps.

## گزینه‌های بررسی‌شده
- **TinyCMS**: سبک و Git-based، اما ممکن است برای ادغام با Hasura/TiDB نیاز به customize داشته باشد.
- **Sveltia CMS**: جایگزین عالی برای Netlify CMS، Git-based، framework-agnostic، با i18n و mobile support. بهترین گزینه برای SvelteKit.
- **SveltyCMS**: مخصوص SvelteKit، با Vite/Tailwind، سریع و developer-friendly.
- **DatoCMS**: خوب برای image optimization و SEO، اما خارجی (نه داخلی).
- **Hasura-based Custom**: خودمان بسازیم روی Hasura/TiDB برای ادغام کامل.

## پیشنهاد: CMS داخلی بر پایه Hasura + Sveltia CMS
- **هسته**: از Sveltia CMS استفاده کنیم (open-source، Git-based) و آن را با Hasura ادغام کنیم تا محتوا در TiDB ذخیره شود (نه فقط Git).
- **معماری**:
  - **ذخیره‌سازی**: محتوا در جدول‌های TiDB (مثل `content_pages`, `content_blocks`) با فیلدهای SEO (meta_title, meta_description, slug).
  - **API**: Hasura GraphQL برای CRUD + Subscriptions.
  - **UI Editor**: یک صفحه در پورتال SvelteKit که از Sveltia CMS الهام گرفته، برای ویرایش محتوا.
  - **ادغام ایجنت**: ایجنت‌های Wasm می‌توانند mutationهای Hasura را صدا بزنند برای آپدیت محتوا (مثل generate meta tags بر اساس محصول).

## مزایا
- **داخلی و امن**: همه چیز داخل Super Node.
- **SEO-Friendly**: ایجنت‌ها می‌توانند محتوا را dynamic optimize کنند.
- **سبک**: بدون overhead خارجی.

## قدم‌های بعدی
1. نصب Sveltia CMS در پروژه SvelteKit.
2. ادغام با Hasura برای ذخیره‌سازی.
3. تست آپدیت SEO توسط یک ایجنت نمونه.