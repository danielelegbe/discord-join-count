package commands

import "github.com/bwmarrin/discordgo"

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "zoomer-stats-individual",
			Description: "Get your voice channel statistics",
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
			Name:        "zoomer-stats-all",
			Description: "Show the voice channel leaderboard",
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
	}
)
