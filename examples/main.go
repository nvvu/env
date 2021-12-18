package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nvvu/env"
)

type Config struct {
	Api        ApiServer
	Username   string   `env:"AUTH_USER" yaml:"username"`
	Password   string   `env:"AUTH_PASS" yaml:"password"`
	Age        uint16   `env:"USER_AGE"`
	KafkaHosts []string `env:"KAFKA_HOSTS"`
	AA         []int    `env:"AA"`
}

type ApiServer struct {
	Host string   `json:"addr" env:"API_SERVER_HOST"`
	Port int      `json:"port" env:"API_SERVER_PORT"`
	X    []string `env:"API_SERVER_X"`
}

func main() {
	cfg := Config{}

	fmt.Println(cfg.Api)

	if err := env.OverwriteFromEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	d, _ := json.Marshal(cfg)
	fmt.Println(string(d))
}
