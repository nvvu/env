package main

import (
	"fmt"
	"log"

	"github.com/nvvu/env"
)

type Config struct {
	Api      ApiServer
	Username string `env:"AUTH_USER" yaml:"username"`
	Password string `env:"AUTH_PASS" yaml:"password"`
	Age      uint16 `env:"USER_AGE"`
}

type ApiServer struct {
	Host string `json:"addr" env:"API_SERVER_HOST"`
	Port int    `json:"port" env:"API_SERVER_PORT"`
}

func main() {
	cfg := Config{
		Username: "user_1",
		Password: "passss",
		Api: ApiServer{
			Host: "0.0.0.0",
			Port: 3000,
		},
	}

	if err := env.OverwriteFromEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}
