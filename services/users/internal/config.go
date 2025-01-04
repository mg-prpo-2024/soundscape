package internal

// Options for the CLI.
type Options struct {
	AppEnv              string `help:"Application environment (production|local)" default:"local"`
	Port                int    `help:"Port to listen on" short:"p" default:"8888"`
	PostgresUrl         string `help:"Host for the Postgres database"`
	Auth0Domain         string `help:"Auth0 domain"`
	Auth0Audience       string `help:"Auth0 audience"`
	StripeSecretKey     string `help:"Stripe secret key"`
	StripeWebhookSecret string `help:"Stripe webhook secret"`
}
