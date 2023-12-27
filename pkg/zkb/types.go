package zkb

import "time"

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
