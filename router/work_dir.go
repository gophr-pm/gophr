package main

import (
	"math/rand"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateWorkDirName() string {
	unixTimestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	randomString := randStringRunes(16)
	return unixTimestamp + randomString
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
