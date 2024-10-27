import { prometheus } from "@hono/prometheus";
import { OpenAPIHono } from "@hono/zod-openapi";
import { compress } from "hono/compress";
import { etag } from "hono/etag";
import { requestId } from "hono/request-id";
import { timeout } from "hono/timeout";
import { trimTrailingSlash } from "hono/trailing-slash";
import { notFound, onError, serveEmojiFavicon } from "stoker/middlewares";
import { defaultHook } from "stoker/openapi";

import { pinoLogger } from "@/middlewares/pino-logger";

import type { AppBindings, AppOpenAPI } from "./types";

const { printMetrics, registerMetrics } = prometheus();

export function createRouter() {
  return new OpenAPIHono<AppBindings>({
    strict: false,
    defaultHook,
  });
}

const DEFAULT_TIMEOUT = 20_000;

export default function createApp() {
  const app = createRouter();

  app.get("/health", (c) => {
    return c.text("OK");
  });
  // TODO: cors?
  app.use("*", registerMetrics);
  app.get("/metrics", printMetrics);
  app.use(trimTrailingSlash());
  app.use(requestId());
  app.use(pinoLogger());
  app.use(serveEmojiFavicon("🎤"));

  app.use(compress());
  app.use(etag());
  app.use(timeout(DEFAULT_TIMEOUT));

  app.notFound(notFound);
  app.onError(onError);

  return app;
}

export function createTestApp<R extends AppOpenAPI>(router: R) {
  return createApp().route("/", router);
}
