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
  if (event.user.app_metadata.localUserCreated) {
    return;
  }
  const endpoint = "https://75b6-93-103-238-188.ngrok-free.app";
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
