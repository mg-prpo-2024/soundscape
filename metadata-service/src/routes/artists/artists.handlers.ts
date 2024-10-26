import type { AppRouteHandler } from "@/lib/types";
import type { GetOne } from "@/routes/artists/artists.routes";

export const getOne: AppRouteHandler<GetOne> = async (c) => {
  const id = c.req.param("id");
  return c.json({ id });
};
