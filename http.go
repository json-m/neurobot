package main

import (
	"encoding/json"
	"net/http"
	"neurobot/pkg/esi"
	"neurobot/pkg/zkb"
)

// kmReceiver accepts a POST request with json data containing a Killmail
func kmReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var km zkb.Killmail

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&km)
	if err != nil {
		http.Error(w, "Error decoding request", http.StatusInternalServerError)
		return
	}

	// for ie. POST /killmail?type=whoops
	kmType := r.URL.Query()["type"][0]
	//log.Println("got km POST data:", kmType)

	//if km.Victim.ShipTypeID != 670 {
	//	return
	//}

	// return if km.Zkb.TotalValue is 10000
	if km.Zkb.TotalValue == 1000000.0 {
		return
	}

	if kmType == "testing" {
		err = killmailsWhoops(km)
		if err != nil {
			http.Error(w, "Error sending message to Discord", http.StatusInternalServerError)
			return
		}
	}

	if kmType == "finishedkillfeed" {
		err = killmailsKillfeed(km)
		if err != nil {
			http.Error(w, "Error sending message to Discord", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// cacheStats prints the esi.CS struct as json
func cacheStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(esi.CS)
	if err != nil {
		return
	}
}

// httpListener starts an http listener for the local api on port 9292
func httpListener() {
	http.HandleFunc("/killmail", kmReceiver)
	http.HandleFunc("/cs", cacheStats)
	err := http.ListenAndServe(":9292", nil)
	if err != nil {
		go httpListener()
	}
}
