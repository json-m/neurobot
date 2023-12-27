package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"log"
	"net/http"
	"neurobot/pkg/cartographer"
	"neurobot/pkg/esi"
	"neurobot/pkg/inventory"
	"neurobot/pkg/zkb"
	"time"
)

// kmReceiver accepts a POST request with json data containing a Killmail
func kmReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var km zkb.Killmail

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&km)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	// for ie. POST /killmail?type=whoops
	kmType := r.URL.Query()["type"][0]
	log.Println("got km POST data:", kmType)

	if kmType == "testing" {
		channel := "1189353671213981798"

		// get final blow on km
		var finalblow zkb.Killmail
		for _, a := range km.Attackers {
			if a.FinalBlow == true {
				finalblow.Attackers = append(finalblow.Attackers, a)
				break
			}
		}

		// create a discord embed
		embed := discordgo.MessageEmbed{
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    fmt.Sprintf("https://images.evetech.net/types/%d/render", km.Victim.ShipTypeID),
				Width:  64,
				Height: 64,
			},
			Title: fmt.Sprintf("Killmail: %s (%s)\n%s", esi.EsiCharacter(km.Victim.CharacterID).Name, inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Zkb.URL),
			//Description: fmt.Sprintf("[link to killmail](%s)", km.Zkb.URL),
			Color: 0xFF0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Killmail Time",
					Value:  fmt.Sprintf("%s\n<t:%d:R>", km.KillmailTime.Format(time.RFC822), km.KillmailTime.Unix()),
					Inline: true,
				},
				{
					Name:   "System",
					Value:  fmt.Sprintf("[%s](https://zkillboard.com/system/%d/)", cartographer.EveNavigation(0, km.SolarSystemID).Name, km.SolarSystemID),
					Inline: true,
				},
				{
					Name:   "Ship",
					Value:  fmt.Sprintf("[%s](https://zkillboard.com/kill/%d/)", inventory.SdeGetItemName(km.Victim.ShipTypeID), km.Victim.ShipTypeID),
					Inline: true,
				},
				{
					Name:   "Victim",
					Value:  fmt.Sprintf("[%s](https://zkillboard.com/character/%d/)", esi.EsiCharacter(km.Victim.CharacterID).Name, km.Victim.CharacterID),
					Inline: true,
				},
				{
					Name:   "Victim/Corp",
					Value:  fmt.Sprintf("[%s](https://zkillboard.com/corporation/%d/)", esi.EsiCorporation(int(esi.EsiCharacter(km.Victim.CharacterID).CorporationId)).Name, km.Victim.CorporationID),
					Inline: true,
				},
				{
					Name:   "Victim/Alliance",
					Value:  "[N/A]",
					Inline: true,
				},
				{
					Name:   "Value",
					Value:  fmt.Sprintf("%s ISK (%dpts)", humanize.Comma(int64(km.Zkb.TotalValue)), km.Zkb.Points),
					Inline: true,
				},
				{
					Name: "Final Blow",
					Value: fmt.Sprintf("[%s](https://zkillboard.com/character/%d/) in a [%s](https://zkillboard.com/ship/%d/)",
						esi.EsiCharacter(finalblow.Attackers[0].CharacterID).Name, finalblow.Attackers[0].CharacterID,
						inventory.SdeGetItemName(finalblow.Attackers[0].ShipTypeID), finalblow.Attackers[0].ShipTypeID,
					),
					Inline: true,
				},
				{
					Name:   "",
					Value:  "",
					Inline: true,
				},
			},
		}

		// update Victim/Alliance in embed.Fields
		if km.Victim.AllianceID != 0 {
			// update Victim/Alliance Value to placeholder
			embed.Fields[5].Value = fmt.Sprintf("[[%s]](https://zkillboard.com/alliance/%d/)", esi.EsiAlliance(km.Victim.AllianceID).Ticker, km.Victim.AllianceID)
		}

		_, err = Config.session.ChannelMessageSendEmbed(channel, &embed)
		if err != nil {
			http.Error(w, "Error sending message to Discord", http.StatusInternalServerError)
			return
		}
	}

	if kmType == "whoops" {

	}

	w.WriteHeader(http.StatusOK)
}

// httpListener starts an http listener for the local api on port 9292
func httpListener() {
	http.HandleFunc("/killmail", kmReceiver)
	err := http.ListenAndServe(":9292", nil)
	if err != nil {
		go httpListener()
	}
}
