package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/redis/go-redis/v9"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	host = os.Getenv("REDIS_HOST") //"redis:6379"
	port = os.Getenv("REDIS_PORT") //"redis:6379"
	addr = fmt.Sprintf("%s:%s", host, port)
	db   = 15
)

func init() {
	// default the redis host and port
	if host == "" {
		host = "localhost"
		addr = fmt.Sprintf("%s:%s", host, port)
	}
	if port == "" {
		port = "6379"
		addr = fmt.Sprintf("%s:%s", host, port)
	}
}

func TestTokenStore(t *testing.T) {
	Convey("Test redis token store", t, func() {
		opts := &redis.UniversalOptions{
			Addrs: []string{addr},
			DB:    db,
		}
		store := NewRedisStore(opts)
		ctx := context.Background()

		Convey("Test authorization code store", func() {
			info := &models.Token{
				ClientID:      "1",
				UserID:        "1_1",
				RedirectURI:   "http://localhost/",
				Scope:         "all",
				Code:          "11_11_11",
				CodeCreateAt:  time.Now(),
				CodeExpiresIn: time.Second * 5,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			cinfo, err := store.GetByCode(ctx, info.Code)
			So(err, ShouldBeNil)
			So(cinfo.GetUserID(), ShouldEqual, info.UserID)

			err = store.RemoveByCode(ctx, info.Code)
			So(err, ShouldBeNil)

			cinfo, err = store.GetByCode(ctx, info.Code)
			So(err, ShouldBeNil)
			So(cinfo, ShouldBeNil)
		})

		Convey("Test access token store", func() {
			info := &models.Token{
				ClientID:        "1",
				UserID:          "1_1",
				RedirectURI:     "http://localhost/",
				Scope:           "all",
				Access:          "1_1_1",
				AccessCreateAt:  time.Now(),
				AccessExpiresIn: time.Second * 5,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			ainfo, err := store.GetByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)

			ainfo, err = store.GetByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo, ShouldBeNil)
		})

		Convey("Test refresh token store", func() {
			info := &models.Token{
				ClientID:         "1",
				UserID:           "1_2",
				RedirectURI:      "http://localhost/",
				Scope:            "all",
				Access:           "1_2_1",
				AccessCreateAt:   time.Now(),
				AccessExpiresIn:  time.Second * 5,
				Refresh:          "1_2_2",
				RefreshCreateAt:  time.Now(),
				RefreshExpiresIn: time.Second * 15,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			rinfo, err := store.GetByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)

			rinfo, err = store.GetByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo, ShouldBeNil)
		})
	})
}

func TestTokenStoreWithKeyNamespace(t *testing.T) {
	ctx := context.Background()
	Convey("Test redis token store", t, func() {
		opts := &redis.UniversalOptions{
			Addrs: []string{addr},
			DB:    db,
		}
		store := NewRedisStore(opts, "test:")

		Convey("Test authorization code store", func() {
			info := &models.Token{
				ClientID:      "1",
				UserID:        "1_1",
				RedirectURI:   "http://localhost/",
				Scope:         "all",
				Code:          "11_11_11",
				CodeCreateAt:  time.Now(),
				CodeExpiresIn: time.Second * 5,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			cinfo, err := store.GetByCode(ctx, info.Code)
			So(err, ShouldBeNil)
			So(cinfo.GetUserID(), ShouldEqual, info.UserID)

			err = store.RemoveByCode(ctx, info.Code)
			So(err, ShouldBeNil)

			cinfo, err = store.GetByCode(ctx, info.Code)
			So(err, ShouldBeNil)
			So(cinfo, ShouldBeNil)
		})

		Convey("Test access token store", func() {
			info := &models.Token{
				ClientID:        "1",
				UserID:          "1_1",
				RedirectURI:     "http://localhost/",
				Scope:           "all",
				Access:          "1_1_1",
				AccessCreateAt:  time.Now(),
				AccessExpiresIn: time.Second * 5,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			ainfo, err := store.GetByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)

			ainfo, err = store.GetByAccess(ctx, info.GetAccess())
			So(err, ShouldBeNil)
			So(ainfo, ShouldBeNil)
		})

		Convey("Test refresh token store", func() {
			info := &models.Token{
				ClientID:         "1",
				UserID:           "1_2",
				RedirectURI:      "http://localhost/",
				Scope:            "all",
				Access:           "1_2_1",
				AccessCreateAt:   time.Now(),
				AccessExpiresIn:  time.Second * 5,
				Refresh:          "1_2_2",
				RefreshCreateAt:  time.Now(),
				RefreshExpiresIn: time.Second * 15,
			}
			err := store.Create(ctx, info)
			So(err, ShouldBeNil)

			rinfo, err := store.GetByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo.GetUserID(), ShouldEqual, info.GetUserID())

			err = store.RemoveByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)

			rinfo, err = store.GetByRefresh(ctx, info.GetRefresh())
			So(err, ShouldBeNil)
			So(rinfo, ShouldBeNil)
		})
	})
}
