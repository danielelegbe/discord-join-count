package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	DiscordToken     string `env:"DISCORD_TOKEN"`
	SPOST_CHANNEL_ID string `env:"SPOST_CHANNEL_ID"`
	MainChannelId    string `env:"MAIN_CHANNEL_ID"`
	UserId           string `env:"USER_ID"`
}

var ConfigInstance *Config

func GetConfig() *Config {

	ConfigInstance = &Config{
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		SPOST_CHANNEL_ID: os.Getenv("SPOST_CHANNEL_ID"),
		MainChannelId:    os.Getenv("MAIN_CHANNEL_ID"),
		UserId:           os.Getenv("USER_ID"),
	}

	err := verityConfig()

	if err != nil {
		log.Fatal(err)
	}

	return ConfigInstance
}

func verityConfig() error {
	if ConfigInstance.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is not set")
	}

	if ConfigInstance.SPOST_CHANNEL_ID == "" {
		return fmt.Errorf("SPOST_CHANNEL_ID is not set")
	}

	if ConfigInstance.MainChannelId == "" {
		return fmt.Errorf("MAIN_CHANNEL_ID is not set")
	}

	if ConfigInstance.UserId == "" {
		return fmt.Errorf("USER_ID is not set")
	}

	return nil
}
