# SoundScape

Run the app locally:

```bash
docker compose up --build
```

## Setup

### Auth0 configuration

Set up post login action that will call /sign-in endpoint of user service
The action needs:

- users service URL
- AUTH0_HOOK_SECRET

```ts
exports.onExecutePostLogin = async (event, api) => {
  const app = event.client.name;
  if (app !== "soundscape" && app !== "soundscape-prod") {
    return;
  }
  if (event.user.app_metadata.localUserCreated) {
    return;
  }
  const endpoint =
    app === "soundscape" ? "https://keen-wealthy-bengal.ngrok-free.app" : "http://72.144.96.197";
  const user = {
    email: event.user.email,
    id: event.user.user_id,
  };
  await fetch(`${endpoint}/users`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "X-Auth0-Webhook-Secret": event.secrets.AUTH0_HOOK_SECRET, // this is set in the auth0 dashboard
    },
    body: JSON.stringify(user),
  });
  api.user.setAppMetadata("localUserCreated", true);
};
```
