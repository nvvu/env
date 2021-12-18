package main

import (
	"fmt"
	"log"

	"github.com/nvvu/env"
)

type Config struct {
	HTTP   HTTP
	Kafka  *Kafka
	Secret string `env:"SECRET"`
}

type HTTP struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type Kafka struct {
	Brokers []string `env:"KAFKA_BROKERS"`
	Topic   string   `env:"KAFKA_TOPIC"`
	Group   string   `env:"KAFKA_GROUP"`
}

func main() {
	cfg := Config{}

	if err := env.OverwriteFromEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
	fmt.Println(cfg.Kafka)
}
