package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	//TODO: Handle the received Killmail
	kmType := r.URL.Query()["type"][0] // for ie. POST /killmail?type=whoops

	w.WriteHeader(http.StatusOK)

	fmt.Println(r.URL.Query(), kmType)

}

// httpListener starts an http listener for the local api on port 9292
func httpListener() {
	http.HandleFunc("/killmail", kmReceiver)
	err := http.ListenAndServe(":9292", nil)
	if err != nil {
		go httpListener()
	}
}
