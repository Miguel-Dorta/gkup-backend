package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/Miguel-Dorta/gkup-backend/pkg/repository/settings"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
)

type Hasher struct {
	h   hash.Hash
	buf []byte
}

var Algorithms = map[string]func() hash.Hash{
	"md5": md5.New,
	"sha1": sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
	"sha3-256": sha3.New256,
	"sha3-512": sha3.New512,
}

func NewHasher(s *settings.Settings) (*Hasher, error) {
	h, ok := Algorithms[s.HashAlgorithm]
	if !ok {
		return nil, fmt.Errorf("invalid hash algorithm: %s", s.HashAlgorithm)
	}

	return &Hasher{
		h:   h(),
		buf: make([]byte, s.BufferSize),
	}, nil
}

func (h *Hasher) HashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file (%s): %w", path, err)
	}
	defer f.Close()

	h.h.Reset()
	if _, err = io.CopyBuffer(h.h, f, h.buf); err != nil {
		return nil, fmt.Errorf("error hashing file (%s): %w", path, err)
	}
	return h.h.Sum(nil), err
}
