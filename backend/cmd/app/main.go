package main

import (
	"fmt"
	"net/http"
	"server/internal/server"
)

func main() {

	server, err := server.NewServer()
	if err != nil {
		fmt.Println("Could not start the server:")
		fmt.Println(err.Error())
	}

	http.HandleFunc("/api/submit/", server.H.HandleSubmit)

	fmt.Println("Server started")
	err = http.ListenAndServe(":8080", nil)
}
