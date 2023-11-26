package pow

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"math/rand"
	"time"
	"world-of-wisdom/internal/utils"
)

const (
	version1 = 1
	zero     = '0'
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

func (h *Hashcash) Check() bool {
	hashString := utils.Data2Sha1Hash(h.String())
	if h.Bits > len(hashString) {
		return false
	}

	for _, ch := range hashString[:h.Bits] {
		if ch != zero {
			return false
		}
	}
	return true
}

func (h *Hashcash) Compute(maxIterations int64) error {
	for h.Counter <= maxIterations || maxIterations <= 0 {
		if h.Check() {
			return nil
		}
		h.Counter++
	}
	return ErrMaxIterExceeded
}

func randBytes() []byte {
	return big.NewInt(int64(rand.Uint64())).Bytes()
}
