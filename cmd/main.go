package main

import (
	"flag"
	"github.com/ashurbekovz/vktexbot/internal/app/build"
	"log"
)

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	secretPath := flag.String("secret", "", "Path to the secret file")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("config path is required")
		return
	}

	if *secretPath == "" {
		log.Fatal("secret path is required")
		return
	}

	app := build.NewVkApp(*configPath, *secretPath)
	app.Run()
}
