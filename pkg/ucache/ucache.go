package ucache

// μcache - absolutely diminutive in-memory caching
// just kidding, it's just slices all the way down.

import (
	"errors"
	"log"
	"time"
)

var c Cache

func init() {
	go c.startExpiryThread()
}

// cache structure, errors todo: value should be interface{} or something to be more generic
type Record struct {
	Value   int       `json:"value"`
	Created time.Time `json:"created"`
}
type Cache struct {
	cache  []Record
	Expiry time.Duration
}

var ErrNotFound = errors.New("not found")
var ErrRecordExists = errors.New("record exists")

// get record
func (c *Cache) Get(val int) (Record, error) {
	for _, r := range c.cache {
		if r.Value == val {
			return r, nil
		}
	}
	// wasn't found in slice
	return Record{}, ErrNotFound
}

// sets record
func (c *Cache) Set(val int) error {
	_, err := c.Get(val)
	if errors.Is(err, ErrNotFound) {
		// if it doesn't exist, then set it
		r := Record{
			Value:   val,
			Created: time.Now().UTC(),
		}
		c.cache = append(c.cache, r)
		return nil
	}

	// otherwise..
	return ErrRecordExists
}

// gets the length of the cache
func (c *Cache) Len() int {
	return len(c.cache)
}

// remove the individual Record from cache
func (c *Cache) expire(r []Record, i int) []Record {
	r[i] = r[len(r)-1]
	return r[:len(r)-1]
}

// bg thread supervisor
func (c *Cache) startExpiryThread() {
	go func() {
		if err := c.expireThread(); err != nil {
			log.Println("[ucache] problem in expiry thread: ", err)
			go c.startExpiryThread()
		}
	}()
}

// actual thread to check cache records at interval
func (c *Cache) expireThread() error {
	for {
		now := time.Now().UTC()
		for i, r := range c.cache {
			age := now.Sub(r.Created)
			if age > c.Expiry {
				c.cache = c.expire(c.cache, i)
			}
		}
		var err error
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 30)
	}
}
