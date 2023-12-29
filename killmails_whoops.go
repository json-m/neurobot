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

func killmailsWhoops(km zkb.Killmail) error {
	// where to send
	channel := "985308569077440584"
	e := generateKillmailWhoops(km)

	_, err := Config.session.ChannelMessageSendEmbed(channel, &e)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending message to Discord: %+v", http.StatusInternalServerError))
	}

	return nil
}

func generateKillmailWhoops(km zkb.Killmail) discordgo.MessageEmbed {
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

	// loop to add as many poop emoji to desc for as many points in km.Zkb.Points
	var desc string
	for i := 0; i < km.Zkb.Points; i++ {
		desc += "💩"
	}

	// were they blobbed?
	blobbed := "❌"
	if len(km.Attackers) >= 5 {
		blobbed = "✅"
	}
	if len(km.Attackers) >= 20 {
		blobbed = "🤣"
	}
	if len(km.Attackers) >= 30 {
		blobbed = "<:kerdesk:964049663164567572>" // 964049663164567572
	}

	// was kill honorable? (combat recons maybe?)

	// other ideas for last field randomness

	// create a discord embed
	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    fmt.Sprintf("https://images.evetech.net/types/%d/render", km.Victim.ShipTypeID),
			Width:  64,
			Height: 64,
		},
		Title:       fmt.Sprintf("Feed: %s (%s)\n%s", esi.EsiCharacter(km.Victim.CharacterID).Name, inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Zkb.URL),
		Description: desc,
		Color:       0x654321,
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
				Name: finalblowTitleText,
				Value: fmt.Sprintf("[%s](https://zkillboard.com/ship/%d/)",
					inventory.SdeGetItemName(finalblow.Attackers[0].ShipTypeID), finalblow.Attackers[0].ShipTypeID,
				),
				Inline: true,
			},
			{ //todo: handle faction somehow too
				Name:   "Killer/Affiliation",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/corporation/%d/)\n", esi.EsiCorporation(finalblow.Attackers[0].CorporationID).Name, finalblow.Attackers[0].CorporationID),
				Inline: true,
			},
			{
				Name:   "Value",
				Value:  fmt.Sprintf("`%s ISK`", humanize.Comma(int64(km.Zkb.TotalValue))),
				Inline: true,
			},

			{
				Name:   "Victim",
				Value:  fmt.Sprintf("[%s](https://zkillboard.com/character/%d/)", esi.EsiCharacter(km.Victim.CharacterID).Name, km.Victim.CharacterID),
				Inline: true,
			},
			{ // todo: make this one a random one from a list of meme things
				Name:   "Blobbed?",
				Value:  blobbed,
				Inline: true,
			},
		},
	}

	return embed
}
