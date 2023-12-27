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

var cs stats

func ss() {
	for {
		time.Sleep(3 * time.Minute)
		// print hit/miss and ratio based on lookup count
		log.Printf("CACHE: Hit: %d, Miss: %d, Ratio: %.2f\n", cs.Hit, cs.Miss, float64(cs.Hit)/float64(cs.Lookups))
	}
}

// EsiCharacter searches ESI for char ID, gets char struct
func EsiCharacter(id int) *esi.GetCharactersCharacterIdOk {
	if id == 0 {
		return &esi.GetCharactersCharacterIdOk{Name: "N/A", CorporationId: 0, AllianceId: 0}
	}
	cs.Lookups++
	// check if id exists in cache first, then lookup at esi
	n, f := cache.Get(id)
	if f != true {
		if char, ok := n.(*esi.GetCharactersCharacterIdOk); ok {
			cs.Hit++
			return char
		}
	}

	c, _, err := esiClient.ESI.CharacterApi.GetCharactersCharacterId(context.Background(), int32(id), nil)
	if err != nil {
		log.Println("couldn't search character")
		return nil
	}
	log.Println("esiCharacterName:", c.Name)

	cache.Set(id, c)
	cs.Miss++
	return &c
}

// EsiCorporation searches ESI for corporation ID and retrieves corporation struct
func EsiCorporation(id int) *esi.GetCorporationsCorporationIdOk {
	if id == 0 {
		return &esi.GetCorporationsCorporationIdOk{Name: "N/A", Ticker: "N/A", AllianceId: 0}
	}
	cs.Lookups++
	n, f := cache.Get(id)
	if f != true {
		if corp, ok := n.(*esi.GetCorporationsCorporationIdOk); ok {
			cs.Hit++
			return corp
		}
	}

	c, _, err := esiClient.ESI.CorporationApi.GetCorporationsCorporationId(context.Background(), int32(id), nil)
	if err != nil {
		log.Println("couldn't search corporation")
		return nil
	}
	log.Println("esiCorporationName:", c.Name)

	cache.Set(id, c)
	cs.Miss++
	return &c
}

// EsiAlliance searches ESI for alliance ID and gets the alliance struct
func EsiAlliance(id int) *esi.GetAlliancesAllianceIdOk {
	if id == 0 {
		return &esi.GetAlliancesAllianceIdOk{Name: "N/A", Ticker: "N/A"}
	}
	cs.Lookups++
	// check if id exists in cache first, then lookup at esi
	n, f := cache.Get(id)
	if f != true {
		if alliance, ok := n.(*esi.GetAlliancesAllianceIdOk); ok {
			cs.Hit++
			return alliance
		}
	}

	c, _, err := esiClient.ESI.AllianceApi.GetAlliancesAllianceId(context.Background(), int32(id), &esi.GetAlliancesAllianceIdOpts{})
	if err != nil {
		log.Println("couldn't search alliance")
		return nil
	}
	log.Println("esiAllianceName:", c.Name)

	cache.Set(id, c.Name)
	cs.Miss++
	return &c
}
