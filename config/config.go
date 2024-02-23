package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config stores all the configuration for the application
// using values read from the config file or env variables
type Config struct {
	PORT string `mapstructure:"PORT"`
	ENV  string `mapstructure:"GIN_MODE"`

	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	RedisUrl string `mapstructure:"REDIS_URL"`

	CookieDomain       string `mapstructure:"COOKIE_DOMAIN"`
	AuthSecret         string `mapstructure:"AUTH_SECRET"`
	AuthCallbackURL    string `mapstructure:"AUTH_CALLBACK_URL"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	AWSRegion          string `mapstructure:"AWS_REGION"`
	AWSAccessKey       string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`

	EmailSenderAddress  string `mapstructure:"EMAIL_SENDER_ADDR"`
	EmailSenderPassword string `mapstructure:"EMAIL_SENDER_PASSWORD"`
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

	EnvVars.PORT = viper.GetString("PORT")
	EnvVars.ENV = viper.GetString("GIN_MODE")

	EnvVars.DBDriver = viper.GetString("DB_DRIVER")
	EnvVars.DBSource = viper.GetString("DB_SOURCE")
	EnvVars.RedisUrl = viper.GetString("REDIS_URL")

	EnvVars.CookieDomain = viper.GetString("COOKIE_DOMAIN")
	EnvVars.AuthSecret = viper.GetString("AUTH_SECRET")
	EnvVars.AuthCallbackURL = viper.GetString("AUTH_CALLBACK_URL")
	EnvVars.GoogleClientID = viper.GetString("GOOGLE_CLIENT_ID")
	EnvVars.GoogleClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")

	EnvVars.AWSRegion = viper.GetString("AWS_REGION")
	EnvVars.AWSAccessKey = viper.GetString("AWS_ACCESS_KEY_ID")
	EnvVars.AWSSecretAccessKey = viper.GetString("AWS_SECRET_ACCESS_KEY")

	EnvVars.EmailSenderAddress = viper.GetString("EMAIL_SENDER_ADDR")
	EnvVars.EmailSenderPassword = viper.GetString("EMAIL_SENDER_PASSWORD")
}
