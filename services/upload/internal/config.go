package internal

// Options for the CLI.
type Options struct {
	AppEnv                       string `help:"Application environment (production|local)" default:"local"`
	Port                         int    `help:"Port to listen on" short:"p" default:"8080"`
	Auth0Domain                  string `help:"Auth0 domain"`
	Auth0Audience                string `help:"Auth0 audience"`
	PostgresUrl                  string `help:"Host for the Postgres database"`
	AzureStorageConnectionString string `help:"Azure storage connection string"`
}

type Config struct {
	Options
}

func BuildConfig(opts Options) *Config {
	return &Config{
		Options: opts,
	}
}
