package redis

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/redis.v5"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// NewTokenStore Create a token store instance based on redis
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

// TokenStore redis token store
type TokenStore struct {
	cli *redis.Client
}

// Create Create and store the new token information
func (rs *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}
	pipe := rs.cli.Pipeline()
	if code := info.GetCode(); code != "" {
		pipe.Set(code, jv, info.GetCodeExpiresIn())
	} else {
		basicID := uuid.NewV4().String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp

		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			pipe.Set(refresh, basicID, rexp)
		}
		pipe.Set(info.GetAccess(), basicID, aexp)
		pipe.Set(basicID, jv, rexp)
	}

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

// RemoveByCode Use the authorization code to delete the token information
func (rs *TokenStore) RemoveByCode(code string) (err error) {
	err = rs.remove(code)
	return
}

// RemoveByAccess Use the access token to delete the token information
func (rs *TokenStore) RemoveByAccess(access string) (err error) {
	err = rs.remove(access)
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (rs *TokenStore) RemoveByRefresh(refresh string) (err error) {
	err = rs.remove(refresh)
	return
}

func (rs *TokenStore) getData(key string) (ti oauth2.TokenInfo, err error) {
	result := rs.cli.Get(key)
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

func (rs *TokenStore) getBasicID(token string) (basicID string, err error) {
	tv, verr := rs.cli.Get(token).Result()
	if verr != nil {
		if verr == redis.Nil {
			return
		}
		err = verr
		return
	}
	basicID = tv
	return
}

// GetByCode Use the authorization code for token information data
func (rs *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = rs.getData(code)
	return
}

// GetByAccess Use the access token for token information data
func (rs *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := rs.getBasicID(access)
	if err != nil || basicID == "" {
		return
	}
	ti, err = rs.getData(basicID)
	return
}

// GetByRefresh Use the refresh token for token information data
func (rs *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := rs.getBasicID(refresh)
	if err != nil || basicID == "" {
		return
	}
	ti, err = rs.getData(basicID)
	return
}
