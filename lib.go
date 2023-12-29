package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
