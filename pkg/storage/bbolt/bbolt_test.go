//nolint

package bbolt_test

import (
	"testing"

	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/bbolt"
	"github.com/vkuksa/shortly/pkg/storage/test"
)

// Ensure service implements interface.
var _ storage.Storage[int] = (*bbolt.Storage[int])(nil)

func TestSetGetDelete(t *testing.T) {
	s := test.MustCreateStorage[int](t, "bbolt")
	defer test.MustCleanupStorage(t, s)

	test.SetGetDelete(t, s)
}

func TestClose(t *testing.T) {
	s := test.MustCreateStorage[int](t, "bbolt")
	defer test.MustCleanupStorage(t, s)

	test.Close(t, s)
}
