package utils

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const UINT24_MAX_VALUE = 16_777_215

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func PutUint24(buf []byte, v uint32) error {
	if len(buf) < 3 {
		return errors.New("buffer is smaller then 3 byte int")
	}
	if v > UINT24_MAX_VALUE {
		return errors.New("exceeded maximum value for 3 bytes int")
	}
	inner := make([]byte, 4)
	binary.BigEndian.PutUint32(inner, v)
	copy(buf, inner[1:])
	return nil
}
