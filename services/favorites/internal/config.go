package internal

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Options for the CLI.
type Options struct {
	AppEnv             string `help:"Application environment (production|local)" default:"local"`
	Port               int    `help:"Port to listen on" short:"p" default:"8001"`
	Auth0Domain        string `help:"Auth0 domain"`
	Auth0Audience      string `help:"Auth0 audience"`
	PostgresUrl        string `help:"Host for the Postgres database"`
	MetadataServiceUrl string `help:"URL for the metadata service" default:"http://localhost:8000"`
}

type Config struct {
	Options
}

func BuildConfig(opts Options) *Config {
	PrettyPrint(opts)
	return &Config{
		Options: opts,
	}
}

func PrettyPrint(v interface{}) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Println("Pretty Print Error:", err)
		return
	}
	logrus.Println(string(prettyJSON))
}
