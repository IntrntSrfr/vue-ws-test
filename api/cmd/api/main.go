package main

import (
	"encoding/json"
	"fmt"
	"os"

	api "github.com/intrntsrfr/vue-ws-test"
)

type Config struct {
	JWTKey string `json:"jwt_key"`
}

func main() {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		panic("config file not found")
	}
	var config *Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic("mangled config file, fix it")
	}

	jwtUtil := api.NewJWTUtil([]byte(config.JWTKey))

	handler := api.NewHandler(&api.Config{JwtUtil: jwtUtil})

	err = handler.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
