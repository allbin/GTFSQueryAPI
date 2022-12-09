package main

import (
	"github.com/allbin/gtfsQueryGoApi/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app.Run()
}
