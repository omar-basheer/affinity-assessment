package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// initialize router
	router := Router()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// initialize db
	if err := InitDB("voucher.db"); err != nil {
		log.Fatal("failed to initialize database: ", err)
	}

	// start server
	fmt.Println("Server running on port:", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
