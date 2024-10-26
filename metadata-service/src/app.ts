import configureOpenAPI from "@/lib/configure-open-api";
import createApp from "@/lib/create-app";
import artists from "@/routes/artists/artists.index";
import index from "@/routes/index.route";

const app = createApp();

configureOpenAPI(app);

const routes = [index, artists] as const;

routes.forEach((route) => {
  app.route("/", route);
});

export type AppType = (typeof routes)[number];

export default app;