package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := readConfig()
	if err != nil {
		log.Fatal("error reading config file:", err)
	}
}

// Timer is an individual timer
type Timer struct {
	Message     string    `json:"message"`
	MessageID   string    `json:"messageID"`
	Channel     string    `json:"channel"`
	Expires     time.Time `json:"expires"`
	Notify      string    `json:"notify"`
	HasNotified bool      `json:"hasNotified"`
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
	Config.session.Identify.Intents = discordgo.IntentsGuildMessages

	// command handlers
	Config.session.AddHandler(timerHandler)
	Config.session.AddHandler(showTimersHandler)

	// Open a websocket connection to Discord and begin listening.
	err = Config.session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// start background tasks
	go timerMonitor()
	go httpListener()
	go func() { // periodic config writer
		for {
			time.Sleep(2 * time.Minute)
			bgerr := writeConfig()
			if bgerr != nil {
				_, _ = Config.session.ChannelMessageSend("1189353671213981798", "<@201538116664819712> i can't write config.json: "+bgerr.Error())
			}
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	Config.session.Close()
}

// timerHandler handles timer messages for calculating or adding a timer to track
func timerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if bot is allowed to talk in this channel
	if blocked(m.ChannelID) {
		return
	}

	if strings.Contains(m.Content, "<@1189348098695237662> timer ") {
		log.Println("in cmd:", m.Content)

		args := strings.Split(stripCommand(m.Content), " ")
		dd, hh, mm := processTimerInput(args[1])
		now := time.Now().UTC()

		// now nudge the now time by the time in the command, and return the new unixtime as string
		var newTime time.Time
		if hh == 0 && mm != 0 {
			newTime = now.Add(time.Minute * time.Duration(mm))
		} else if dd == 0 {
			newTime = now.Add(time.Hour*time.Duration(hh) + time.Minute*time.Duration(mm))
		} else {
			newTime = now.AddDate(0, 0, dd).Add(time.Hour*time.Duration(hh) + time.Minute*time.Duration(mm))
		}
		newUnixTime := strconv.FormatInt(newTime.Unix(), 10)

		// create embed of each type of timestamp message
		embed := &discordgo.MessageEmbed{
			Type:  "rich",
			Title: "Timer calculator",
			Description: fmt.Sprintf("<t:%s:F> :: <t:%s:R>\n\n", newUnixTime, newUnixTime) +
				fmt.Sprintf("Relative: `<t:%s:R>`\n", newUnixTime) +
				fmt.Sprintf("Full: `<t:%s:F>` \n", newUnixTime),
			Color: 0x0000ff,
		}

		// pin arg handler
		for _, a := range args {
			if a == "pin" {
				log.Println("got pin request")

				// the message format is:
				// timer ddhhmm pin group message
				// ddhhmm is arg 1

				// set notify and msg args
				notify := args[3]
				message := strings.Join(args[4:], " ")
				timer := Timer{
					Message:   message,
					MessageID: m.ID,
					Channel:   m.ChannelID,
					Expires:   newTime,
					Notify:    notify,
				}

				embed.Color = 0x00ff00
				embed.Title = message

				// send pin embed and get response for tracking
				response, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					log.Println("Error sending message:", err)
					return
				}

				// pin the message
				err = s.ChannelMessagePin(m.ChannelID, response.ID)
				if err != nil {
					log.Println("Error pinning message:", err)
					return
				}

				// append to timers
				Config.Timers = append(Config.Timers, timer)
				log.Printf("%+v", timer)

				// write config
				err = writeConfig()
				if err != nil {
					log.Println("writing config during pin:", err)
					_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}

				// react to msg
				_ = s.MessageReactionAdd(m.ChannelID, m.ID, "⏰")
				_ = s.MessageReactionAdd(m.ChannelID, m.ID, "👍")

				return
			}
		}

		// send normal timer calc message
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}

}

// showTimersHandler prints a list of timers that are currently being tracked
func showTimersHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if bot is allowed to talk in this channel
	if blocked(m.ChannelID) {
		return
	}

	// if command is timers
	if strings.Contains(m.Content, "<@1189348098695237662> timers") {
		log.Println("in cmd:", m.Content)

		// if there are any timers
		if len(Config.Timers) > 0 {
			// generate the string for message content
			timers := "\n"
			for _, timer := range Config.Timers {
				tttt := strconv.FormatInt(timer.Expires.Unix(), 10)
				timers += fmt.Sprintf("* %s :: <t:%s:R>\n", timer.Message, tttt)
			}

			// send
			_, err := s.ChannelMessageSend(m.ChannelID, timers)
			if err != nil {
				log.Println(err)
				return
			}

		} else { // otherwise if no timers
			_, err := s.ChannelMessageSend(m.ChannelID, "no timers running")
			if err != nil {
				log.Println(err)
				return
			}
		}

	}
}
