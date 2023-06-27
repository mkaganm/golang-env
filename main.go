package main

import (
	"fmt"
	"golang-env/config"
)

func main() {

	config.InitEnvConfigs()

	fmt.Println("SERVER ADDRESS : ", config.EnvConfigs.ServerAddress)
	fmt.Println("SECRET KEY : ", config.EnvConfigs.SecretKey)
	fmt.Println("DB DRIVER : ", config.EnvConfigs.DBDriver)
	fmt.Println("DB SOURCE : ", config.EnvConfigs.DBSource)
}
