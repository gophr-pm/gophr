package awesome

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randStringRunes generates a random string n runes long.
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateRandomAwesomePackages(numPackages int) []PackageTuple {
	var PackageTuples []PackageTuple
	for i := 1; i < numPackages; i++ {
		PackageTuples = append(PackageTuples,
			PackageTuple{
				author: randStringRunes(i),
				repo:   randStringRunes(i),
			},
		)
	}
	return PackageTuples
}
