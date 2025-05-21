package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/ashurbekovz/vktexbot/internal/app/build"
)

func main() {
	configPath := flag.String("config", "", "Path to the config file")
	secretPath := flag.String("secret", "", "Path to the secret file")
	flag.Parse()

	config, secret, err := getPaths(*configPath, *secretPath)
	if err != nil {
		log.Fatal(err)
	}

	app := build.NewVkApp(config, secret)
	app.Run()
}

func getPaths(configFlag, secretFlag string) (string, string, error) {
	config := os.Getenv("CONFIG_PATH")
	if config == "" {
		config = configFlag
	}

	secret := os.Getenv("SECRET_PATH")
	if secret == "" {
		secret = secretFlag
	}

	if config == "" {
		return "", "", errors.New("config path is required (set via CONFIG_PATH env var or --config flag)")
	}

	if secret == "" {
		return "", "", errors.New("secret path is required (set via SECRET_PATH env var or --secret flag)")
	}

	return config, secret, nil
}
