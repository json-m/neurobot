package main

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
					err := sendTimerWarning(t)
					if err != nil {
						log.Println("Error sending timer message:", err)
					}
					Config.Timers[i].HasNotified = true
					err = writeConfig()
					if err != nil {
						_, _ = Config.session.ChannelMessageSend("1189353671213981798", "<@201538116664819712> i can't write config.json: "+err.Error())
					}
				}

				// 48 hours after t.Expiry, remove it from the slice and unpin it from the channel
				if time.Since(t.Expires) > 48*time.Hour {
					log.Println("removing a timer because of 48hr expired removal")
					Config.Timers = append(Config.Timers[:i], Config.Timers[i+1:]...)
					i--
					err := writeConfig()
					if err != nil {
						_, _ = Config.session.ChannelMessageSend("1189353671213981798", "<@201538116664819712> i can't write config.json: "+err.Error())
					}
					err = Config.session.ChannelMessageUnpin(t.Channel, t.MessageID)
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

// blocked checks if the message source channel is in a blocklist, used for disallowing commands in public channels
func blocked(id string) bool {
	// if m.ChannelID is in blockedChannels, just return
	for _, channel := range blockedChannels {
		if id == channel {
			return true
		}
	}
	return false
}

// stripCommand cleans up the command
func stripCommand(arg string) string {
	return strings.TrimSpace(strings.TrimLeft(arg, "<@1189348098695237662>"))
}

func MemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc)) +
		fmt.Sprintf("\nTotalAlloc = %v MiB", bToMb(m.TotalAlloc)) +
		fmt.Sprintf("\nSys = %v MiB", bToMb(m.Sys)) +
		fmt.Sprintf("\nNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
