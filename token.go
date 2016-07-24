package redis

import (
	"encoding/json"
	"strconv"

	"gopkg.in/redis.v4"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// DefaultIncrKey TokenStore incr id
const DefaultIncrKey = "oauth2_incr"

// NewTokenStore Create a token TokenStore instance based on redis
func NewTokenStore(cfg *Config) (ts oauth2.TokenStore, err error) {
	opt := &redis.Options{
		Network:      cfg.Network,
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  cfg.PoolTimeout,
	}
	cli := redis.NewClient(opt)
	if verr := cli.Ping().Err(); verr != nil {
		err = verr
		return
	}
	ts = &TokenStore{cli: cli}
	return
}

// TokenStore Redis Token Store
type TokenStore struct {
	cli *redis.Client
}

func (rs *TokenStore) getBasicID(id int64, info oauth2.TokenInfo) string {
	return "oauth2_" + info.GetClientID() + "_" + strconv.FormatInt(id, 10)
}

// Create Create and TokenStore the new token information
func (rs *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}
	id, err := rs.cli.Incr(DefaultIncrKey).Result()
	if err != nil {
		return
	}
	pipe := rs.cli.Pipeline()
	basicID := rs.getBasicID(id, info)
	aexp := info.GetAccessExpiresIn()
	rexp := aexp

	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshExpiresIn()
		ttl := rs.cli.TTL(refresh)
		if verr := ttl.Err(); verr != nil {
			err = verr
			return
		}
		if v := ttl.Val(); v.Seconds() > 0 {
			rexp = v
		}
		if aexp.Seconds() > rexp.Seconds() {
			aexp = rexp
		}
		pipe.Set(refresh, basicID, rexp)
	}
	pipe.Set(info.GetAccess(), basicID, aexp)
	pipe.Set(basicID, jv, rexp)

	if _, verr := pipe.Exec(); verr != nil {
		err = verr
	}
	return
}

// remove
func (rs *TokenStore) remove(key string) (err error) {
	_, verr := rs.cli.Del(key).Result()
	if verr != redis.Nil {
		err = verr
	}
	return
}

// RemoveByAccess Use the access token to delete the token information(Along with the refresh token)
func (rs *TokenStore) RemoveByAccess(access string) (err error) {
	err = rs.remove(access)
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (rs *TokenStore) RemoveByRefresh(refresh string) (err error) {
	err = rs.remove(refresh)
	return
}

// get
func (rs *TokenStore) get(token string) (ti oauth2.TokenInfo, err error) {
	tv, verr := rs.cli.Get(token).Result()
	if verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	result := rs.cli.Get(tv)
	if verr := result.Err(); verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	iv, err := result.Bytes()
	if err != nil {
		return
	}
	var tm models.Token
	if verr := json.Unmarshal(iv, &tm); verr != nil {
		err = verr
		return
	}
	ti = &tm
	return
}

// GetByAccess Use the access token for token information data
func (rs *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(access)
	return
}

// GetByRefresh Use the refresh token for token information data
func (rs *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.get(refresh)
	return
}
