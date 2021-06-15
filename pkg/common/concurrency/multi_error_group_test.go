package concurrency

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMultiErrorGroupExecuteWithoutError(t *testing.T) {
	meg := NewMultiErrorGroup()

	meg.Go(func() error {
		return nil
	})

	meg.Go(func() error {
		return nil
	})
	meg.Go(func() error {
		return nil
	})

	err := meg.Wait()
	require.NoError(t, err)
}

func TestMultiErrorGroupExecuteWithError(t *testing.T) {
	var errMessage = "error"

	meg := NewMultiErrorGroup()

	meg.Go(func() error {
		return nil
	})

	meg.Go(func() error {
		return errors.New(errMessage)
	})
	meg.Go(func() error {
		return nil
	})

	err := meg.Wait()
	require.Error(t, err)
	require.EqualError(t, err, errMessage)
}

func TestMultiErrorGroupExecuteWithMultipleErrors(t *testing.T) {
	meg := NewMultiErrorGroup()

	meg.Go(func() error {
		return errors.New("error-0")
	})

	meg.Go(func() error {
		return nil
	})
	meg.Go(func() error {
		return errors.New("error-0")
	})

	err := meg.Wait()
	require.Error(t, err)
	require.EqualError(t, err, "error-0\nerror-0")
}
