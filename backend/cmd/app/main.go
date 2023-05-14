package main

import (
	"log"
	"net/http"
	"server/internal/server"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(color.RedString("FATAL ERROR: %s", err.Error()))
	}

	server := server.NewServer()

	http.HandleFunc("/api/submit/", server.H.HandleSubmit)

	log.Println(color.GreenString("SERVER STARTED"))
	http.ListenAndServe(":8080", nil)
}
