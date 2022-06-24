package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// A discord bot that updates a channel state every 6 hours through the colors of the rainbow
// button-🟣 | button-🔵 | button-🟢 | button-🟡 | button-🟠 | button-🔴
// when a user uses the "/push" command the button state resets and the user's role is updated to the color of the button when they pushed it
// if a button ever pasts "red" all users roles are removed and the button is reset

// a channel that users are added to when they push the button
var pushers chan string = make(chan string)

func main() {
	fmt.Println("Starting button bot...")

	// Load discord token from .env file
	godotenv.Load()
	token := os.Getenv("DISCORD_TOKEN")
	guild := os.Getenv("DISCORD_GUILD")
	channel := os.Getenv("DISCORD_CHANNEL")

	if token == "" {
		fmt.Println("Discord token not found in .env file")
		os.Exit(1)
	}

	// Create a new discord bot
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		os.Exit(1)
	}

	ch := make(chan struct{})
	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		fmt.Println("Connected to Discord")
		ch <- struct{}{}
	})

	// Handle application commands
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	session.Open()
	<-ch

	// verify all the roles needed for the button are present
	verifyRoles(session, guild)

	// Update the bot's interactions
	created, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", commands)
	if err != nil {
		fmt.Println("Error overwriting commands: ", err)
	} else {
		fmt.Printf("Overwrote %d commands\n", len(created))
	}

	// check if the button was already running and if we can pick up where we left off
	state := ButtonPurple
	c, err := session.Channel(channel)
	if err != nil {
		fmt.Println("Error getting channel: ", err)
		os.Exit(1)
	}

	if ButtonStateFromState(c.Name) != ButtonDead {
		state = ButtonStateFromState(c.Name)
		fmt.Println("countinuing from state: ", state.Role())
	}

	_, err = updateState(session, guild, channel, state)
	if err != nil {
		fmt.Println("Error updating state: ", err)
		os.Exit(1)
	}

	// run the bot forever
	timer := time.NewTimer(24 * time.Hour)
	for {
		select {
		case <-timer.C:
			state = state.Next()

			// if the button just died
			if state == ButtonDead {
				endButton(session, guild)

				// sleep for 1 week
				timer = time.NewTimer(7 * 24 * time.Hour)
			} else {
				t, err := updateState(session, guild, channel, state)
				if err != nil {
					fmt.Println("Error updating state: ", err)
				}

				// choose new timer duration between 12 hours and 24 hours
				wait := time.Duration(rand.Intn(12)+12) * time.Hour

				timer = time.NewTimer(wait - t)
			}
		case pusher := <-pushers:
			err = updateUser(session, guild, pusher, state)
			if err != nil {
				fmt.Println("Error updating user: ", err)
			} else {
				if state != ButtonPurple {
					state = ButtonPurple
					_, err := updateState(session, guild, channel, state)
					if err != nil {
						fmt.Println("Error updating state: ", err)
					}
				}

				timer.Stop()
				timer = time.NewTimer(24 * time.Hour)
			}
		}
	}
}
