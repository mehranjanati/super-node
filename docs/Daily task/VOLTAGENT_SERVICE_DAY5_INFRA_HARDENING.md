# VoltAgent Service Day 5: Infra Hardening

هدف روز:

- harden کردن deployment و config برای مسیر `Go -> voltagent-service`
- کم کردن ریسک‌های runtime بعد از روشن شدن health و observability baseline
- آماده‌کردن stack برای validation پایدار در روزهای بعد

## وضعیت شروع روز

موارد done قبل از شروع:

- `voltagent-service` در composeهای اصلی stack حضور دارد
- contract داخلی `Go -> voltagent-service` تعریف شده است
- wiring backend برای `deploy_website` با fallback embedded شروع شده است
- محورهای health و observability در روز 4 مشخص شده‌اند

gapهای باز:

- env matrix بین `Go`, compose و `voltagent-service` هنوز به‌صورت کامل harden نشده است
- internal-only exposure برای `voltagent-service` هنوز باید explicit بررسی شود
- readiness/startup ordering و failure modeهای compose هنوز نیاز به بازبینی دارند
- runbook بالا آوردن stack و عیب‌یابی اولیه هنوز کامل نیست

## وضعیت تاییدشده در repo فعلی

- [x] `docker-compose.yml` و `docker-compose.app.yml` شامل `voltagent-service` هستند
- [x] config رسمی `voltagent` در `config/config.yaml` و `internal/config/config.go` وجود دارد
- [x] defaultهای `base_url`, `timeout`, `enabled`, `use_embedded_fallback`, `contract_version` در backend تعریف شده‌اند
- [ ] internal-only exposure برای `voltagent-service` هنوز به‌صورت policy روشن مستند یا enforce نشده است
- [ ] health/readiness compose برای این migration هنوز harden نشده است
- [ ] runbook اجرایی stack هنوز ثبت نشده است

## کارهای امروز

1. بازبینی و تثبیت envهای runtime:
   - `NEXUS_VOLTAGENT_*`
   - base URL
   - timeout
   - fallback policy
2. بررسی composeها برای:
   - dependency ordering
   - health check
   - readiness behavior
   - internal-only exposure
3. اطمینان از اینکه `voltagent-service` فقط از مسیر موردنیاز backend قابل دسترسی است
4. ثبت runbook کوتاه برای:
   - بالا آوردن stack
   - تشخیص failure اولیه
   - بررسی health بین Go و `voltagent-service`
5. بازبینی defaultهای config تا fallback ناخواسته، خطاهای remote را پنهان نکند

## فایل‌های اصلی

- `docker-compose.yml`
- `docker-compose.app.yml`
- `docker-compose.infra.yml`
- `config/config.yaml`
- `internal/config/config.go`
- `internal/adapters/voltagentclient/client.go`

## ریسک‌های روز

- exposed ماندن `voltagent-service` خارج از boundary داخلی موردنظر
- drift بین envهای compose و config واقعی backend
- سبز بودن ظاهری health در حالی که planner یا dependency آماده نیست
- masking شدن failureهای remote به‌خاطر fallback بیش از حد permissive

## معیار done

- envهای اصلی `voltagent-service` و backend explicit و قابل‌ردیابی باشند
- compose behavior برای startup و health ambiguity نداشته باشد
- policy دسترسی `voltagent-service` روشن باشد و exposure ناخواسته نداشته باشیم
- runbook اولیه برای بالا آوردن stack و debug failureهای رایج ثبت شده باشد

## carry-over به روز بعد

- unit و integration test
- smoke test compose
- validation مسیر end-to-end برای `deploy_website`
