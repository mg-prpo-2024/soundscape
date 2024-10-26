import { createRoute, z } from "@hono/zod-openapi";
import * as HttpStatusCodes from "stoker/http-status-codes";
import { jsonContent } from "stoker/openapi/helpers";

export const getOne = createRoute({
  tags: ["Artists"],
  method: "get",
  path: "/artists/{id}",
  responses: {
    [HttpStatusCodes.OK]: jsonContent(
      z.object({
        id: z.string().openapi({ example: "123" }),
      }),
      "Artist",
    ),
  },
});

export type GetOne = typeof getOne;
