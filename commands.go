package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "push",
		Description: "Pushes the button",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "source",
		Description: "Returns a link to the project's GitHub repository",
		Type:        discordgo.ChatApplicationCommand,
	},
}

var handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"push": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Member == nil {
			err := sendError(s, i.Interaction, "You must be in the button channel to use this command")
			if err != nil {
				fmt.Println("Error sending error message: ", err)
			}
			return
		}

		if i.GuildID != os.Getenv("DISCORD_GUILD") {
			err := sendError(s, i.Interaction, "You must be in the button channel to use this command")
			if err != nil {
				fmt.Println("Error sending error message: ", err)
			}
			return
		}

		if i.ChannelID != os.Getenv("DISCORD_CHANNEL") {
			err := sendError(s, i.Interaction, "You must be in the button channel to use this command")
			if err != nil {
				fmt.Println("Error sending error message: ", err)
			}
			return
		}

		// update the user's role
		pushers <- i.Member.User.ID

		// send a message to the channel
		err := sendMessage(s, i.Interaction, "PUSHED!")
		if err != nil {
			fmt.Println(err)
		}
	},
	"source": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		sendMessage(s, i.Interaction, "https://github.com/Alextopher/discord-the-button")
	},
}

func sendError(s *discordgo.Session, i *discordgo.Interaction, msg string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         msg,
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           uint64(discordgo.MessageFlagsEphemeral),
			Files:           []*discordgo.File{},
			Choices:         []*discordgo.ApplicationCommandOptionChoice{},
			CustomID:        "",
			Title:           "",
		},
	})
}

func sendMessage(s *discordgo.Session, i *discordgo.Interaction, msg string) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			TTS:             false,
			Content:         msg,
			Components:      []discordgo.MessageComponent{},
			Embeds:          []*discordgo.MessageEmbed{},
			AllowedMentions: &discordgo.MessageAllowedMentions{},
			Flags:           0,
			Files:           []*discordgo.File{},
			Choices:         []*discordgo.ApplicationCommandOptionChoice{},
			CustomID:        "",
			Title:           "",
		},
	})
}
