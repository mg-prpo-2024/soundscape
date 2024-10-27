import { serve } from "@hono/node-server";

import app from "@/app";
import env from "@/env";
import { log } from "@/lib/log";

const port = env.PORT;

const server = serve(
  {
    fetch: app.fetch,
    port,
  },
  ({ address, port }) => {
    log.info(`Server listening at http://${address}:${port}`);
  },
);

const SHUTDOWN_TIMEOUT = env.NODE_ENV === "production" ? 30_000 : 1000;

async function gracefullShutdown() {
  log.info("Shutting down server...");
  server.close(() => {
    log.info("Server shut down successfully");
  });
  setTimeout(() => {
    log.error("Server could not close connections in time, forcefully shutting down");
    process.exit(1);
  }, SHUTDOWN_TIMEOUT);
}

process.on("SIGINT", gracefullShutdown);
process.on("SIGTERM", gracefullShutdown);
process.on("SIGQUIT", gracefullShutdown);
