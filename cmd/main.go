package main

import "urlshortener/cmd/app"

func main() {
	app := app.NewApp()
	app.Run()
}
