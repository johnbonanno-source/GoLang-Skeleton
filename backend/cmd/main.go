package main

import (
	"cartrader/backend/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed to load dotenv")
	}

	port := getEnvOrFail("PORT")
	dbConnStr := getEnvOrFail("DATABASE_CONNECTION")

	s, err := server.NewServer(dbConnStr, port)
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}

func getEnvOrFail(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s env var not set\n", key)
	}

	return val
}
