package main

import (
	"fmt"
	"net/http"
	"server/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Could not load the .env file")
	}

	server, err := server.NewServer()
	if err != nil {
		fmt.Println("Could not start the server:")
		fmt.Println(err.Error())
	}

	http.HandleFunc("/api/submit/", server.H.HandleSubmit)

	fmt.Println("Server started")
	err = http.ListenAndServe(":8080", nil)
}
