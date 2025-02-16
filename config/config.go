package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken     string `env:"DISCORD_TOKEN"`
	SPOST_CHANNEL_ID string `env:"SPOST_CHANNEL_ID"`
	MainChannelId    string `env:"MAIN_CHANNEL_ID"`
	UserId           string `env:"USER_ID"`
	AppId            string `env:"APP_ID"`
	GuildId          string `env:"GUILD_ID"`
}

var ConfigInstance *Config

func GetConfig() *Config {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	ConfigInstance = &Config{
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		SPOST_CHANNEL_ID: os.Getenv("SPOST_CHANNEL_ID"),
		MainChannelId:    os.Getenv("MAIN_CHANNEL_ID"),
		UserId:           os.Getenv("USER_ID"),
		AppId:            os.Getenv("APP_ID"),
		GuildId:          os.Getenv("GUILD_ID"),
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
