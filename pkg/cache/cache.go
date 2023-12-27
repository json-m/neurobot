package cache

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("created cache")
}

var cm = map[int]interface{}{} // Declare your cache map to hold any type of value.

// Set writes an int:key to any kind of value (interface{}) in the map
func Set(id int, val interface{}) {
	cm[id] = val
}

// Get retrieves an int:key value from the map and returns it as an interface{} type.
func Get(id int) (interface{}, bool) {
	val, ok := cm[id]
	return val, ok
}
