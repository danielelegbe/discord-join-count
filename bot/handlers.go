package bot

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getUserStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	options := i.ApplicationCommandData().Options
	period := options[0].StringValue()

	user, err := b.store.GetUser(b.ctx, userID)

	userFound, err := HandleUserExists(err)
	checkNilErr(err)

	if !userFound {
		sendNotFoundResponse(s, i)
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("**Voice Statistics for %s (%s)**\n", user.Name, period))

	switch period {
	case "today":
		userStats, err := b.store.GetUserTodayTimeSpent(b.ctx, user.ID)

		checkNilErr(err)

		response.WriteString(fmt.Sprintf("Time spent today: %s\n", FormatNullIntDuration(userStats.MinutesToday)))

	case "week":
		userStats, err := b.store.GetUserWeeklyTimeSpent(b.ctx, user.ID)

		checkNilErr(err)

		response.WriteString(fmt.Sprintf("Time spent this week: %s\n", FormatNullIntDuration(userStats.MinutesThisWeek)))

	case "all":
		userStats, err := b.store.GetUserTotalTimeSpent(b.ctx, user.ID)

		checkNilErr(err)

		response.WriteString(fmt.Sprintf("Total time spent: %s\n", FormatNullIntDuration(userStats.TotalMinutes)))
		response.WriteString(fmt.Sprintf("Total sessions: %d\n", userStats.TotalJoins))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response.String(),
		},
	})
}

// HandleLeaderboard displays the voice channel leaderboard
func (b *Bot) getAllUserStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	period := options[0].StringValue()

	var response strings.Builder

	switch period {
	case "today":
		users, err := b.store.GetAllUsersTodayTimeSpent(b.ctx)
		checkNilErr(err)

		if len(users) == 0 {
			response.WriteString("No activity recorded for today!")
			return
		}

		response.WriteString(fmt.Sprintln("Total time spent today:"))

		for _, user := range users {
			response.WriteString(fmt.Sprintf("%s - %s\n", user.Name, FormatNullIntDuration(user.MinutesToday)))
		}
	case "week":
		users, err := b.store.GetAllUsersWeeklyTimeSpent(b.ctx)
		checkNilErr(err)

		if len(users) == 0 {
			response.WriteString("No activity recorded for this week!")
			return
		}

		response.WriteString(fmt.Sprintln("Total time spent this week:"))
		for _, user := range users {
			response.WriteString(fmt.Sprintf("%s - %s\n", user.Name, FormatNullIntDuration(user.MinutesThisWeek)))
		}
	case "all":
		users, err := b.store.GetAllTimeStats(b.ctx)
		checkNilErr(err)

		if len(users) == 0 {
			response.WriteString("No activity recorded")
			return
		}

		response.WriteString(fmt.Sprintln("Total time spent (all time):"))
		for _, user := range users {
			response.WriteString(fmt.Sprintf("%s - %s\n", user.Name, FormatNullIntDuration(user.TotalMinutes)))
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response.String(),
		},
	})
}

// FormatNullIntDuration converts duration to a readable string
func FormatNullIntDuration(minutes sql.NullFloat64) string {
	// Handle nullable values

	if minutes.Valid {
		return fmt.Sprintf("%dh %dm", int(minutes.Float64/60), int(minutes.Float64)%60)
	}

	return "0:00"
}

func FormatDuration(minutes int64) string {
	return fmt.Sprintf("%dh %dm", int(minutes/60), int(minutes)%60)
}

func sendNotFoundResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You haven't spent any time in voice channels yet!",
		},
	})
}
