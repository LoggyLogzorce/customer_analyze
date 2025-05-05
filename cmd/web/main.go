package main

import (
	"first_static_analiz/internal/routers"
	"log"
)

func main() {
	router := routers.SetupRouter()

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
