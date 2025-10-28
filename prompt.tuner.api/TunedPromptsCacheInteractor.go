package main

import (
	"sync"
	"time"
)

type TunedPromptCacheItem struct {
	prompt     string
	expiration time.Time
}

type TunedPromptsCache struct {
	items map[string]TunedPromptCacheItem
	mutex sync.RWMutex
}

var promptsCache = newPromptsCache()

func newPromptsCache() *TunedPromptsCache {
	cache := &TunedPromptsCache{
		items: make(map[string]TunedPromptCacheItem),
	}
	return cache
}

func GetTunedPromptsCache() *TunedPromptsCache {
	return promptsCache
}

func (c *TunedPromptsCache) Add(gitlabProjectId string, prompt string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.items[gitlabProjectId] = TunedPromptCacheItem{prompt: prompt, expiration: time.Now().Add(1 * time.Hour)}
}

func (c *TunedPromptsCache) Exists(gitlabProjectId string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.items[gitlabProjectId]
	if !exists {
		return false
	}
	if time.Now().After(item.expiration) {
		return false
	}
	return true
}

func (c *TunedPromptsCache) Get(gitlabProjectId string) *string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, exists := c.items[gitlabProjectId]
	if !exists {
		return nil
	}
	if time.Now().After(item.expiration) {
		return nil
	}
	return &item.prompt
}
