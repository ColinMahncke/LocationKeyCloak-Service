package main

import "os"

//fallback for different enviroment variables
func Getenv(name string, fallback string) string {
	variable, found := os.LookupEnv(name)

	if !found {
		variable = fallback
	}
	return variable
}
