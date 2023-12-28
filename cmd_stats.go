package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"neurobot/pkg/cache"
	"neurobot/pkg/esi"
)

func botStats(s *discordgo.Session, m *discordgo.MessageCreate) {

	stats := fmt.Sprintf("```\nESI Cache: Hit: %d | Miss: %d | Ratio: %.2f | Size: %d objects\nActive Timers: %d\n\n\n%+v```\n", esi.CS.Hit, esi.CS.Miss, float64(esi.CS.Hit)/float64(esi.CS.Lookups), cache.Len(), len(Config.Timers), MemUsage())

	// send discord msg
	_, err := s.ChannelMessageSend(m.ChannelID, stats)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}
