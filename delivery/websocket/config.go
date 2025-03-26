package websocket

import (
	"net/url"
	"sync"
)

type Limitation struct {
	MaxMessageLength    int32 // todo?.
	MaxSubscriptions    int32
	MaxSubidLength      int32
	MinPowDifficulty    int32
	AuthRequired        bool
	PaymentRequired     bool
	RestrictedWrites    bool
	MaxEventTags        int32
	MaxContentLength    int32
	CreatedAtLowerLimit int64
	CreatedAtUpperLimit int64
}

type Config struct {
	mu sync.RWMutex

	Bind       string `yaml:"bind"`
	Port       uint16 `yaml:"port"`
	url        *url.URL
	limitation *Limitation
}

// GetURL safely retrieves the URL.
func (c *Config) GetURL() *url.URL {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.url
}

// SetURL safely sets the URL.
func (c *Config) SetURL(u *url.URL) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.url = u
}

// GetLimitation safely retrieves the Limitation.
func (c *Config) GetLimitation() *Limitation {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.limitation
}

// SetLimitation safely sets the Limitation.
func (c *Config) SetLimitation(l *Limitation) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.limitation = l
}
