import { convertToModelMessages, streamText, tool } from 'ai';
import { google } from '@ai-sdk/google';
import { z } from 'zod';

const PORT = Number(process.env.PORT || 3001);
const GO_SERVER_URL = process.env.GO_SERVER_URL || 'http://nexus-super-node:3000';
const GO_TIMEOUT_MS = 10000;

type DeployWebsiteInput = {
  prompt: string;
  projectName?: string;
  framework?: string;
  theme?: string;
  template?: string;
};

type SelectedAgentContext = {
  id?: string;
  name?: string;
  type?: string;
  config?: {
    capability?: string;
    executionMode?: string;
    resultSurface?: string;
    tools?: string[];
    systemPrompt?: string;
  };
};

type WorkflowInsightInput = {
  question: string;
  workflowId?: string;
  status?: string;
  limit?: number;
};

type WorkflowRecord = {
  workflowId: string;
  name: string;
  status: string;
  currentStep: string;
  planningSource?: string;
  logs: string[];
  updatedAt?: string;
};

type WorkflowLogRecord = {
  workflowId: string;
  message: string;
  time?: string;
  status?: string;
  currentStep?: string;
};

type ChatRequestMessage = {
  role?: string;
  parts?: Array<{
    type?: string;
    text?: string;
  }>;
};

function slugifyProjectName(value: string) {
  const normalized = value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 40);

  return normalized || `site-${Date.now()}`;
}

function deriveProjectName(projectName: string | undefined, prompt: string) {
  if (projectName && projectName.trim().length > 0) {
    return slugifyProjectName(projectName);
  }

  const candidate = prompt
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 4)
    .join('-');

  return slugifyProjectName(candidate);
}

function getErrorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Unknown error';
}

async function fetchGoJson(path: string, init?: RequestInit) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), GO_TIMEOUT_MS);

  try {
    const response = await fetch(`${GO_SERVER_URL}${path}`, {
      ...init,
      signal: controller.signal,
    });
    const payload = await response.json().catch(() => ({}));

    if (!response.ok) {
      const message =
        typeof payload?.message === 'string'
          ? payload.message
          : typeof payload?.error === 'string'
            ? payload.error
            : response.statusText;
      throw new Error(message);
    }

    return payload;
  } catch (error: any) {
    if (error?.name === 'AbortError') {
      throw new Error('Request to Go gateway timed out.');
    }
    throw error;
  } finally {
    clearTimeout(timeoutId);
  }
}

function buildSelectedAgentPrompt(selectedAgent: SelectedAgentContext | null | undefined) {
  if (!selectedAgent) {
    return 'No draft agent is currently selected.';
  }

  const capability = selectedAgent.config?.capability || 'unknown';
  const executionMode = selectedAgent.config?.executionMode || 'unknown';
  const resultSurface = selectedAgent.config?.resultSurface || 'unknown';
  const tools = Array.isArray(selectedAgent.config?.tools)
    ? selectedAgent.config?.tools.join(', ')
    : 'none';

  return `Selected draft agent:
- name: ${selectedAgent.name || 'unknown'}
- type: ${selectedAgent.type || 'unknown'}
- capability: ${capability}
- executionMode: ${executionMode}
- resultSurface: ${resultSurface}
- tools: ${tools}`;
}

function getLatestUserText(messages: unknown) {
  if (!Array.isArray(messages)) {
    return '';
  }

  for (let i = messages.length - 1; i >= 0; i -= 1) {
    const message = messages[i] as ChatRequestMessage;
    if (message?.role !== 'user' || !Array.isArray(message.parts)) {
      continue;
    }

    return message.parts
      .filter((part) => part?.type === 'text')
      .map((part) => part.text ?? '')
      .join(' ')
      .trim();
  }

  return '';
}

function isDeployIntent(text: string) {
  return /deploy|website|site|landing|launch|build|create app|create website|سایت|وب.?سایت|لندینگ|دیپلوی|بساز/i.test(text);
}

function shouldForceWorkflowInsightTool(
  selectedAgent: SelectedAgentContext | null,
  latestUserText: string,
) {
  return selectedAgent?.config?.capability === 'workflow_insight' && !isDeployIntent(latestUserText);
}

console.log(`BFF Server is running on http://localhost:${PORT}`);

