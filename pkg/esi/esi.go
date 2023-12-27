package esi

import (
	"context"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"log"
	"net/http"
	"neurobot/pkg/cache"
	"neurobot/pkg/etagtripper"
	"time"
)

func init() {
	log.Println("initializing esi")
	// initalize a goesi client
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: &etagtripper.ETagTransport{Next: &http.Transport{}},
	}
	//_, err := httpcache.NewWithInmemoryCache(client, true, time.Hour*168)
	//if err != nil {
	//	log.Fatal(err)
	//}
	esiClient = goesi.NewAPIClient(client, "neurotoxin discord bot (ign: Maxx Ibanez / jason@xax.li")

	go ss()
}

var esiClient *goesi.APIClient

type stats struct {
	Hit     int `json:"hit"`
	Miss    int `json:"miss"`
	Lookups int `json:"lookups"`
}

var CS stats

func ss() {
	for {
		time.Sleep(3 * time.Minute)
		// print hit/miss and ratio based on lookup count
		log.Printf("CACHE: Hit: %d, Miss: %d, Ratio: %.2f\n", CS.Hit, CS.Miss, float64(CS.Hit)/float64(CS.Lookups))
	}
}

// EsiCharacter searches ESI for char ID, gets char struct
func EsiCharacter(id int) *esi.GetCharactersCharacterIdOk {
	if id == 0 {
		return &esi.GetCharactersCharacterIdOk{Name: "NPC", CorporationId: 0, AllianceId: 0}
	}
	CS.Lookups++
	// check if id exists in cache first, then lookup at esi
	n, f := cache.Get(id)
	if f == true {
		if char, ok := n.(*esi.GetCharactersCharacterIdOk); ok {
			CS.Hit++
			log.Println("EsiCharacter (HIT):", char.Name)
			return char
		}
	}

	c, _, err := esiClient.ESI.CharacterApi.GetCharactersCharacterId(context.Background(), int32(id), nil)
	if err != nil {
		log.Println("couldn't search character")
		return nil
	}
	log.Println("EsiCharacter (MISS):", c.Name)

	cache.Set(id, &c)
	CS.Miss++
	return &c
}

// EsiCorporation searches ESI for corporation ID and retrieves corporation struct
func EsiCorporation(id int) *esi.GetCorporationsCorporationIdOk {
	if id == 0 {
		return &esi.GetCorporationsCorporationIdOk{Name: "N/A", Ticker: "N/A", AllianceId: 0}
	}
	CS.Lookups++
	n, f := cache.Get(id)
	if f == true {
		if corp, ok := n.(*esi.GetCorporationsCorporationIdOk); ok {
			CS.Hit++
			log.Println("EsiCorporation (HIT):", corp.Name)
			return corp
		}
	}

	c, _, err := esiClient.ESI.CorporationApi.GetCorporationsCorporationId(context.Background(), int32(id), nil)
	if err != nil {
		log.Println("couldn't search corporation")
		return nil
	}
	log.Println("EsiCorporation (MISS):", c.Name)

	cache.Set(id, &c)
	CS.Miss++
	return &c
}

// EsiAlliance searches ESI for alliance ID and gets the alliance struct
func EsiAlliance(id int) *esi.GetAlliancesAllianceIdOk {
	if id == 0 {
		return &esi.GetAlliancesAllianceIdOk{Name: "N/A", Ticker: "N/A"}
	}
	CS.Lookups++
	// check if id exists in cache first, then lookup at esi
	n, f := cache.Get(id)
	if f == true {
		if alliance, ok := n.(*esi.GetAlliancesAllianceIdOk); ok {
			CS.Hit++
			log.Println("EsiAlliance (HIT):", alliance.Name)
			return alliance
		}
	}

	c, _, err := esiClient.ESI.AllianceApi.GetAlliancesAllianceId(context.Background(), int32(id), &esi.GetAlliancesAllianceIdOpts{})
	if err != nil {
		log.Println("couldn't search alliance")
		return nil
	}
	log.Println("EsiAlliance (MISS):", c.Name)

	cache.Set(id, &c)
	CS.Miss++
	return &c
}
