package main

import (
	"errors"

	"github.com/gofrs/flock"
)

var lock *flock.Flock

func AquireLock() error {
	lock = flock.New("status.lock")

	locked, err := lock.TryLock()
	if err != nil {
		return err
	}

	if !locked {
		return errors.New("unable to lock lockfile")
	}

	return nil
}

func ReleaseLock() {
	lock.Unlock()
}
