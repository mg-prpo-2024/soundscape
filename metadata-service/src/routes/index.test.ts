import { testClient } from "hono/testing";
import { expect, it } from "vitest";

import app from "@/routes/index.route";

it("get / should return a description message for the api", async () => {
  const res = await testClient(app).index.$get();

  expect(await res.json()).toEqual({
    message: "Music Metadata API",
  });
  expect(res.status).toBe(200);
});
