package inventory

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"github.com/ulikunitz/xz"
	"io"
	"log"
	"os"
	"runtime/debug"
)

// invItem represents an item in the sde
type invItems []struct {
	TypeID        int         `json:"typeID"`
	GroupID       int         `json:"groupID"`
	TypeName      string      `json:"typeName"`
	Description   string      `json:"description"`
	Mass          float64     `json:"mass"`
	Volume        int         `json:"volume"`
	Capacity      int         `json:"capacity"`
	PortionSize   int         `json:"portionSize"`
	RaceID        interface{} `json:"raceID"`
	BasePrice     interface{} `json:"basePrice"`
	Published     int         `json:"published"`
	MarketGroupID interface{} `json:"marketGroupID"`
	IconID        *int        `json:"iconID"`
	SoundID       *int        `json:"soundID"`
	GraphicID     int         `json:"graphicID"`
}

//go:embed invTypes.json.xz
var compressedinvTypes []byte
var invTypes []byte
var items invItems

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	// Decompress invTypes
	log.Println("decompressing invTypes.json.xz")
	r, err := xz.NewReader(bytes.NewReader(compressedinvTypes))
	if err != nil {
		log.Fatalf("Failed to create xz reader: %v", err)
	}

	invTypes, err = io.ReadAll(r)
	if err != nil {
		log.Fatalf("Failed to read xz compressed data: %v", err)
	}

	err = json.Unmarshal(invTypes, &items)
	if err != nil {
		// struct was autogenerated, might need type corrections later
		log.Fatalln("problem importing sde, probably a bad type:", err.Error())
	}

	// free memory
	log.Println("freeing memory")
	compressedinvTypes = nil
	invTypes = nil
	debug.FreeOSMemory()
}

func SdeGetItemName(id int) string {
	for _, item := range items {
		if item.TypeID == id {
			return item.TypeName
		}
	}
	return ""
}
