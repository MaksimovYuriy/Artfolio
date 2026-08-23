package main

import (
	"log"

	"github.com/maksimovyuriy/artfolio/backend/internal/migrator"
)

func main() {
	if err := migrator.Run(); err != nil {
		log.Fatal(err)
	}
}
