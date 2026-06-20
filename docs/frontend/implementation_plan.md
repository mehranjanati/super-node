> **⚠️ توجه (آپدیت استراتژی):** این سند ممکن است شامل جزئیات قدیمی مربوط به نسخه‌های اولیه فرانت‌اند باشد. سیستم اکنون از رویکرد **Template-Driven** و **GitOps کاربر‌محور** استفاده می‌کند. برای مشاهده استراتژی نهایی حتماً به [MASTER_PLAN.md](../../MASTER_PLAN.md) و [MVP_DEVELOPMENT_PLAN.md](../MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# Implementation Plan - Sveltia CMS Integration

This plan outlines the steps to integrate Sveltia CMS into the Nexus Portal, transitioning from mocked data to a dynamic, CMS-driven architecture.

## User Objective
Integrate Sveltia CMS to manage content (Service Status, Announcements) dynamically within the portal, replacing hardcoded values.

## Proposed Plan

### Step 1: CMS Initialization & Configuration
**Goal**: Enable the CMS admin interface.
- Create `static/admin/index.html` loading the Sveltia CMS script.
- Create `static/admin/config.yml` with basic configuration (GitHub backend or local git proxy for dev).
- **Verification**: User can access `/admin` and see the Sveltia login/interface.

### Step 2: Define Content Collections (Schema)
**Goal**: Model the data structures needed for the Dashboard.
- Update `config.yml` to include collections for:
    - **System Alerts** (Announcements/Recent Activity)
    - **Service Status** (Overrides for Matrix/LiveKit status messages)
- **Verification**: User can create a new "Alert" or "Service Status" entry in the CMS.

### Step 3: Frontend Data Integration
**Goal**: Display CMS-managed content in the Svelte application.
- Create a utility (e.g., `src/lib/utils/content.ts`) to import generic markdown/JSON content.
- Update `Dashboard.svelte` to fetch "Recent Activity" from the CMS content instead of using hardcoded rows.
- **Verification**: Dashboard reflects changes made in the CMS admin panel.
