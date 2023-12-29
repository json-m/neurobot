package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// cmdHandler switch statement to pass args in message to other command handler functions
func cmdHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// is me?
	if m.Author.ID == s.State.User.ID {
		return
	}

	// allowed to chat here?
	if blocked(m.ChannelID) {
		return
	}

	// talking to me?
	if !strings.Contains(m.Content, "<@1189348098695237662>") {
		return
	}

	// switch on each command
	args := strings.Split(stripCommand(m.Content), " ")
	switch args[0] {
	case "calc":
		timerCalcHandler(s, m)
	case "timer":
		timerHandler(s, m)
	case "timers":
		showTimersHandler(s, m)
	case "stats":
		botStats(s, m)
	default:
		return
	}
}
