package commands

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "stats",
			Description: "Get your voice channel statistics",
		},
		{
			Name:        "leaderboard",
			Description: "Show voice channel leaderboard",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "period",
					Description: "Time period (today/week/all)",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Today", Value: "today"},
						{Name: "This Week", Value: "week"},
						{Name: "All Time", Value: "all"},
					},
				},
			},
		},
		{
			Name:        "user",
			Description: "Get stats for a specific user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to get stats for",
					Required:    true,
				},
			},
		},
	}
)
