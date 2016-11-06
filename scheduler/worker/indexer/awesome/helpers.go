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

func generateMapOfAwesomePackages(pkgs []packageTuple) map[string]packageTuple {
	m := make(map[string]packageTuple)
	for _, pkg := range pkgs {
		key := pkg.author + "/" + pkg.repo
		m[key] = pkg
	}

	return m
}

func generateRandomAwesomePackages(numPackages int) []packageTuple {
	var packageTuples []packageTuple
	for i := 0; i < numPackages; i++ {
		packageTuples = append(packageTuples,
			packageTuple{
				author: randStringRunes(26),
				repo:   randStringRunes(26),
			},
		)
	}
	return packageTuples
}
