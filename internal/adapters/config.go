package config

import "github.com/spf13/viper"

type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT""`
	ServiceName       string `mapstructure:"SERVICE_NAME"`
	SymbolServiceAddr string `mapstructure:"SYMBOL_SERVICE_ADDR"`
	OtelCollectorAddr string `mapstructure:"OTEL_COLLECTOR_STRUCTURE"`
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
	viper.BindEnv("OTEL_COLLECTOR_STRUCTURE")
	viper.BindEnv("SYMBOL_SERVICE_ADDR")
	viper.BindEnv("SERVICE_NAME")

	err = viper.Unmarshal(&config)
	return config, err
}
