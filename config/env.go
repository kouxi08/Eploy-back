package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Env() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}
}
