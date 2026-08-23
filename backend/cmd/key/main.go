package main

import (
	"log"

	"github.com/maksimovyuriy/artfolio/backend/internal/keygen"
)

func main() {
	if err := keygen.Run(); err != nil {
		log.Fatal(err)
	}
}
