package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	_ "embed"

	"github.com/bwmarrin/discordgo"
	"github.com/danielelegbe/discord-join-count/bot"
	"github.com/danielelegbe/discord-join-count/config"
	"github.com/danielelegbe/discord-join-count/schedule"
	"github.com/danielelegbe/discord-join-count/service"
	"github.com/danielelegbe/discord-join-count/storage"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func main() {
	config := config.GetConfig()

	// Context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Make sure the data directory exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.MkdirAll("./data", os.ModePerm)
	}

	db, err := sql.Open("sqlite", "./data/discord.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer db.Close()

	// Initialize storage
	store := storage.CreateAndMigrateStore(db, ddl, ctx)

	// Initialize Discord bot
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	gocron, err := gocron.NewScheduler()

	if err != nil {
		log.Fatalf("Error creating GoCron scheduler: %v", err)
	}

	b := bot.New(discord, store, ctx)
	svc := service.New(b, store, ctx)

	scheduler := schedule.New(gocron, ctx, store, b, svc)

	go func() {
		slog.Info("Starting scheduler...")
		scheduler.Start()
	}()

	// Start the Discord bot on the main thread
	log.Println("Starting Discord bot...")
	b.Run(config.DiscordToken)

	// Wait for termination signal
	slog.Info("Bot and Scheduler running. Waiting for termination signal...")
	<-signalChan

	slog.Warn("Shutdown signal received. Cleaning up...")

	// Stop the scheduler and close the Discord session
	scheduler.Stop()
	b.Close()
	cancel()
}

func loadEnv() (string, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return "", errors.New("Warning: .env file not found")
	}

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		return "", errors.New("DISCORD_BOT_TOKEN is not set")
	}

	return botToken, nil
}
