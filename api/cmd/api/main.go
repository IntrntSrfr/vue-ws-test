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

	// dependencies
	db, err := api.Open("./data.json")
	if err != nil {
		panic(err)
	}
	defer func(db *api.JsonDB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	jwtUtil := api.NewJWTUtil([]byte(config.JWTKey))

	// server
	handler := api.NewHandler(&api.Config{JwtUtil: jwtUtil, DB: db})

	// run server
	// this will block
	err = handler.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
