package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
	"time"
)

// Timer is an individual timer
type Timer struct {
	Owner       string    `json:"owner"`
	Message     string    `json:"message"`
	MessageID   string    `json:"messageID"`
	PinnedID    string    `json:"pinnedID"`
	Channel     string    `json:"channel"`
	Expires     time.Time `json:"expires"`
	Notify      string    `json:"notify"`
	HasNotified bool      `json:"hasNotified"`
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
					Owner:     m.Author.ID,
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

				// update response id
				timer.PinnedID = response.ID

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
				_ = s.MessageReactionAdd(m.ChannelID, m.ID, "❌")

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
	// if there are any timers
	if len(Config.Timers) > 0 {
		// generate the string for message content
		timers := "\n"
		for _, timer := range Config.Timers {
			tttt := strconv.FormatInt(timer.Expires.Unix(), 10)
			timers += fmt.Sprintf("* [%s](<https://discord.com/channels/%s/%s/%s>) :: <t:%s:R> :: <@%s> 📨 %s\n", timer.Message, m.GuildID, timer.Channel, timer.MessageID, tttt, timer.Owner, timer.Notify)
		}

		// inject embed into a message that allows mentions for timer.Notify
		msg := &discordgo.MessageSend{
			Content: timers,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: nil,
			},
		}

		// send
		_, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
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

// deleteTimer checks every 30 seconds for timers stored in Config.Timers to see if the timer Owner has added an X reaction to the MessageID
func deleteTimerHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	// check if bot is allowed to talk in this channel
	if blocked(m.ChannelID) {
		return
	}

	// is this MessageID a message in Config.Timers?
	for ti, timer := range Config.Timers {
		if timer.MessageID == m.MessageReaction.MessageID {
			//log.Println("got a reaction to a tracked message")
			// was the person who added the reaction the owner of the timer?
			if timer.Owner != m.UserID {
				return
			}
			if m.Emoji.Name == "❌" {
				log.Println("deleting timer:", m.UserID, m.MessageID, timer.Message)
				// delete the timer from Config.Timers, unpin the message, and add 🚮 emoji to m.MessageReaction.MessageID
				Config.Timers = append(Config.Timers[:ti], Config.Timers[ti+1:]...)
				err := writeConfig()
				if err != nil {
					log.Println("writing config during timer deletion:", err)
					return
				}
				err = s.ChannelMessageUnpin(m.ChannelID, timer.PinnedID)
				if err != nil {
					log.Println("Error unpinning message:", err)
					return
				}
				err = s.MessageReactionAdd(m.ChannelID, m.MessageReaction.MessageID, "🚮")
				if err != nil {
					log.Println("Error adding reaction:", err)
					return
				}
				break
			}
		}
	}
}

// sendTimerWarning
func sendTimerWarning(timer Timer) error {
	// create an embed
	newUnixTime := strconv.FormatInt(timer.Expires.Unix(), 10)
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("%s :: 30 minute warning!", timer.Message),
		Description: fmt.Sprintf("%s <t:%s:R>", timer.Message, newUnixTime),
		Color:       0xffa500, // Orange color
	}

	// inject embed into a message that allows mentions for timer.Notify
	msg := &discordgo.MessageSend{
		Content: fmt.Sprintf(timer.Notify),
		Embed:   embed,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{
				discordgo.AllowedMentionTypeRoles,
				discordgo.AllowedMentionTypeUsers,
			},
		},
	}

	// send
	_, err := Config.session.ChannelMessageSendComplex(timer.Channel, msg)
	if err != nil {
		return err
	}

	return nil
}
