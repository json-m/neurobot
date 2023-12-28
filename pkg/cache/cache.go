package cache

import (
	"log"
	"os"
	"sync"
)

var (
	cm   = make(map[int]interface{})
	lock = sync.RWMutex{}
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("created cache")
}

// Set writes an int:key to any kind of value (interface{}) in the map
func Set(id int, val interface{}) {
	lock.Lock()
	defer lock.Unlock()

	cm[id] = val
}

// Get retrieves an int:key value from the map and returns it as an interface{} type.
func Get(id int) (interface{}, bool) {
	lock.RLock()
	defer lock.RUnlock()

	val, ok := cm[id]
	return val, ok
}

// Len prints len of cm as int
func Len() int {
	return len(cm)
}
