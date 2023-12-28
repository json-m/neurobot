package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := readConfig()
	if err != nil {
		log.Fatal("error reading config file:", err)
	}
}

// mostly just public+whoops channel ids
var blockedChannels = []string{"617541872033726466", "658565710121009172", "709389086774919209", "617541872033726468", "958645719579893780", "985308569077440584"}

func main() {
	var err error
	Config.session, err = discordgo.New("Bot " + Config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	Config.session.Identify.Intents |= discordgo.IntentMessageContent | discordgo.IntentsGuildMessages

	// command handlers
	Config.session.AddHandler(cmdHandler)

	// todo: move these under cmdHandler
	Config.session.AddHandler(timerHandler)
	Config.session.AddHandler(deleteTimerHandler)

	// open ws connection
	err = Config.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// start background tasks
	go timerMonitor()
	go httpListener()

	// wait for ^C
	log.Println("startup done, bot should be up :: ^C to quit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	Config.session.Close()
}
