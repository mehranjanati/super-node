import { randomUUID } from "node:crypto";

import { serve } from "@hono/node-server";
import { Hono } from "hono";
import { createPinoLogger } from "@voltagent/logger";
import { z } from "zod";

const port = Number(process.env.PORT ?? process.env.VOLTAGENT_PORT ?? 3141);
const serviceName = "voltagent-service";
const contractVersion = "v1alpha1";
type AppLogLevel = "fatal" | "error" | "warn" | "info" | "debug" | "trace";

const appLogLevels = new Set<AppLogLevel>([
  "fatal",
  "error",
  "warn",
  "info",
  "debug",
  "trace",
]);

const rawLogLevel = process.env.LOG_LEVEL;
const logLevel: AppLogLevel =
  rawLogLevel && appLogLevels.has(rawLogLevel as AppLogLevel)
    ? (rawLogLevel as AppLogLevel)
    : "info";

const logger = createPinoLogger({
  name: serviceName,
  level: logLevel,
});

const deployWebsiteInputSchema = z.object({
  project_name: z.string().trim().min(1, "project_name is required"),
  prompt: z.string().trim().min(1, "prompt is required"),
  framework: z.enum(["svelte", "react", "vue"]).optional(),
  theme: z.enum(["light", "dark", "modern", "minimal"]).optional(),
});

const planRequestSchema = z.object({
  contract_version: z.string(),
  intent: z.string(),
  input: deployWebsiteInputSchema,
  context: z
    .object({
      user_id: z.string().trim().min(1).optional(),
      session_id: z.string().trim().min(1).optional(),
      source: z.string().trim().min(1).optional(),
    })
    .optional(),
});

const planPayloadSchema = z.object({
  intent: z.literal("deploy_website"),
  kind: z.literal("workflow"),
  execution_target: z.literal("go-temporal"),
  workflow: z.object({
    name: z.literal("website-deployment-v1"),
    action: z.literal("start_dynamic_pipeline"),
    input: deployWebsiteInputSchema,
  }),
  artifacts: z.object({
    project_name: z.string().trim().min(1),
  }),
  warnings: z.array(z.string()),
});

type ErrorCode =
  | "VOLTAGENT_INVALID_REQUEST"
  | "VOLTAGENT_UNSUPPORTED_CONTRACT_VERSION"
  | "VOLTAGENT_UNSUPPORTED_INTENT"
  | "VOLTAGENT_PLAN_INVALID"
  | "VOLTAGENT_PLAN_GENERATION_FAILED"
  | "VOLTAGENT_INTERNAL_ERROR"
  | "VOLTAGENT_UNAVAILABLE";

type ErrorResponse = {
  status: "error";
  contract_version: typeof contractVersion;
  request_id: string;
  error: {
    code: ErrorCode;
    message: string;
    retryable: boolean;
    details?: Record<string, unknown>;
  };
};

const app = new Hono();

function getRequestID(headerValue?: string): string {
  const requestID = headerValue?.trim();
  return requestID && requestID.length > 0 ? requestID : randomUUID();
}

function jsonError(
  requestID: string,
  statusCode: number,
  code: ErrorCode,
  message: string,
  retryable: boolean,
  details?: Record<string, unknown>,
) {
  const body: ErrorResponse = {
    status: "error",
    contract_version: contractVersion,
    request_id: requestID,
    error: {
      code,
      message,
      retryable,
      ...(details ? { details } : {}),
    },
  };

  return Response.json(body, { status: statusCode });
}

app.get("/health", (c) => {
  const requestID = getRequestID(c.req.header("X-Request-Id"));

  logger.info("Health check completed", {
    request_id: requestID,
    route: "/health",
  });

  return c.json({
    status: "ok",
    service: serviceName,
    contract_version: contractVersion,
    request_id: requestID,
    checks: {
      api: "ok",
      planner: "ok",
    },
  });
});

