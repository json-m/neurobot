package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"neurobot/pkg/cartographer"
	"neurobot/pkg/esi"
	"neurobot/pkg/zkb"
	"strconv"
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

		// create a discord embed
		embed := discordgo.MessageEmbed{
			Title:       "Killmail",
			Description: fmt.Sprintf("[link to kill](https://zkillboard.com/kill/%d/)", km.KillmailID),
			Color:       0xFF0000,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Killmail ID",
					Value:  strconv.Itoa(km.KillmailID),
					Inline: true,
				},
				{
					Name:   "Killmail Time",
					Value:  km.KillmailTime.Format(time.RFC3339),
					Inline: true,
				},
				{
					Name:   "Solar System Name",
					Value:  cartographer.EveNavigation(0, km.SolarSystemID).Name,
					Inline: true,
				},
				{
					Name:   "Victim",
					Value:  esi.EsiCharacterName(km.Victim.CharacterID),
					Inline: true,
				},
			},
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
