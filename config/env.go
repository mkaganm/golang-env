package config

import "log"
import "github.com/spf13/viper"

var EnvConfigs *envConfigs

type envConfigs struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
}

func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

func loadEnvVariables() (config *envConfigs) {

	viper.AddConfigPath(".")

	viper.SetConfigName("app")
	//viper.SetConfigName("app1")
	//viper.SetConfigName("app2")

	viper.SetConfigType("env")
	//viper.SetConfigType("yaml")
	//viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic("Error reading config file! ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panic("Unable to decode into struct! ", err)
	}

	return
}
