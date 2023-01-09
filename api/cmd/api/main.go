package main

import (
	"fmt"
	"net/http"

	api "github.com/intrntsrfr/vue-ws-test"
)

func main() {
	fmt.Println("setting up hub")
	hub := api.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", hub.Handler())

	port := ":8080"
	fmt.Println("running hub at port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
