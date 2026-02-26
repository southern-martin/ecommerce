package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type UUID [16]byte

var Nil UUID

func New() UUID {
	var u UUID
	_, err := rand.Read(u[:])
	if err != nil {
		panic(fmt.Errorf("uuid: failed to read randomness: %w", err))
	}
	// UUIDv4
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return u
}

func Parse(s string) (UUID, error) {
	var u UUID
	s = strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return Nil, errors.New("uuid: invalid format")
	}
	if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
		return Nil, errors.New("uuid: invalid length")
	}
	raw := strings.Join(parts, "")
	buf, err := hex.DecodeString(raw)
	if err != nil {
		return Nil, fmt.Errorf("uuid: invalid hex: %w", err)
	}
	if len(buf) != 16 {
		return Nil, errors.New("uuid: invalid byte length")
	}
	copy(u[:], buf)
	return u, nil
}

func MustParse(s string) UUID {
	u, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func (u UUID) String() string {
	var b [36]byte
	hex.Encode(b[0:8], u[0:4])
	b[8] = '-'
	hex.Encode(b[9:13], u[4:6])
	b[13] = '-'
	hex.Encode(b[14:18], u[6:8])
	b[18] = '-'
	hex.Encode(b[19:23], u[8:10])
	b[23] = '-'
	hex.Encode(b[24:36], u[10:16])
	return string(b[:])
}
