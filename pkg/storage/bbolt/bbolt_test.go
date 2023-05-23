//nolint

package bbolt_test

import (
	"testing"

	"github.com/vkuksa/shortly/pkg/storage/bbolt"
	"github.com/vkuksa/shortly/pkg/storage/test"
)

var (
	options = bbolt.Options{File: "test.db", Bucket: "test"}
)

func TestClose(t *testing.T) {
	s, _ := test.MustCreateStorage[int](t, "bbolt", options)

	test.Close(t, s)
}

func TestSetGetDelete(t *testing.T) {
	s, closer := test.MustCreateStorage[int](t, "bbolt", options)
	defer closer()

	test.SetGetDelete(t, s)
}
