package internal

// Options for the CLI.
type Options struct {
	AppEnv              string `help:"Application environment (production|local)" default:"local"`
	Port                int    `help:"Port to listen on" short:"p" default:"8888"`
	PostgresUrl         string `help:"Host for the Postgres database"`
	Auth0Domain         string `help:"Auth0 domain"`
	Auth0Audience       string `help:"Auth0 audience"`
	Auth0HookSecret     string `help:"Auth0 webhook secret"`
	StripeSecretKey     string `help:"Stripe secret key"`
	StripeWebhookSecret string `help:"Stripe webhook secret"`
}

type Config struct {
	Options
	PlanStripeId PlanStripeId
}

type PlanStripeId struct {
	Student string
	Premium string
}

func BuildConfig(opts Options) *Config {
	return &Config{
		Options:      opts,
		PlanStripeId: GetStripePriceIds(opts.AppEnv),
	}
}

// hardcoded for now
func GetStripePriceIds(env string) PlanStripeId {
	if env == "local" {
		return PlanStripeId{
			Student: "price_1QdueRRxs8oYOJV19FDaJ3XG",
			Premium: "price_1QdudiRxs8oYOJV1tQ9sFvJ2",
		}
	}
	return PlanStripeId{
		Student: "price_1QdurZRrEPjnlg9fDMA3HYzP",
		Premium: "price_1Qdus0RrEPjnlg9frV08BsQt",
	}
}
