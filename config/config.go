package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	ConfigurationPath = "./config"
	ConfigurationName = "config"
	ConfigurationType = "json"
)

type Configuration struct {
	MigrateToVersion          string
	MigrationLocation         string
	DisableSwaggerHttpHandler string
	GinMode                   string
	PostgreSQLUrl             string
	RabbitMQUrl               string
	AppName                   string
	AppVersion                string
	HttpPort                  string
	LogLevel                  string
	PgPoolMax                 string
	RmqRpcServer              string
	RmqRpcClient              string
}

// NewConfig returns app config.
func NewConfig() (*Configuration, error) {
	v := viper.New()
	v.AddConfigPath(ConfigurationPath)
	v.SetConfigName(ConfigurationName)
	v.SetConfigType(ConfigurationType)

	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("NewConfig: %v", err)
		return &Configuration{}, err
	}

	return &Configuration{
		MigrateToVersion:          v.Get("develop.migrateToVersion").(string),
		MigrationLocation:         v.Get("develop.migrationLocation").(string),
		DisableSwaggerHttpHandler: v.Get("develop.disableSwaggerHttpHandler").(string),
		GinMode:                   v.Get("develop.ginMode").(string),
		PostgreSQLUrl:             v.Get("develop.postgreSQLUrl").(string),
		RabbitMQUrl:               v.Get("develop.rabbitMQUrl").(string),
		AppName:                   v.Get("develop.appName").(string),
		AppVersion:                v.Get("develop.appVersion").(string),
		HttpPort:                  v.Get("develop.httpPort").(string),
		LogLevel:                  v.Get("develop.logLevel").(string),
		PgPoolMax:                 v.Get("develop.pgPoolMax").(string),
		RmqRpcServer:              v.Get("develop.rmqRpcServer").(string),
		RmqRpcClient:              v.Get("develop.rmqRpcClient").(string),
	}, err
}
