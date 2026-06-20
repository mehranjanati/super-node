> **⚠️ توجه (آپدیت معماری):** این سند ممکن است شامل جزئیات قدیمی باشد. برای مشاهده معماری نهایی سیستم (شامل اضافه شدن Rivet، Matrix، OpenClaw، Redpanda و LiveKit) حتماً به [FULL_SYSTEM_ARCHITECTURE.md](./FULL_SYSTEM_ARCHITECTURE.md) و برای برنامه اجرایی به [MVP_DEVELOPMENT_PLAN.md](./MVP_DEVELOPMENT_PLAN.md) مراجعه کنید.

# 🔗 پلن اتصال BFF به بک‌اند Go (فاز ۵)

این سند مراحل دقیق جایگزینی Mock فعلی در لایه BFF با یک ارتباط واقعی HTTP/REST به هسته پردازشی Go (Super Node) را شرح می‌دهد.

---

## 🎯 هدف
زمانی که هوش مصنوعی (در لایه BFF) تصمیم می‌گیرد ابزاری مثل `deploy_website` را اجرا کند، BFF باید یک درخواست به سرور Go بفرستد. Go عملیات سنگین (مثل ساخت کانتینر) را انجام داده و نتیجه را به BFF برمی‌گرداند تا در نهایت به کاربر نمایش داده شود.

---

## 📝 قرارداد API (API Contract)

برای اینکه Bun و Go بتوانند با هم صحبت کنند، باید روی یک ساختار داده توافق کنیم:

*   **آدرس (Endpoint):** `POST http://localhost:8080/internal/tools/deploy`
*   **هدرها:** `Content-Type: application/json`
*   **بدنه درخواست (Request Body) از Bun به Go:**
    ```json
    {
      "project_name": "my-awesome-site",
      "template": "svelte"
    }
    ```
*   **بدنه پاسخ (Response Body) از Go به Bun:**
    ```json
    {
      "status": "success",
      "url": "https://my-awesome-site.nexus.app",
      "message": "Deployment completed in 4.2s"
    }
    ```

---

## 🛠️ گام‌های اجرایی

### گام ۱: پیاده‌سازی هندلر در Go (Super Node)
در پروژه Go خود، باید یک سرور HTTP داخلی (Internal API) راه‌اندازی کنید که به درخواست‌های BFF گوش دهد.

```go
// مسیر فرضی: cmd/server/internal_api.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeployRequest struct {
	ProjectName string `json:"project_name"`
	Template    string `json:"template"`
}

type DeployResponse struct {
	Status  string `json:"status"`
	URL     string `json:"url"`
	Message string `json:"message"`
}

func DeployToolHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("[Go SuperNode] Received deploy request for: %s
", req.ProjectName)

	// 🔴 در اینجا منطق واقعی دیپلوی (Docker/Podman) اجرا می‌شود 🔴
	time.Sleep(3 * time.Second) // شبیه‌سازی زمان دیپلوی

	response := DeployResponse{
		Status:  "success",
		URL:     fmt.Sprintf("https://%s.nexus.app", req.ProjectName),
		Message: "Container deployed successfully via Go",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/internal/tools/deploy", DeployToolHandler)
	fmt.Println("Go Internal API running on :8080")
	http.ListenAndServe(":8080", nil)
}
```

### گام ۲: بروزرسانی کد BFF (حذف Mock)
در فایل `BFF/index.ts`، باید بخش `execute` ابزار `deploy_website` را تغییر دهیم تا واقعاً به Go درخواست بفرستد.

```typescript
// تغییرات در BFF/index.ts
execute: async ({ projectName, template }) => {
  console.log(`[BFF] Sending deploy request to Go for ${projectName}...`);
  
  try {
    const goResponse = await fetch('http://localhost:8080/internal/tools/deploy', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        project_name: projectName, 
        template: template || 'default' 
      })
    });

    if (!goResponse.ok) {
      throw new Error(`Go server responded with status: ${goResponse.status}`);
    }

    const data = await goResponse.json();
    console.log(`[BFF] Received response from Go:`, data);
    
    return data; // این داده مستقیماً به Vercel AI SDK و سپس به فرانت‌اند می‌رود

  } catch (error) {
    console.error("[BFF] Failed to connect to Go Super Node:", error);
    return {
      status: "error",
      message: "Failed to execute tool in Super Node. Is the Go server running?"
    };
  }
}
```

---

## 🚀 گام‌های پیشرفته (برای آینده)

اگر فرآیند دیپلوی در Go بیش از ۳۰ ثانیه طول بکشد، ممکن است درخواست HTTP تایم‌اوت (Timeout) شود. در این حالت باید معماری را به **Redis Pub/Sub** ارتقا دهیم:

1. **Bun** یک درخواست سریع به Go می‌فرستد و یک `job_id` می‌گیرد.
2. **Bun** به کانال Redis با نام `job_<job_id>` سابسکرایب می‌کند.
3. **Go** دیپلوی را در بک‌گراند انجام می‌دهد و لاگ‌ها را در Redis پابلیش (Publish) می‌کند.
4. **Bun** لاگ‌ها را از Redis می‌خواند و از طریق استریم به فرانت‌اند می‌فرستد.
