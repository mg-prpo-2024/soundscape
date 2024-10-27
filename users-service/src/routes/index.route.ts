import { createRoute } from "@hono/zod-openapi";
import * as HttpStatusCodes from "stoker/http-status-codes";
import { jsonContent } from "stoker/openapi/helpers";
import { createMessageObjectSchema } from "stoker/openapi/schemas";

import { createRouter } from "@/lib/create-app";

const router = createRouter().openapi(
  createRoute({
    tags: ["Index"],
    method: "get",
    path: "/",
    summary: "Users API Index",
    responses: {
      [HttpStatusCodes.OK]: jsonContent(createMessageObjectSchema("Users API"), "Users API Index"),
    },
  }),
  (c) => {
    return c.json(
      {
        message: "Users API",
      },
      HttpStatusCodes.OK,
    );
  },
);

export default router;
