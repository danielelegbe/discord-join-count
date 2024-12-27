package bot

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/danielelegbe/discord-join-count/config"
	"github.com/danielelegbe/discord-join-count/sqlc"
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

func sendNotFoundResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "You haven't spent any time in voice channels yet!",
		},
	})
}

func (b *Bot) HandleChannelJoinLeave(discord *discordgo.Session, message *discordgo.VoiceStateUpdate) {
	// Don't track bots
	if message.Member.User.Bot {
		return
	}

	name := message.Member.DisplayName()

	if name == "" {
		name = message.Member.User.Username
	}

	// User has joined the channel
	if message.BeforeUpdate == nil {
		slog.Info("joined", "name", name)

		err := b.store.UpsertUser(b.ctx, sqlc.UpsertUserParams{
			ID:   message.Member.User.ID,
			Name: name,
		})

		checkNilErr(err)

		err = b.store.InsertUserJoin(b.ctx, sqlc.InsertUserJoinParams{
			UserID:    message.Member.User.ID,
			GuildID:   message.GuildID,
			ChannelID: message.ChannelID,
		})

		checkNilErr(err)

		return
	}

	// User has left the channel
	if message.BeforeUpdate != nil && message.ChannelID == "" {
		slog.Info("left", "name", name)

		err := b.store.UpdateUserLeave(b.ctx, message.Member.User.ID)
		checkNilErr(err)
	}
}

func CreateUserChannel(discord *discordgo.Session) (*discordgo.Channel, error) {
	userChannel, err := discord.UserChannelCreate(config.ConfigInstance.UserId)

	checkNilErr(err)

	return userChannel, nil
}

func (b *Bot) SendWeeklyLeaderboardScores(channelID string) error {
	var response strings.Builder
	users, err := b.store.GetAllUsersWeeklyTimeSpent(b.ctx)

	if err != nil {
		return err
	}

	response.WriteString(fmt.Sprint("Zoomerlympics Weekly Leaderboard 🏆\n\n"))

	if len(users) == 0 {
		response.WriteString("No activity recorded for this week!")
		return nil
	}

	for _, user := range users {
		response.WriteString(fmt.Sprintf("%s - %s - %d session(s)\n", user.Name, FormatNullIntDuration(user.MinutesThisWeek), user.JoinsThisWeek))
	}

	_, err = b.Discord.ChannelMessageSend(channelID, response.String())

	slog.Info("sending weekly leaderboards")

	return err
}
