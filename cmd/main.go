package main

import (
	"flag"
	"log"
	"os"

	"github.com/ashurbekovz/vktexbot/internal/app/build"
)

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	secretPath := flag.String("secret", "", "Path to the secret file")
	flag.Parse()

	config := os.Getenv("CONFIG_PATH")
	if config == "" {
		config = *configPath
	}

	secret := os.Getenv("SECRET_PATH")
	if secret == "" {
		secret = *secretPath
	}

	if config == "" {
		log.Fatal("config path is required (set via CONFIG_PATH env var or --config flag)")
		return
	}

	if secret == "" {
		log.Fatal("secret path is required (set via SECRET_PATH env var or --secret flag)")
		return
	}

	app := build.NewVkApp(config, secret)
	app.Run()
}
