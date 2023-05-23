package redis_test

import (
	"testing"

	"github.com/vkuksa/shortly/pkg/storage/redis"
	"github.com/vkuksa/shortly/pkg/storage/test"
)

var (
	options = redis.Options{
		Address:  "localhost:6379",
		Password: "",
		DB:       15,
	}
)

func TestClose(t *testing.T) {
	s, _ := test.MustCreateStorage[int](t, "redis", options)

	test.Close(t, s)
}

func TestSetGetDelete(t *testing.T) {
	s, closer := test.MustCreateStorage[int](t, "redis", options)
	defer closer()

	test.SetGetDelete(t, s)
}
