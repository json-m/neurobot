package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"math/rand"
	"net/http"
	"neurobot/pkg/cartographer"
	"neurobot/pkg/esi"
	"neurobot/pkg/inventory"
	"neurobot/pkg/zkb"
	"time"
)

// facts for kms
var facts = []string{
	"Neurotoxin Control was founded on July 21st 2019",
	"Neurotoxin have lost `4` Alliance Tournament ships",
	"Neurotoxin won `The Great War of WANGS` in Tama",
	"Neurotoxin provides femboy leasing at all major EVE Online events",
	"In 2020 <@228304135412383746> was appointed head diplomat of Neurotoxin Control?", // @'s liam
	"Since 2019 Neurotoxin has yet to get along with a single neighbor?",
	"Neurotoxin Control is a proud Triple Platinum sponsor of [Femboy Hooters](<https://zkillboard.com/corporation/98647355/>)",
	"The last Neurotoxin AT loss was <t:1680339240:R>",
	"In Neurotoxin, AWOXing is a rite of passage",
}

func killmailsKillfeed(km zkb.Killmail) error {
	// where to send
	channel := "658565710121009172"
	e := generateKillmailKillfeed(km)

	// disallow pings from facts
	msg := &discordgo.MessageSend{
		Embed: &e,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: nil,
		},
	}

	_, err := Config.session.ChannelMessageSendComplex(channel, msg)
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

	// select random string in facts
	randomFact := facts[rand.Intn(len(facts))]
	descStr := fmt.Sprintf("#### Did you know?\n%s", randomFact)

	// create embed
	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    fmt.Sprintf("https://images.evetech.net/types/%d/render", km.Victim.ShipTypeID),
			Width:  64,
			Height: 64,
		},
		Title:       fmt.Sprintf("Kill: %s (%s)\n%s", esi.EsiCharacter(km.Victim.CharacterID).Name, inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Zkb.URL),
		Description: descStr,
		Color:       0xFF0000,
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
