package esi

import (
	"context"
	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"log"
	"net/http"
	"neurobot/pkg/etagtripper"
	"time"
)

func init() {
	log.Println("initializing esi client")
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
}

var esiClient *goesi.APIClient

// EsiCharacterName searches ESI for char ID, gets char name
func EsiCharacterName(id int) string {
	c, _, err := esiClient.ESI.CharacterApi.GetCharactersCharacterId(context.Background(), int32(id), &esi.GetCharactersCharacterIdOpts{})
	if err != nil {
		log.Println("couldn't search character")
		return ""
	}
	log.Println("esiCharacterName:", c.Name)

	return c.Name
}