app.post("/plan", async (c) => {
  const requestID = getRequestID(c.req.header("X-Request-Id"));
  const correlationID = c.req.header("X-Correlation-Id");
  let payload: unknown;

  try {
    payload = await c.req.json();
  } catch (error) {
    logger.warn("Received malformed JSON payload", {
      request_id: requestID,
      correlation_id: correlationID,
      route: "/plan",
      error,
    });

    return jsonError(
      requestID,
      400,
      "VOLTAGENT_INVALID_REQUEST",
      "Request body must be valid JSON",
      false,
    );
  }

  const parsedRequest = planRequestSchema.safeParse(payload);
  if (!parsedRequest.success) {
    const firstIssue = parsedRequest.error.issues[0];

    logger.warn("Plan request validation failed", {
      request_id: requestID,
      correlation_id: correlationID,
      route: "/plan",
      issues: parsedRequest.error.issues,
    });

    return jsonError(
      requestID,
      400,
      "VOLTAGENT_INVALID_REQUEST",
      firstIssue?.message ?? "Plan request is invalid",
      false,
      firstIssue?.path?.length
        ? { field: firstIssue.path.join(".") }
        : undefined,
    );
  }

  const request = parsedRequest.data;

  if (request.contract_version !== contractVersion) {
    return jsonError(
      requestID,
      400,
      "VOLTAGENT_UNSUPPORTED_CONTRACT_VERSION",
      `contract_version \`${request.contract_version}\` is not supported`,
      false,
      {
        contract_version: request.contract_version,
        supported_versions: [contractVersion],
      },
    );
  }

  if (request.intent !== "deploy_website") {
    return jsonError(
      requestID,
      400,
      "VOLTAGENT_UNSUPPORTED_INTENT",
      `intent \`${request.intent}\` is not supported in phase 1`,
      false,
      {
        intent: request.intent,
        supported_intents: ["deploy_website"],
      },
    );
  }

  try {
    const warnings: string[] = [];
    const normalizedInput = {
      ...request.input,
      framework: request.input.framework ?? "svelte",
      theme: request.input.theme ?? "minimal",
    };

    if (!request.input.framework) {
      warnings.push("framework defaulted to `svelte`");
    }

    if (!request.input.theme) {
      warnings.push("theme defaulted to `minimal`");
    }

    const validatedPlan = planPayloadSchema.safeParse({
      intent: "deploy_website",
      kind: "workflow",
      execution_target: "go-temporal",
      workflow: {
        name: "website-deployment-v1",
        action: "start_dynamic_pipeline",
        input: normalizedInput,
      },
      artifacts: {
        project_name: normalizedInput.project_name,
      },
      warnings,
    });

    if (!validatedPlan.success) {
      const firstIssue = validatedPlan.error.issues[0];

      logger.error("Generated plan failed contract validation", {
        request_id: requestID,
        correlation_id: correlationID,
        route: "/plan",
        issues: validatedPlan.error.issues,
      });

      return jsonError(
        requestID,
        422,
        "VOLTAGENT_PLAN_INVALID",
        firstIssue?.message ?? "Generated plan does not match the contract",
        false,
        firstIssue?.path?.length
          ? { field: firstIssue.path.join(".") }
          : undefined,
      );
    }

    logger.info("Plan request completed", {
      request_id: requestID,
      correlation_id: correlationID,
      route: "/plan",
      intent: request.intent,
      source: request.context?.source ?? "unknown",
    });

    return c.json({
      status: "ok",
      contract_version: contractVersion,
      request_id: requestID,
      plan: validatedPlan.data,
    });
  } catch (error) {
    logger.error("Unexpected plan generation failure", {
      request_id: requestID,
      correlation_id: correlationID,
      route: "/plan",
      error,
    });

    return jsonError(
      requestID,
      500,
      "VOLTAGENT_INTERNAL_ERROR",
      "Unexpected error while generating plan",
      true,
    );
  }
});

serve(
  {
    fetch: app.fetch,
    port,
  },
  (info) => {
    logger.info("VoltAgent internal API started", {
      port: info.port,
      contract_version: contractVersion,
    });
  },
);
