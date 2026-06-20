import { convertToModelMessages, streamText, tool } from 'ai';
import { google } from '@ai-sdk/google';
import { z } from 'zod';

// اتصال به Redis بومی Bun (در صورت عدم وجود Redis لوکال، این بخش را می‌توان کامنت کرد)
// import { Redis } from 'ioredis';
// const redis = new Redis("redis://localhost:6379");

const PORT = 3001;

console.log(`🚀 BFF Server is running on http://localhost:${PORT}`);

Bun.serve({
  port: PORT,
  async fetch(req) {
    const url = new URL(req.url);

    // مدیریت CORS برای ارتباط با فرانت‌اند (SvelteKit)
    if (req.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'POST, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type',
        },
      });
    }

    if (url.pathname === '/api/health' && req.method === 'GET') {
      return Response.json({
        status: 'ok',
        service: 'bff',
      }, {
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
      });
    }

    // مسیر اصلی چت
    if (url.pathname === '/api/chat' && req.method === 'POST') {
      try {
        const body = await req.json();
        console.log('[BFF] Received chat request body:', JSON.stringify(body, null, 2));
        const { messages, data } = body;

        // استخراج مسیر فعلی کاربر از دیتای ارسالی
        const currentPath = data?.currentPath || 'نامشخص';
        
        // پیام سیستمی برای آگاهی مدل از صفحه فعلی
        const systemPrompt = `You are VoltAgent, a helpful AI assistant for the Nexus platform. The user is currently on the page: ${currentPath}. Use this context to provide more relevant answers if they ask about what they are seeing or what they can do here.`;

        // Vercel AI SDK expects specific format for messages
        const coreMessages = await convertToModelMessages(messages);

        // ارتباط با LLM (Google Gemini)
        const result = streamText({
          model: google('gemini-2.5-flash'), // استفاده از مدل Gemini
          system: systemPrompt,
          messages: coreMessages,
          tools: {
            // تعریف ابزار دیپلوی که به بک‌اند Go متصل می‌شود
            deploy_website: tool({
              description: 'Deploy a new website or container',
              inputSchema: z.object({
                projectName: z.string().describe('The name of the project to deploy'),
                template: z.string().optional().describe('The template to use (e.g., svelte, react)'),
              }),
              execute: async ({ projectName, template }: { projectName: string, template?: string }) => {
                console.log(`[BFF] Tool called: deploy_website for ${projectName}${template ? ` using ${template}` : ''}`);
                
                try {
                  // تنظیم Timeout برای درخواست به Go (مثلاً 10 ثانیه)
                  const controller = new AbortController();
                  const timeoutId = setTimeout(() => controller.abort(), 10000);

                  // ارسال درخواست واقعی به Go Server
                  // در محیط داکر، نام سرویس nexus-super-node است و پورت 3000
                  const goUrl = process.env.GO_SERVER_URL || 'http://nexus-super-node:3000';
                  
                  const goResponse = await fetch(`${goUrl}/internal/tools/deploy`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 
                      project_name: projectName, 
                      template: template || 'default' 
                    }),
                    signal: controller.signal
                  });

                  clearTimeout(timeoutId);

                  if (!goResponse.ok) {
                    const errorData = await goResponse.json().catch(() => ({}));
                    console.error(`[BFF] Go Server Error (${goResponse.status}):`, errorData);
                    return {
                      status: 'error',
                      message: `Failed to trigger deployment: ${errorData.message || goResponse.statusText}`
                    };
                  }

                  const result = await goResponse.json();
                  console.log('[BFF] Go Server Response:', result);
                  return result;

                } catch (error: any) {
                  console.error('[BFF] Error calling Go Server:', error);
                  if (error.name === 'AbortError') {
                    return {
                      status: 'error',
                      message: 'Deployment request timed out. The Go server took too long to respond.'
                    };
                  }
                  return {
                    status: 'error',
                    message: `Failed to connect to Go server: ${error.message}`
                  };
                }
              },
            }),
          },
        });

        // Chat v6 در فرانت، UI message stream انتظار دارد نه data stream خام
        const response = result.toUIMessageStreamResponse();
        
        // اضافه کردن هدرهای CORS به پاسخ استریم
        response.headers.set('Access-Control-Allow-Origin', '*');
        return response;

      } catch (error) {
        console.error('[BFF] Error:', error);
        return new Response(JSON.stringify({ error: 'Internal Server Error' }), { 
          status: 500,
          headers: { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' }
        });
      }
    }

    return new Response('Not Found', { status: 404 });
  },
});
