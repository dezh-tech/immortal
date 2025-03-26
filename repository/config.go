package repository

import "sync"

type Config struct {
	mu sync.RWMutex

	defaultQueryLimit uint32
	maxQueryLimit     uint32
}

// GetDefaultQueryLimit safely retrieves the DefaultQueryLimit.
func (c *Config) GetDefaultQueryLimit() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.defaultQueryLimit
}

// SetDefaultQueryLimit safely sets the DefaultQueryLimit.
func (c *Config) SetDefaultQueryLimit(limit uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.defaultQueryLimit = limit
}

// GetMaxQueryLimit safely retrieves the MaxQueryLimit.
func (c *Config) GetMaxQueryLimit() uint32 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.maxQueryLimit
}

// SetMaxQueryLimit safely sets the MaxQueryLimit.
func (c *Config) SetMaxQueryLimit(limit uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxQueryLimit = limit
}
