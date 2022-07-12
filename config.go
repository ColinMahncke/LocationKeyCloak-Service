package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load() error {
	err := LoadEnv()
	if err != nil {
		return err

	}
	return nil
}

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")

		return err
	}
	return nil
}
