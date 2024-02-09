package util

import "github.com/spf13/viper"

// config stores all configuration of the application
// the value are read by viper from a config file or environment variables
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// loadConfig reads configuration from file path or environment
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