Bun.serve({
  port: PORT,
  async fetch(req) {
    const url = new URL(req.url);

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
      return Response.json(
        {
          status: 'ok',
          service: 'bff',
        },
        {
          headers: {
            'Access-Control-Allow-Origin': '*',
          },
        },
      );
    }

    if (url.pathname === '/api/chat' && req.method === 'POST') {
      try {
        const body = await req.json();
        console.log('[BFF] Received chat request body:', JSON.stringify(body, null, 2));
        const { messages, data } = body;

        const currentPath = data?.currentPath || 'unknown';
        const currentRoute = data?.currentRoute || currentPath;
        const selectedAgent = (data?.selectedAgent ?? null) as SelectedAgentContext | null;
        const latestUserText = getLatestUserText(messages);

        const systemPrompt = `You are VoltAgent, a helpful AI assistant for the Nexus platform. The user is currently on the page: ${currentPath} (route: ${currentRoute}). Use this context to provide more relevant answers if they ask about what they are seeing or what they can do here.

${buildSelectedAgentPrompt(selectedAgent)}

When the user wants to inspect workflow status, failures, stuck runs, latest runtime activity, or log-backed workflow summaries, call the workflow_insight tool.
For workflow_insight:
- Always pass the user's request text as question.
- Include workflowId only when the user clearly mentions a workflow identifier.
- Include status when the user asks for failed, running, completed, or similar workflow states.
- Keep limit small, usually between 1 and 5.

When the user wants to build, create, deploy, launch, or generate a website/app/landing page, call the deploy_website tool.
For deploy_website:
- Always pass the user's request text as prompt.
- Prefer a short lowercase projectName slug if the user provides one explicitly.
- If projectName is not explicit, you may omit it and the backend will derive one.
- Include framework, theme, and template only when the user expresses them or they are clear from context.

If a selected draft agent has capability workflow_insight, you must call workflow_insight for workflow, runtime, log, failure, status, or summary questions unless the user explicitly asks to deploy something new.
Do not answer with a generic limitation like "I cannot summarize workflows" when workflow_insight is available.
If the user's request is about workflows in Persian or English, treat it as a workflow_insight request by default.`;

        const coreMessages = await convertToModelMessages(messages);
        const toolChoice = shouldForceWorkflowInsightTool(selectedAgent, latestUserText)
          ? { type: 'tool' as const, toolName: 'workflow_insight' as const }
          : undefined;

        const result = streamText({
          model: google('gemini-2.5-flash'),
          system: systemPrompt,
          messages: coreMessages,
          toolChoice,
          tools: {
            workflow_insight: tool({
              description:
                'Inspect workflow executions and logs, then return a concise runtime insight summary.',
              inputSchema: z.object({
                question: z.string().describe('The user request about workflows, runtime status, failures, or logs'),
                workflowId: z.string().optional().describe('Optional workflow identifier when the user asks about a specific workflow'),
                status: z.string().optional().describe('Optional workflow status filter such as RUNNING, COMPLETED, or FAILED'),
                limit: z.number().int().min(1).max(5).optional().describe('Maximum number of workflows to include in the result'),
              }),
              execute: async ({ question, workflowId, status, limit }: WorkflowInsightInput) => {
                const normalizedQuestion = question.trim();
                const normalizedWorkflowId = workflowId?.trim();
                const normalizedStatus = status?.trim().toUpperCase();
                const normalizedLimit = Math.min(Math.max(limit ?? 3, 1), 5);

                if (!normalizedQuestion) {
                  return {
                    status: 'error',
                    capability: 'workflow_insight',
                    error: 'question_required',
                    message: 'A workflow insight question is required.',
                  };
                }

                try {
                  const executeResult = await fetchGoJson('/api/agents/execute', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                      question: normalizedQuestion,
                      workflowId: normalizedWorkflowId,
                      status: normalizedStatus,
                      limit: normalizedLimit,
                      selectedAgent: selectedAgent,
                      currentPath: currentPath,
                      currentRoute: currentRoute
                    }),
                  });
                  return executeResult;
                } catch (error) {
                  const errorMessage = getErrorMessage(error);
                  const isNotFound = errorMessage.toLowerCase().includes('not found');
                  
                  console.error('[BFF] Error calling workflow insight endpoint:', error);
                  return {
                    status: 'error',
                    capability: 'workflow_insight',
                    question: normalizedQuestion,
                    error: isNotFound ? 'workflow_not_found' : 'workflow_data_unavailable',
                    message: `Failed to retrieve workflow insight: ${errorMessage}`,
                  };
                }
              },
            }),
            deploy_website: tool({
              description: 'Deploy a new website or container',
              inputSchema: z.object({
                prompt: z.string().describe('The natural-language request describing the website or app to create'),
                projectName: z.string().optional().describe('Optional short slug for the project, e.g. ai-studio-site'),
                framework: z.string().optional().describe('Optional framework, e.g. svelte, react, vue'),
                theme: z.string().optional().describe('Optional theme such as modern, minimal, dark, light'),
                template: z.string().optional().describe('Optional template name to use for the deployment'),
              }),
              execute: async ({ prompt, projectName, framework, theme, template }: DeployWebsiteInput) => {
                const normalizedPrompt = prompt.trim();
                if (!normalizedPrompt) {
                  return {
                    status: 'error',
                    error: 'Prompt is required for deployment.',
                    message: 'Prompt is required for deployment.',
                  };
                }

                const normalizedProjectName = deriveProjectName(projectName, normalizedPrompt);
                console.log(
                  `[BFF] Tool called: deploy_website for ${normalizedProjectName}` +
                  `${framework ? ` using ${framework}` : ''}`,
                );

                try {
                  const deployResult = await fetchGoJson('/internal/tools/deploy', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                      project_name: normalizedProjectName,
                      prompt: normalizedPrompt,
                      framework: framework?.trim() || undefined,
                      theme: theme?.trim() || undefined,
                      template: template?.trim() || framework?.trim() || 'default',
                    }),
                  });

                  console.log('[BFF] Go Server Response:', deployResult);
                  return deployResult;
                } catch (error) {
                  console.error('[BFF] Error calling Go Server:', error);
                  return {
                    status: 'error',
                    error: 'deployment_unreachable',
                    message: `Failed to connect to Go server: ${getErrorMessage(error)}`,
                  };
                }
              },
            }),
          },
        });

        const response = result.toUIMessageStreamResponse();
        response.headers.set('Access-Control-Allow-Origin', '*');
        return response;
      } catch (error) {
        console.error('[BFF] Error:', error);
        return new Response(JSON.stringify({ error: 'Internal Server Error' }), {
          status: 500,
          headers: { 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' },
        });
      }
    }

    return new Response('Not Found', { status: 404 });
  },
});
