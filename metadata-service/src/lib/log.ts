import pino from "pino";
import pretty from "pino-pretty";

import env from "@/env";

export const log = pino(
  {
    level: env.LOG_LEVEL || "info",
  },
  env.NODE_ENV === "production" ? undefined : pretty(),
);
