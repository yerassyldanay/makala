package configx

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Configuration struct {
	ListenAddr string `long:"port" env:"LISTEN_ADDR" description:"Listen to port (format: :8080|127.0.0.1:8080)" required:"false" default:"8080"`
	// PostgresHostname Postgres host address or name
	PostgresHostname string `long:"psql_host" env:"POSTGRESQL_HOST" description:"Postgres hostname" required:"false" default:"postgres"`
	// PostgresPort Postgres server port
	PostgresPort int64 `long:"psql_port" env:"POSTGRESQL_PORT" description:"Postgres port" required:"false" default:"5432"`
	// PostgresDBName Postgres server models name
	PostgresDBName string `long:"psql_db" env:"POSTGRESQL_DB" description:"Postgres models name" required:"false" default:"simple"`
	// PostgresUsername Postgres server user name
	PostgresUsername string `long:"psql_user" env:"POSTGRESQL_USER" description:"Postgres username" required:"false" default:"simple"`
	// PostgresPassword Postgres server password
	PostgresPassword string `long:"psql_password" env:"POSTGRESQL_PASSWORD" description:"Postgres password" required:"false" default:"simple"`
	PostgresSSLMode  string `long:"psql_ssl_mode" env:"POSTGRESQL_SSL_MODE" description:"Postgres SSL mode" required:"false" default:"disable"`
	// playground for trying and learning https://crontab.guru/
	// 0 * * * * - every hour
	// 0/1 * * * * - every minute
	CronSyncFeed string `long:"cron_sync_feed" env:"CRON_SYNC_FEED" description:"time period for feed update in in-memory" required:"false" default:"0/2 * * * *"`
	TimeFormat   string `long:"time_format" env:"TIME_FORMAT" description:"default time format" required:"false" default:"15:04:05 02/01/2006"`
	// redis host
	RedisHost              string `long:"redis_host" env:"REDIS_HOST" description:"redis host" required:"false" default:"redis"`
	RedisPort              int32  `long:"redis_port" env:"REDIS_PORT" description:"redis port" required:"false" default:"6379"`
	RedisFeedLogicDatabase int32  `long:"redis_feed_database" env:"REDIS_FEED_DATABASE" description:"a logic datastore for feed" required:"false" default:"0"`
	RedisAdLogicDatabase   int32  `long:"redis_ad_database" env:"REDIS_AD_DATABASE" description:"a logic datastore for promoted posts/ads" required:"false" default:"1"`
	GinMode                string `long:"gin_mode" env:"GIN_MODE" description:"set to debug mode to get more information" required:"false" default:"debug"`
	FeedSearchIntervalMin  int    `long:"feed_search_interval_min" env:"FEED_SEARCH_INTERVAL_MIN" description:"time that will probably required to scroll feed" required:"false" default:"1"`
	//  absolute path of the project
	absPath string
}

// NewConfiguration
// either parses env variables if it exists
// or parses default values if it is not required field (has required tag set to false)
func NewConfiguration() (*Configuration, error) {
	var c Configuration
	p := flags.NewParser(&c, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash|flags.IgnoreUnknown)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return nil, fmt.Errorf("this err indicates that the built-in help was shown (the error contains the help message). err: %w", err)
		} else {
			return nil, fmt.Errorf("failed to parse conf. err: %w", err)
		}
	}
	return &c, nil
}

// SetAbsPath sets the path of the project.
// It makes finding folder/file path easier
func (c *Configuration) SetAbsPath(p string) {
	c.absPath = p
}

// GetAbsPath returns the path of the project.
// It makes finding folder/file path easier
func (c *Configuration) GetAbsPath() string {
	return c.absPath
}
