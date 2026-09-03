package main

import (
	"log"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
