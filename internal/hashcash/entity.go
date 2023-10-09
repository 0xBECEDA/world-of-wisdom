package hashcash

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"math/rand"
	"time"
	"world-of-wisdom/utils"
)

const (
	version1 = 1

	zeroByte = '0'
)

type Hashcash struct {
	Version  int
	Bits     int
	Date     time.Time
	Resource string
	Rand     []byte
	Counter  int64
}

func NewHashcash(bits int, resource string) *Hashcash {
	return &Hashcash{
		Version:  version1,
		Bits:     bits,
		Date:     time.Now(),
		Resource: resource,
		Rand:     randBytes(),
	}
}

func (h *Hashcash) String() string {
	return fmt.Sprintf("%d:%d:%s:%s:%s",
		h.Version,
		h.Bits,
		h.Resource,
		base64.StdEncoding.EncodeToString(h.Rand),
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", h.Counter))),
	)
}

func randBytes() []byte {
	return big.NewInt(int64(rand.Uint64())).Bytes()
}

func (h *Hashcash) Verify() bool {
	hashString := utils.Data2Sha1Hash(h.String())
	if h.Bits > len(hashString) {
		return false
	}

	for _, ch := range hashString[:h.Bits] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}

func (h *Hashcash) ComputeHash(maxIterations int64) error {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		if h.Verify() {
			return nil
		}
		h.Counter++
	}
	return ErrMaxIterationsExceeded
}
