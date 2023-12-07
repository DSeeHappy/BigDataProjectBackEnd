package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	config := viper.New()

	config.SetConfigName(fileName)
	config.SetConfigType("toml")

	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")

	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}

	//if config.GetString("rabbitmq_instance_type") == "publisher" || config.GetString("rabbitmq_instance_type") == "PUBLISHER" {
	//	log.Println("Initializing RabbitMQ Publisher")
	//	server.InitRabbitMQPublisher()
	//} else {
	//	log.Println("Initializing RabbitMQ Subscriber")
	//	server.InitRabbitMQSubscriber()
	//}

	return config
}
