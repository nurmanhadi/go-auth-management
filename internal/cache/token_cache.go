package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/goccy/go-json"
)

type TokenCache struct {
	cache *memcache.Client
}

func NewTokenCache(cache *memcache.Client) *TokenCache {
	return &TokenCache{
		cache: cache,
	}
}
func (c *TokenCache) SetRefreshToken(refreshToken string, value *RefreshData, exp int32) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	hashKey := sha256.Sum256([]byte(refreshToken))
	err = c.cache.Set(&memcache.Item{
		Key:        fmt.Sprintf("refresh:%s", hex.EncodeToString(hashKey[:8])),
		Value:      data,
		Expiration: exp,
	})
	if err != nil {
		return err
	}
	return nil
}
