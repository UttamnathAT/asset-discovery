package main

import (
	"github.com/PAM-IDAM-Org/asset-discovery/internal/application"
)

func main() {
	app := application.New()
	
	if !app.Initialize() {
		return
	}
	
	if err := app.Run(); err != nil {
		panic(err)
	}
}
