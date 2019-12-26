package checkHash

import "github.com/Miguel-Dorta/gkup-backend/pkg/hash"

func ValidAlgorithm(algorithm string) bool {
	_, ok := hash.Algorithms[algorithm]
	return ok
}
