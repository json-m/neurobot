package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
	"time"
)

// find does a regex search and converts the result to an integer
func find(regex *regexp.Regexp, input string) int {
	match := regex.FindStringSubmatch(input)
	if len(match) == 0 {
		return 0
	}

	output, err := strconv.Atoi(match[1])
	if err != nil {
		fmt.Printf("Error in converting string to int: %s\n", err)
		return 0
	}

	return output
}

// processTimerInput takes the input string and returns the dd,hh,mm ints
func processTimerInput(input string) (int, int, int) {
	dayRegex := regexp.MustCompile(`(\d+)d`)
	hourRegex := regexp.MustCompile(`(\d+)h`)
	minuteRegex := regexp.MustCompile(`(\d+)m`)
	days := find(dayRegex, input)
	hours := find(hourRegex, input)
	minutes := find(minuteRegex, input)

	return days, hours, minutes
}

// timerMonitor is a background goroutine for checking on Timers in the Config to see if any are expiring soon or have expired
func timerMonitor() {
	for {
		if len(Config.Timers) > 0 {
			for i, t := range Config.Timers {
				// notify 30minutes before a timer expires, then update HasNotified for that timer so that it doesn't fire again
				if time.Until(t.Expires) <= 30*time.Minute && !t.HasNotified {
					log.Println("sending timer message:", t)
					err := sendTimerMessage(t)
					if err != nil {
						log.Println("Error sending timer message:", err)
					}
					Config.Timers[i].HasNotified = true
					_ = writeConfig()
				}

				// 48 hours after a timer has expired, remove it from the slice and unpin it from the channel
				if time.Until(t.Expires) >= 48*time.Hour {
					Config.Timers = append(Config.Timers[:i], Config.Timers[i+1:]...)
					i--
					_ = writeConfig()
					err := Config.session.ChannelMessageUnpin(t.Channel, t.MessageID)
					if err != nil {
						log.Println("Error unpinning message:", err)
					}
				}

				// do another thing..

			}
		}
		time.Sleep(time.Minute)
	}

}

func sendTimerMessage(timer Timer) error {
	// create an embed
	newUnixTime := strconv.FormatInt(timer.Expires.Unix(), 10)
	embed := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       fmt.Sprintf("%s :: 30 minute warning!", timer.Message),
		Description: fmt.Sprintf("%s\n%s\n<t:%s:R>", timer.Notify, timer.Message, newUnixTime),
		Color:       0xffa500, // Orange color
	}

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

	_, err := Config.session.ChannelMessageSendComplex(timer.Channel, msg)
	if err != nil {
		return err
	}

	return nil
}

func blocked(id string) bool {
	// if m.ChannelID is in blockedChannels, just return
	for _, channel := range blockedChannels {
		if id == channel {
			return true
		}
	}
	return false
}
