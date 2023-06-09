package utils

import (
	"github.com/spf13/viper"
	"log"
)

// Config stores all the configuration for the application
// using values read from the config file or env variables
type Config struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	PORT         string `mapstructure:"PORT"`
	AuthSecret   string `mapstructure:"AUTH_SECRET"`
	ClientDomain string `mapstructure:"CLIENT_DOMAIN"`
}

var EnvVars Config

// LoadConfig loads the configuration from the config file or env variable
func LoadConfig() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Cannot read config file:", err)
	}
	viper.AutomaticEnv()

	EnvVars.DBDriver = viper.GetString("DB_DRIVER")
	EnvVars.DBSource = viper.GetString("DB_SOURCE")
	EnvVars.PORT = viper.GetString("PORT")
	EnvVars.AuthSecret = viper.GetString("AUTH_SECRET")
	EnvVars.ClientDomain = viper.GetString("CLIENT_DOMAIN")
}
