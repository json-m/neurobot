package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

type Conf struct {
	Token      string  `json:"token"`
	ListenPort int     `json:"port,omitempty"`
	Timers     []Timer `json:"timers,omitempty"`
	session    *discordgo.Session
}

var Config Conf

// load config file
func readConfig() error {
	// open config file
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("readConfig.open:", err)
	}

	// read config file
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("readConfig.Decode:", err)
	}
	f.Close()

	return nil
}

// write config file
func writeConfig() error {
	// open config file
	f, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("writeConfig.open:", err)
	}

	// write config file
	encoder := json.NewEncoder(f)
	err = encoder.Encode(&Config)
	if err != nil {
		log.Fatal("writeConfig.Encode:", err)
	}
	f.Close()

	log.Println("wrote to config")

	return nil
}
