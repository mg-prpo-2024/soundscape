import { pinoLogger as logger } from "hono-pino";

import { log } from "@/lib/log";

export function pinoLogger() {
  return logger({
    pino: log,
    http: {
      reqId: () => crypto.randomUUID(),
    },
  });
}
