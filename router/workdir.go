package main

import (
	"math/rand"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// generateWorkDirName generates a unique name for a directory in the
// construction zone.
func generateWorkDirName() string {
	unixTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	randomString := randStringRunes(16)
	return unixTimestamp + randomString
}

// randStringRunes generates a random string n runes long.
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
