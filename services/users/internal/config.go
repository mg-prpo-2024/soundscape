package internal

// Options for the CLI.
type Options struct {
	AppEnv              string `help:"Application environment (production|local)" default:"local"`
	Port                int    `help:"Port to listen on" short:"p" default:"8888"`
	PostgresHost        string `help:"Host for the Postgres database" default:"localhost"`
	PostgresUser        string `help:"Username for the Postgres database" default:"5432"`
	PostgresPassword    string `help:"Password for the Postgres database"`
	PostgresDB          string `help:"DB name for the Postgres database"`
	PostgresPort        int    `help:"Port for the Postgres database" default:"5432"`
	Auth0Domain         string `help:"Auth0 domain"`
	Auth0Audience       string `help:"Auth0 audience"`
	StripeSecretKey     string `help:"Stripe secret key"`
	StripeWebhookSecret string `help:"Stripe webhook secret"`
}
