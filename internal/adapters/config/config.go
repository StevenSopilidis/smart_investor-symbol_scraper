package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment        string        `mapstructure:"ENVIRONMENT""`
	AlphaVantageApiKey string        `mapstructure:"ALPHA_VANTAGE_API_KEY"`
	ServiceName        string        `mapstructure:"SERVICE_NAME"`
	SymbolServiceAddr  string        `mapstructure:"SYMBOL_SERVICE_ADDR"`
	OtelCollectorAddr  string        `mapstructure:"OTEL_COLLECTOR_ADDR"`
	ScrapeInterval     time.Duration `mapstructure:"SCRAPE_INTERVAL"`
	SymbolTopic        string        `mapstructure:"SYMBOL_TOPIC"`
	KafkaBroker        string        `mapstructure:"KAFKA_BROKER"`
	ScrapeEndpoint     string        `mapstructure:"SCRAPE_ENDPOINT"`
	SymbolManagerAddr  string        `mapstructure:"SYMBOL_MANAGER_ADDR"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("ALPHA_VANTAGE_API_KEY")
	viper.BindEnv("OTEL_COLLECTOR_ADDR")
	viper.BindEnv("SYMBOL_SERVICE_ADDR")
	viper.BindEnv("SERVICE_NAME")
	viper.BindEnv("SCRAPE_INTERVAL")
	viper.BindEnv("SYMBOL_TOPIC")
	viper.BindEnv("KAFKA_BROKER")
	viper.BindEnv("SCRAPE_ENDPOINT")
	viper.BindEnv("SYMBOL_MANAGER_ADDR")

	err = viper.Unmarshal(&config)
	return config, err
}
