package main

import (
	"github.com/yeremiaaryo/gotu-assignment/internal/apps/server"
	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"log"
)

func main() {
	var (
		cfg *configs.Config
	)

	err := configs.Init(
		configs.WithConfigFolder([]string{
			"./configs/",
			"./internal/configs/", // for local configs file path
		}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatalf("failed to initialize configs: %v", err)
	}
	cfg = configs.Get()

	err = server.InitApps(cfg)
	if err != nil {
		log.Fatalf("failed to initialize apps: %v", err)
	}
}
