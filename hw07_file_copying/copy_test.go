package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	source := "testdata/input.txt"
	destination := "/tmp/copy_test_out.txt"
	t.Run("Unsupported file", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy("/dev/zero", destination, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})
	t.Run("Source path empty", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy("", destination, 0, 0)
		require.ErrorIs(t, err, ErrSourcePathEmpty)
	})
	t.Run("Destination path empty", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy(source, "", 0, 0)
		require.ErrorIs(t, err, ErrDestinationPathEmpty)
	})
	t.Run("Limit over offset", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy(source, destination, 10_000, 10_000)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("Limit over offset and size", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy(source, destination, 10_000, 5_000)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("Limit over offset and zero size", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy(source, destination, 10_000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})
	t.Run("Limit over size", func(t *testing.T) {
		defer os.Remove(destination)
		err := Copy(source, destination, 0, 10_000)
		require.Nil(t, err)
		destinationStat, err := os.Stat(destination)
		require.Nil(t, err)
		sourceStat, err := os.Stat(source)
		require.Nil(t, err)
		require.Equal(t, destinationStat.Size(), sourceStat.Size())
	})
	t.Run("Limit copy", func(t *testing.T) {
		defer os.Remove(destination)
		limit = 100
		err := Copy(source, destination, 0, limit)
		require.Nil(t, err)
		destinationStat, err := os.Stat(destination)
		require.Nil(t, err)
		require.Equal(t, destinationStat.Size(), limit)
	})
}
