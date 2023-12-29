package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"net/http"
	"neurobot/pkg/cartographer"
	"neurobot/pkg/esi"
	"neurobot/pkg/inventory"
	"neurobot/pkg/zkb"
	"time"
)

func killmailsKillfeed(km zkb.Killmail) error {
	// where to send
	channel := "1189353671213981798"
	e := generateKillmailKillfeed(km)

	_, err := Config.session.ChannelMessageSendEmbed(channel, &e)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending message to Discord: %+v", http.StatusInternalServerError))
	}

	return nil
}

func generateKillmailKillfeed(km zkb.Killmail) discordgo.MessageEmbed {

	// get final blow on km
	var finalblow zkb.Killmail
	for _, a := range km.Attackers {
		if a.FinalBlow == true {
			finalblow.Attackers = append(finalblow.Attackers, a)
			break
		}
	}

	// final blow base text format: "Final Blow", will append " (SOLO)" if solo kill, " (+N)" if not solo, but ignore any with character id of 0
	var finalblowTitleText string
	finalblowTitleText = "Final Blow"
	// switch statement here
	switch {
	case km.Zkb.Solo:
		finalblowTitleText += " (SOLO)"
		break
	case len(km.Attackers) > 1:
		atkrs := 0
		for _, a := range km.Attackers {
			if a.CharacterID == 0 {
				continue
			}
			atkrs++
		}
		finalblowTitleText += fmt.Sprintf(" (+%d)", atkrs)
		break
	}

	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    fmt.Sprintf("https://images.evetech.net/types/%d/render", km.Victim.ShipTypeID),
			Width:  64,
			Height: 64,
		},
		Title: fmt.Sprintf("Kill: %s (%s)\n%s", esi.EsiCharacter(km.Victim.CharacterID).Name, inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Zkb.URL),
		// todo: change this to a "### Did you know?" kind of random fact from a list of random neurotoxin facts
		//Description: fmt.Sprintf("[link to killmail](%s)", km.Zkb.URL),
		Color: 0xFF0000,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Time",
				Value:  fmt.Sprintf("`%s`\n<t:%d:R>", km.KillmailTime.Format(time.RFC822), km.KillmailTime.Unix()),
				Inline: true,
			},
			{
				Name:   "System",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/system/%d/)", cartographer.EveNavigation(0, km.SolarSystemID).Name, km.SolarSystemID),
				Inline: true,
			},
			{
				Name:   "Lost Ship",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/kill/%d/)", inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Victim.ShipTypeID),
				Inline: true,
			},
			{ // placeholder space
				Name:   "           ",
				Value:  "           ",
				Inline: true,
			},
			{
				Name:   "Victim",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/character/%d/)", esi.EsiCharacter(km.Victim.CharacterID).Name, km.Victim.CharacterID),
				Inline: true,
			},
			{
				Name: "Victim/Affiliation",
				Value: fmt.Sprintf("[%s](https://zkillboard.com/corporation/%d/)\n",
					esi.EsiCorporation(int(esi.EsiCharacter(km.Victim.CharacterID).CorporationId)).Name, km.Victim.CorporationID),
				Inline: true,
			},
			{
				Name:   "Value/Points",
				Value:  fmt.Sprintf("`%s ISK`\n`%d Points`", humanize.Comma(int64(km.Zkb.TotalValue)), km.Zkb.Points),
				Inline: true,
			},
			{
				Name: finalblowTitleText,
				Value: fmt.Sprintf("[%s](https://zkillboard.com/character/%d/)\n([%s](https://zkillboard.com/ship/%d/))",
					esi.EsiCharacter(finalblow.Attackers[0].CharacterID).Name, finalblow.Attackers[0].CharacterID,
					inventory.SdeGetItemName(finalblow.Attackers[0].ShipTypeID), finalblow.Attackers[0].ShipTypeID,
				),
				Inline: true,
			},
			{ //todo: handle faction somehow too
				Name:   "Killer/Affiliation",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/corporation/%d/)\n", esi.EsiCorporation(finalblow.Attackers[0].CorporationID).Name, finalblow.Attackers[0].CorporationID),
				Inline: true,
			},
		},
	}

	// update Victim/Alliance in embed.Fields
	if km.Victim.AllianceID != 0 {
		// update Victim/Alliance Value to placeholder
		embed.Fields[5].Value += fmt.Sprintf("[[%s](https://zkillboard.com/alliance/%d/)]", esi.EsiAlliance(km.Victim.AllianceID).Ticker, km.Victim.AllianceID)
	}

	// update Corp/Alliance if Alliance on finalblow character
	if finalblow.Attackers[0].AllianceID != 0 {
		embed.Fields[8].Value = embed.Fields[8].Value + fmt.Sprintf(" [[%s](https://zkillboard.com/alliance/%d/)]", esi.EsiAlliance(finalblow.Attackers[0].AllianceID).Ticker, finalblow.Attackers[0].AllianceID)
	}

	return embed
}
