package main

import (
	"bytes"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
)

var (
	packageLocks                 = NewConcurrentMap()
	packageLocksModificationLock sync.Mutex
)

type PackageLock struct {
	BumpLooping      bool
	PendingBumps     int64
	BumpLoopingLock  sync.RWMutex
	PendingBumpsLock sync.Mutex
}

func NewPackageLock() *PackageLock {
	return &PackageLock{}
}

func lockMapKeyOf(author, repo string) string {
	var buffer bytes.Buffer
	buffer.WriteString(author)
	buffer.WriteByte('/')
	buffer.WriteString(repo)
	return buffer.String()
}

func recordPackageInstall(
	session *gocql.Session,
	author string,
	repo string,
) error {
	// Get the key for the author & repo.
	key := lockMapKeyOf(author, repo)

	// Check whether there is a package lock already.
	packageLock, exists := packageLocks.Get(key)
	if !exists {
		// Create the lock since one doesn't exist already.
		packageLocksModificationLock.Lock()
		// Check whether the lock still doesn't exist.
		packageLock, exists = packageLocks.Get(key)
		// If it still doesn't exist, create it doe.
		if !exists {
			packageLock = NewPackageLock()
			packageLocks.Set(key, packageLock)
		}
		packageLocksModificationLock.Unlock()
	}

	// Increment the number of pending bumps.
	packageLock.PendingBumpsLock.Lock()
	packageLock.PendingBumps = packageLock.PendingBumps + 1
	packageLock.PendingBumpsLock.Unlock()

	// Decide whether to enter the bump loop.
	if !packageLock.BumpLooping {
		// Lock and then check again.
		packageLock.BumpLoopingLock.Lock()
		// If someone else started bump looping, then exit here.
		if packageLock.BumpLooping {
			return nil
		}
		// Since nobody else is bump looping, we're going to.
		packageLock.BumpLooping = true
		packageLock.BumpLoopingLock.Unlock()

		for {
			// Now that we have a package lock, figure out how much to bump.
			packageLock.PendingBumpsLock.Lock()
			// The bump amount is however many pending installs there are plus the one
			// we were recording in the first place.
			bumpAmount := packageLock.PendingBumps
			if bumpAmount > 0 {
				packageLock.PendingBumps = 0
			}
			// Release the lock after we have what we need.
			packageLock.PendingBumpsLock.Unlock()

			// If there was something to bump, then carry on. Otherwise, remove this
			// lock and be done with it.
			if bumpAmount > 0 {
				// Perform the bump itself.
				err := common.BumpRangedInstallTotals(
					session,
					time.Now(),
					author,
					repo,
					bumpAmount,
				)

				// We couldn't perform the bump, so put it back and exit.
				if err != nil {
					packageLock.PendingBumpsLock.Lock()
					packageLock.PendingBumps = packageLock.PendingBumps + bumpAmount
					packageLock.PendingBumpsLock.Unlock()
					return err
				} // If there was no error, continue bump looping.
			} else {
				// There's nothing left to bump, so its time to exit.
				return nil
			}
		}
	} else {
		// We're no looping, so we can exit.
		return nil
	}
}
