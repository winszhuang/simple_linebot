package config

import "github.com/spf13/viper"

type Config struct {
	GoogleMapApiKey string `mapstructure:"GOOGLE_MAP_API_KEY"`

	ChannelSecret string `mapstructure:"CHANNEL_SECRET"`
	ChannelToken  string `mapstructure:"CHANNEL_TOKEN"`

	DBUserName string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBDomain   string `mapstructure:"DB_DOMAIN"`
	DBName     string `mapstructure:"DB_NAME"`

	WebHookUrl string `mapstructure:"WEBHOOK_URL"`
	Port       string `mapstructure:"PORT"`
}

var C *Config

func Init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		panic(err)
	}
}
