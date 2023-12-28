package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"neurobot/pkg/esi"
)

func botStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	stats := fmt.Sprintf("```\nCACHE: Hit: %d, Miss: %d, Ratio: %.2f\nActive Timers: %d\n```\n", esi.CS.Hit, esi.CS.Miss, float64(esi.CS.Hit)/float64(esi.CS.Lookups), len(Config.Timers))

	// send discord msg
	_, err := s.ChannelMessageSend(m.ChannelID, stats)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}
