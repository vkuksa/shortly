package redis_test

import (
	"os"
	"testing"

	"github.com/vkuksa/shortly/pkg/storage/redis"
	"github.com/vkuksa/shortly/pkg/storage/test"
)

var (
	RedisHost = "localhost"
	RedisPort = "6379"
	options   = redis.Options{
		Address:  RedisHost,
		Password: RedisPort,
		DB:       15,
	}
)

func init() {
	host, exists := os.LookupEnv("REDIS_HOST")
	if exists {
		RedisHost = host
	}
	port, exists := os.LookupEnv("REDIS_PORT")
	if exists {
		RedisPort = port
	}
}

func TestClose(t *testing.T) {
	s, _ := test.MustCreateStorage[int](t, "redis", options)

	test.Close(t, s)
}

func TestSetGetDelete(t *testing.T) {
	s, closer := test.MustCreateStorage[int](t, "redis", options)
	defer closer()

	test.SetGetDelete(t, s)
}
