package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/utils/uuid"
)

var (
	_             oauth2.TokenStore = &TokenStore{}
	jsonMarshal                     = jsoniter.Marshal
	jsonUnmarshal                   = jsoniter.Unmarshal
)

// NewRedisStore create an instance of a redis store
func NewRedisStore(opts *Options) *TokenStore {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisStoreWithCli(redis.NewClient(opts.redisOptions()))
}

// NewRedisStoreWithCli create an instance of a redis store
func NewRedisStoreWithCli(cli *redis.Client) *TokenStore {
	return &TokenStore{
		cli: cli,
	}
}

// NewRedisClusterStore create an instance of a redis cluster store
func NewRedisClusterStore(opts *ClusterOptions) *TokenStore {
	if opts == nil {
		panic("options cannot be nil")
	}
	return NewRedisClusterStoreWithCli(redis.NewClusterClient(opts.redisClusterOptions()))
}

// NewRedisClusterStoreWithCli create an instance of a redis cluster store
func NewRedisClusterStoreWithCli(cli *redis.ClusterClient) *TokenStore {
	return &TokenStore{
		cli: cli,
	}
}

type clienter interface {
	Get(key string) *redis.StringCmd
	TxPipeline() redis.Pipeliner
	Del(keys ...string) *redis.IntCmd
	Close() error
}

// TokenStore redis token store
type TokenStore struct {
	cli clienter
}

// Close close the store
func (s *TokenStore) Close() error {
	return s.cli.Close()
}

// Create Create and store the new token information
func (s *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	ct := time.Now()
	jv, err := jsonMarshal(info)
	if err != nil {
		return
	}

	pipe := s.cli.TxPipeline()
	if code := info.GetCode(); code != "" {
		pipe.Set(code, jv, info.GetCodeExpiresIn())
	} else {
		basicID := uuid.Must(uuid.NewRandom()).String()
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
func (s *TokenStore) remove(key string) (err error) {
	_, verr := s.cli.Del(key).Result()
	if verr != redis.Nil {
		err = verr
	}
	return
}

// RemoveByCode Use the authorization code to delete the token information
func (s *TokenStore) RemoveByCode(code string) (err error) {
	err = s.remove(code)
	return
}

// RemoveByAccess Use the access token to delete the token information
func (s *TokenStore) RemoveByAccess(access string) (err error) {
	err = s.remove(access)
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (s *TokenStore) RemoveByRefresh(refresh string) (err error) {
	err = s.remove(refresh)
	return
}

func (s *TokenStore) getData(key string) (ti oauth2.TokenInfo, err error) {
	result := s.cli.Get(key)
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
	if verr := jsonUnmarshal(iv, &tm); verr != nil {
		err = verr
		return
	}
	ti = &tm
	return
}

func (s *TokenStore) getBasicID(token string) (basicID string, err error) {
	tv, verr := s.cli.Get(token).Result()
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
func (s *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = s.getData(code)
	return
}

// GetByAccess Use the access token for token information data
func (s *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := s.getBasicID(access)
	if err != nil || basicID == "" {
		return
	}
	ti, err = s.getData(basicID)
	return
}

// GetByRefresh Use the refresh token for token information data
func (s *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := s.getBasicID(refresh)
	if err != nil || basicID == "" {
		return
	}
	ti, err = s.getData(basicID)
	return
}
