import { createRoute, z } from "@hono/zod-openapi";
import * as HttpStatusCodes from "stoker/http-status-codes";
import { jsonContent } from "stoker/openapi/helpers";

const tags = ["Artists"];

export const getOne = createRoute({
  method: "get",
  tags,
  path: "/artists/{id}",
  summary: "Get artist metadata by artist ID",
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
