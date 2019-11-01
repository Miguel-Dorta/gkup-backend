package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
)

var (
	Algorithms = map[string]func() hash.Hash{
		"md5":      md5.New,
		"sha1":     sha1.New,
		"sha256":   sha256.New,
		"sha512":   sha512.New,
		"sha3-256": sha3.New256,
		"sha3-512": sha3.New512,
	}
)

func HashFile(path string, h hash.Hash, buf []byte) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file (%s): %w", path, err)
	}
	defer f.Close()

	h.Reset()
	if _, err = io.CopyBuffer(h, f, buf); err != nil {
		return nil, fmt.Errorf("error hashing file (%s): %w", path, err)
	}
	return h.Sum(nil), err
}
