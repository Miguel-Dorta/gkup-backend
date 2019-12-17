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

func NewHasher(s *settings.Settings) (*Hasher, error) {
	hasher := &Hasher{
		buf: make([]byte, s.BufferSize),
	}

	switch s.HashAlgorithm {
	case "md5":
		hasher.h = md5.New()
	case "sha1":
		hasher.h = sha1.New()
	case "sha256":
		hasher.h = sha256.New()
	case "sha512":
		hasher.h = sha512.New()
	case "sha3-256":
		hasher.h = sha3.New256()
	case "sha3-512":
		hasher.h = sha3.New512()
	default:
		return nil, fmt.Errorf("invalid hash algorithm: %s", s.HashAlgorithm)
	}

	return hasher, nil
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
