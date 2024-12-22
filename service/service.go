package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/danielelegbe/discord-join-count/bot"
	"github.com/danielelegbe/discord-join-count/sqlc"
)

type Service struct {
	bot   *bot.Bot
	store *sqlc.Queries
	ctx   context.Context
}

func New(bot *bot.Bot, store *sqlc.Queries, ctx context.Context) *Service {
	return &Service{
		bot:   bot,
		store: store,
		ctx:   ctx,
	}
}

func (s *Service) CreateLeaderboardScores() (string, error) {
	stats, err := s.store.GetAllTimeStats(s.ctx)

	if err != nil {
		return "", err
	}

	headers := []string{"User", "Total Time", "Total Joins"}
	rows := make([][]string, len(stats))
	for i, stat := range stats {

		// Handle nullable values
		minutes := int64(0)
		if stat.TotalMinutes.Valid {
			minutes = int64(stat.TotalMinutes.Float64)
		}

		joins := stat.TotalJoins // COUNT should never return NULL

		rows[i] = []string{
			stat.Name,
			bot.FormatDuration(minutes),
			fmt.Sprintf("%d", joins),
		}

	}

	// Format the table using the helper function
	return formatTable(headers, rows), nil
}

func (s *Service) SendLeaderboardScores(channelID string) error {
	scores, err := s.CreateLeaderboardScores()

	if err != nil {
		return err
	}

	message := "Zoomer Leaderboards 🏆"
	// add the scores below this message
	message += "\n\n"
	message += scores

	_, err = s.bot.Discord.ChannelMessageSend(channelID, message)

	return err
}

// Formats the headers and data rows into a table with proper alignment.
func formatTable(headers []string, rows [][]string) string {
	// Calculate column widths based on headers and rows
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	// Build the formatted table
	var builder strings.Builder
	builder.WriteString("```\n") // Start Discord code block

	// Format the headers
	for i, header := range headers {
		builder.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], header))
	}
	builder.WriteString("\n")

	// Format the rows
	for _, row := range rows {
		for i, col := range row {
			builder.WriteString(fmt.Sprintf("%-*s  ", colWidths[i], col))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("```") // End Discord code block
	return builder.String()
}
