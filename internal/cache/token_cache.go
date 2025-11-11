package cache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type TokenCache struct {
	cache *memcache.Client
}

func NewTokenCache(cache *memcache.Client) *TokenCache {
	return &TokenCache{
		cache: cache,
	}
}
func (c *TokenCache) SetRefreshToken(key string, value string, exp int32) error {
	err := c.cache.Set(&memcache.Item{
		Key:        fmt.Sprintf("refresh:%s", key),
		Value:      []byte(value),
		Expiration: exp,
	})
	if err != nil {
		return err
	}
	return nil
}
func (c *TokenCache) GetRefreshToken(key string) (string, error) {
	item, err := c.cache.Get(fmt.Sprintf("refresh:%s", key))
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}
func (c *TokenCache) RemoveRefreshToken(key string) error {
	err := c.cache.Delete(fmt.Sprintf("refresh:%s", key))
	if err != nil {
		return err
	}
	return nil
}
