package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Killmail struct {
	Attackers []struct {
		AllianceID     int     `json:"alliance_id"`
		CharacterID    int     `json:"character_id"`
		CorporationID  int     `json:"corporation_id"`
		FactionID      int     `json:"faction_id"`
		DamageDone     int     `json:"damage_done"`
		FinalBlow      bool    `json:"final_blow"`
		SecurityStatus float64 `json:"security_status"`
		ShipTypeID     int     `json:"ship_type_id"`
		WeaponTypeID   int     `json:"weapon_type_id"`
	} `json:"attackers"`
	KillmailID    int       `json:"killmail_id"`
	KillmailTime  time.Time `json:"killmail_time"`
	SolarSystemID int       `json:"solar_system_id"`
	Victim        struct {
		AllianceID    int `json:"alliance_id"`
		CharacterID   int `json:"character_id"`
		CorporationID int `json:"corporation_id"`
		FactionID     int `json:"faction_id"`
		DamageTaken   int `json:"damage_taken"`
		Items         []struct {
			Flag              int `json:"flag"`
			ItemTypeID        int `json:"item_type_id"`
			QuantityDestroyed int `json:"quantity_destroyed,omitempty"`
			Singleton         int `json:"singleton"`
			QuantityDropped   int `json:"quantity_dropped,omitempty"`
		} `json:"items"`
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
		} `json:"position"`
		ShipTypeID int `json:"ship_type_id"`
	} `json:"victim"`
	Zkb struct {
		LocationID  int     `json:"locationID"`
		Hash        string  `json:"hash"`
		FittedValue float64 `json:"fittedValue"`
		TotalValue  float64 `json:"totalValue"`
		Points      int     `json:"points"`
		Npc         bool    `json:"npc"`
		Solo        bool    `json:"solo"`
		Awox        bool    `json:"awox"`
		Esi         string  `json:"esi"`
		URL         string  `json:"url"`
	} `json:"zkb"`
}

// kmReceiver accepts a POST request with json data containing a Killmail
func kmReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var km *Killmail

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
			Description: "New killmail received",
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
					Name:   "Solar System ID",
					Value:  strconv.Itoa(km.SolarSystemID),
					Inline: true,
				},
				{
					Name:   "Victim",
					Value:  strconv.Itoa(km.Victim.CharacterID),
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
