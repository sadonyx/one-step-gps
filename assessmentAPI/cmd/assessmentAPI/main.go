package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sadonyx/assessmentAPI/internal/routes"
	"github.com/sadonyx/assessmentAPI/internal/session"
	"github.com/sadonyx/assessmentAPI/internal/tokens"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tm := tokens.NewTokenManager()

	mongoUri := os.Getenv("MONGO_URI")
	sm, err := session.NewSessionManager(mongoUri)
	if err != nil {
		log.Fatalf("Failed to create session manager: %v", err)
	}

	router := routes.NewRouter(sm, tm)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		panic(err)
	}
}
