package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// verify all the roles needed for the button are present
// if any role is missing it is added
// function returns a slice of strings containing the ids of the roles
// purple - blue - green - yellow - orange - red
func verifyRoles(s *discordgo.Session, guildID string) {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		fmt.Println("Error getting roles: ", err)
		os.Exit(1)
	}

	// map of role names to role ids
	roleToButtonState = make(map[string]ButtonState)
	for _, role := range roles {
		roleToButtonState[role.ID] = ButtonStateFromRole(role.Name)
	}

	// if any role is missing it is added
	// function returns a slice of strings containing the ids of the roles
	// purple - blue - green - yellow - orange - red
	var hasRole map[ButtonState]struct{} = make(map[ButtonState]struct{})
	for _, state := range roleToButtonState {
		hasRole[state] = struct{}{}
	}

	for _, state := range ButtonStates {
		if _, ok := hasRole[state]; !ok {
			role, err := s.GuildRoleCreate(guildID)
			if err != nil {
				fmt.Println("Error creating role: ", err)
				os.Exit(1)
			}

			_, err = s.GuildRoleEdit(guildID, role.ID, state.Role(), state.Color(), true, 0, true)
			if err != nil {
				fmt.Println("Error editing role: ", err)
				os.Exit(1)
			}

			roleToButtonState[role.ID] = state
		}
	}
}

func updateUser(s *discordgo.Session, guildID, userID string, state ButtonState) error {
	// get the new role for the user
	role := state.Role()
	if role == "" {
		return fmt.Errorf("invalid state: %d", state)
	}

	// look up the role in the server
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return err
	}

	var roleID string
	for _, r := range roles {
		if r.Name == role {
			roleID = r.ID
			break
		}
	}

	// get the user's current roles
	user, err := s.GuildMember(guildID, userID)
	if err != nil {
		return err
	}

	// remove all color roles and check if the user has the red role
	hasRed := false
	for _, r := range user.Roles {
		if color, ok := roleToButtonState[r]; ok {
			if color == ButtonRed {
				hasRed = true
			}

			err := s.GuildMemberRoleRemove(guildID, userID, r)
			if err != nil {
				return err
			}
		}
	}

	// Red loses it's role if they push the button while it is red
	// this means red users can not keep the button alive without personal cost
	if hasRed && state == ButtonRed {
		return nil
	} else {
		return s.GuildMemberRoleAdd(guildID, userID, roleID)
	}
}

// sets the name of the button channel to match the state
func updateState(s *discordgo.Session, guildID, channelID string, state ButtonState) (time.Duration, error) {
	// track how long it takes to update the channel
	start := time.Now()
	fmt.Println("Changing state to ", state.Channel())
	_, err := s.ChannelEdit(channelID, state.Channel())
	fmt.Println("Changed!")
	return time.Since(start), err
}

// remove all color roles from all users
func endButton(s *discordgo.Session, guildID string) error {
	// TODO should this work with more than 1000 users?
	users, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, role := range user.Roles {
			if _, ok := roleToButtonState[role]; ok {
				err := s.GuildMemberRoleRemove(guildID, user.User.ID, role)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
