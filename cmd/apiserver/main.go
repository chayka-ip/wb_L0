package main

import (
	"level_zero/internal/app/apiserver"
	"log"
)

func main() {
	config := apiserver.LoadConfig()
	err := apiserver.Start(config)
	if err != nil {
		log.Fatal(err)
	}
}
