package cartographer

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xyz"
	"github.com/ulikunitz/xz"
	"io"
	"log"
	"neurobot/pkg/zkb"
)

// SdeLocations is a slice of structs from the json file
type SdeLocations []struct {
	ItemID          int         `json:"itemID"`
	TypeID          int         `json:"typeID"`
	GroupID         int         `json:"groupID"`
	SolarSystemID   interface{} `json:"solarSystemID"`
	ConstellationID interface{} `json:"constellationID"`
	RegionID        interface{} `json:"regionID"`
	OrbitID         interface{} `json:"orbitID"`
	X               float64     `json:"x"`
	Y               float64     `json:"y"`
	Z               float64     `json:"z"`
	Radius          interface{} `json:"radius"`
	ItemName        string      `json:"itemName"`
	Security        interface{} `json:"security"`
	CelestialIndex  interface{} `json:"celestialIndex"`
	OrbitIndex      interface{} `json:"orbitIndex"`
}

//go:embed mapDenormalize.json.xz
var compressedMapDenormalize []byte
var mapDenormalize []byte
var sdeLocations SdeLocations

// load data for sdeLocations
func init() {
	var err error
	// Decompress mapDenormalize
	r, err := xz.NewReader(bytes.NewReader(compressedMapDenormalize))
	if err != nil {
		log.Fatalf("Failed to create xz reader: %v", err)
	}

	mapDenormalize, err = io.ReadAll(r)
	if err != nil {
		log.Fatalf("Failed to read xz compressed data: %v", err)
	}

	err = json.Unmarshal(mapDenormalize, &sdeLocations)
	if err != nil {
		// struct was autogenerated, might need type corrections later
		log.Fatalln("problem importing sde, probably a bad type:", err.Error())
	}
}

type Location struct {
	ID      int
	Name    string
	Spatial Coordinates
}

// is the location a stargate?
func (l Location) IsStargate() bool {
	// stargates range from 50000000 to 60000000
	return l.ID >= 50000000 && l.ID <= 60000000
}

// is the location a station?
func (l Location) IsStation() bool {
	// stations range from 60000000 to 64000000
	return l.ID >= 60000000 && l.ID <= 64000000
}

// gets index of location from sdeLocations for later use
func getSdeLocIndex(id int) int {
	if id == 0 {
		id = 50014002 // nourv gate in tama
	}
	for i, v := range sdeLocations {
		if v.ItemID == id {
			return i
		}
	}
	return -1
}

// takes index, returns coordinates
func getLocCoordinates(i int) geom.Coord {
	// return coordinates
	return geom.Coord{
		sdeLocations[i].X,
		sdeLocations[i].Y,
		sdeLocations[i].Z,
	}
}

// converts Killmail positional to a geom.Coord
func getKillmailCoordinates(km zkb.Killmail) geom.Coord {
	return geom.Coord{
		km.Victim.Position.X,
		km.Victim.Position.Y,
		km.Victim.Position.Z,
	}
}

// returns a distance as kilometers
func (l Location) Distance() float64 {
	return xyz.Distance(l.Spatial.CoordA, l.Spatial.CoordB) / 1000
}

func EveLocation(km zkb.Killmail) (l Location) {
	locIdx := getSdeLocIndex(km.Zkb.LocationID)
	return Location{
		ID:   km.Zkb.LocationID,
		Name: sdeLocations[locIdx].ItemName,
		Spatial: Coordinates{
			CoordA: getLocCoordinates(locIdx),
			CoordB: getKillmailCoordinates(km),
		},
	}
}