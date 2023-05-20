package inmem_test

import (
	"testing"

	"github.com/vkuksa/shortly/pkg/storage"
	"github.com/vkuksa/shortly/pkg/storage/inmem"
	"github.com/vkuksa/shortly/pkg/storage/test"
)

// Ensure service implements interface.
var _ storage.Storage[int] = (*inmem.Storage[int])(nil)

func TestSetGetDelete(t *testing.T) {
	s := test.MustCreateStorage[int](t, "inmem")
	defer test.MustCleanupStorage(t, s)

	test.SetGetDelete(t, s)
}

func TestClose(t *testing.T) {
	s := test.MustCreateStorage[int](t, "inmem")
	defer test.MustCleanupStorage(t, s)

	test.Close(t, s)
}
