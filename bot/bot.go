package bot

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
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
	// add a event handler
	b.Discord.AddHandler(b.HandleChannelJoinLeave)

	// open session
	b.Discord.Open()

	// keep bot running untill there is NO os interruption (ctrl + C)
	slog.Info("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}

func (b *Bot) HandleChannelJoinLeave(discord *discordgo.Session, message *discordgo.VoiceStateUpdate) {
	name := message.Member.DisplayName()
	if message.BeforeUpdate == nil {
		fmt.Println("joining", name)

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

	if message.BeforeUpdate != nil && message.ChannelID == "" {
		fmt.Println("leaving", name)

		err := b.store.UpdateUserLeave(b.ctx, message.Member.User.ID)
		checkNilErr(err)
		return
	}
}

func (b *Bot) NewPersonStartedStreaming(discord *discordgo.Session, message *discordgo.VoiceStateUpdate) {
	userChannel, err := CreateUserChannel(discord)
	checkNilErr(err)

	if !message.BeforeUpdate.SelfStream {
		_, err := discord.ChannelMessageSend(userChannel.ID, fmt.Sprintf("%s started streaming", message.Member.DisplayName()))

		checkNilErr(err)
	} else if message.BeforeUpdate.SelfStream {

		_, err := discord.ChannelMessageSend(userChannel.ID, fmt.Sprintf("%s stopped streaming", message.Member.DisplayName()))

		checkNilErr(err)
	}

}

func CreateUserChannel(discord *discordgo.Session) (*discordgo.Channel, error) {
	userChannel, err := discord.UserChannelCreate(config.ConfigInstance.UserId)

	checkNilErr(err)

	return userChannel, nil
}

func (b *Bot) GetTableList(discord *discordgo.Session, message *discordgo.VoiceStateUpdate) {
	userChannel, err := CreateUserChannel(discord)
	users, err := b.store.ListUsers(b.ctx)
	checkNilErr(err)

	for _, user := range users {
		_, err := discord.ChannelMessageSend(userChannel.ID, fmt.Sprintf("%s joined the channel", user.Name))
		checkNilErr(err)
	}
}

func (b *Bot) Close() {
	b.Discord.Close()
}

func (b *Bot) GetUsersInVoiceChannel(channelID string) ([]*discordgo.Member, error) {
	// Get the channel details
	channel, err := b.Discord.Channel(channelID)

	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %v", err)
	}

	// Ensure the channel is a voice channel
	if channel.Type != discordgo.ChannelTypeGuildVoice {
		return nil, fmt.Errorf("channel %s is not a voice channel", channelID)
	}

	// Get the guild ID from the channel
	guildID := channel.GuildID

	// Fetch the guild to access VoiceStates
	guild, err := b.Discord.State.Guild(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild: %v", err)
	}

	// Collect user IDs from VoiceStates in the specified channel
	var userIDs []string
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channelID {
			userIDs = append(userIDs, vs.UserID)
		}
	}

	// Get full member details for each user
	var members []*discordgo.Member
	for _, userID := range userIDs {
		member, err := b.Discord.GuildMember(guildID, userID)
		if err != nil {
			log.Printf("failed to get member for user ID %s: %v", userID, err)
			continue
		}
		members = append(members, member)
	}

	return members, nil
}
