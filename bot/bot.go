package bot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/danielelegbe/discord-join-count/bot/commands"
	"github.com/danielelegbe/discord-join-count/config"
	"github.com/danielelegbe/discord-join-count/sqlc"
)

func checkNilErr(e error) {
	if e != nil {
		slog.Error(e.Error())
	}
}

type Bot struct {
	Discord *discordgo.Session
	store   *sqlc.Queries
	ctx     context.Context
}

func New(discord *discordgo.Session, store *sqlc.Queries, ctx context.Context) *Bot {
	return &Bot{discord, store, ctx}
}

func (b *Bot) Run(botToken string) {
	b.Discord.AddHandler(b.HandleChannelJoinLeave)
	b.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		switch i.ApplicationCommandData().Name {
		case "zoomer-stats-individual":
			b.getUserStats(s, i)
		case "zoomer-stats-all":
			b.getAllUserStats(s, i)
		}
	})

	b.Discord.Open()

	registeredCommands, err := b.Discord.ApplicationCommandBulkOverwrite(config.ConfigInstance.AppId, config.ConfigInstance.GuildId, commands.Commands)
	if err != nil {
		slog.Error("Error registering commands", "error", err)
		return
	}

	slog.Info("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	for _, cmd := range registeredCommands {
		err := b.Discord.ApplicationCommandDelete(config.ConfigInstance.AppId, config.ConfigInstance.GuildId, cmd.ID)
		if err != nil {
			slog.Error("Error removing command", "command", cmd.Name, "error", err)
		}
	}
}

func (b *Bot) Close() {
	b.Discord.Close()
}
