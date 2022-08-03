/*
Package sandid implements a unique ID generation algorithm to ensure that every
grain of sand on Earth has its own ID.
*/
package sandid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// SandID is an ID of sand.
type SandID [16]byte

var (
	zero SandID

	storageMutex    sync.Mutex
	luckyNibble     byte
	clockSequence   uint16
	hardwareAddress [6]byte
	lastTime        uint64

	encoding = [64]byte{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_',
	}
	decoding [256]byte
)

func init() {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf(
			"sandid: failed to read random bytes: %v",
			err,
		))
	}

	luckyNibble = b[0]
	clockSequence = binary.BigEndian.Uint16(b[1:3])

	copy(hardwareAddress[:], b[3:])
	hardwareAddress[0] |= 0x01
	netInterfaces, _ := net.Interfaces()
	for _, ni := range netInterfaces {
		if len(ni.HardwareAddr) >= 6 {
			copy(hardwareAddress[:], ni.HardwareAddr)
			break
		}
	}

	for i := 0; i < len(decoding); i++ {
		decoding[i] = 0xff
	}

	for i := 0; i < len(encoding); i++ {
		decoding[encoding[i]] = byte(i)
	}
}

// New returns a new instance of the [SandID].
func New() SandID {
	storageMutex.Lock()
	defer storageMutex.Unlock()

	timeNow := 122192928000000000 + uint64(time.Now().UnixNano()/100)
	if timeNow <= lastTime {
		clockSequence++
	}

	lastTime = timeNow

	sID := SandID{}
	binary.BigEndian.PutUint16(sID[0:], uint16(timeNow>>44))
	binary.BigEndian.PutUint16(sID[2:], uint16(timeNow>>28))
	binary.BigEndian.PutUint32(sID[4:], uint32(timeNow<<4))
	binary.BigEndian.PutUint16(sID[8:], clockSequence)
	copy(sID[10:], hardwareAddress[:])

	sID[7] = sID[7]&0xf0 | luckyNibble&0x0f

	return sID
}

// Parse parses the s into a new instance of the [SandID].
func Parse(s string) (SandID, error) {
	sID := SandID{}
	return sID, sID.UnmarshalText([]byte(s))
}

// MustParse is like the [Parse], but panics if the s cannot be parsed.
func MustParse(s string) SandID {
	sID, err := Parse(s)
	if err != nil {
		panic(err)
	}

	return sID
}

// IsZero reports whether the sID is zero.
func (sID SandID) IsZero() bool {
	return Equal(sID, zero)
}

// String returns the serialization of the sID.
func (sID SandID) String() string {
	b, _ := sID.MarshalText()
	return string(b)
}

// Scan implements the [database/sql.Scanner].
//
// The value must be a string or []byte.
func (sID *SandID) Scan(value interface{}) error {
	switch value := value.(type) {
	case string:
		return sID.UnmarshalText([]byte(value))
	case []byte:
		return sID.UnmarshalBinary(value)
	}

	return errors.New("sandid: invalid type value")
}

// Value implements the [driver.Valuer].
func (sID SandID) Value() (driver.Value, error) {
	return sID.MarshalBinary()
}

// MarshalText implements the [encoding.TextMarshaler].
func (sID SandID) MarshalText() ([]byte, error) {
	d := make([]byte, 22)

	si, di := 0, 0
	for ; si < 15; si, di = si+3, di+4 { // si < (len(sID) / 3) * 3
		v := uint(sID[si])<<16 | uint(sID[si+1])<<8 | uint(sID[si+2])
		d[di] = encoding[v>>18&0x3f]
		d[di+1] = encoding[v>>12&0x3f]
		d[di+2] = encoding[v>>6&0x3f]
		d[di+3] = encoding[v&0x3f]
	}

	v := uint(sID[si]) << 16
	d[di] = encoding[v>>18&0x3f]
	d[di+1] = encoding[v>>12&0x3f]

	return d, nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler].
func (sID *SandID) UnmarshalText(text []byte) error {
	if len(text) != 22 {
		return errors.New("sandid: invalid length string")
	}

	si, n := 0, 0
	if strconv.IntSize >= 64 {
		for ; si <= 22-8 && n <= 16-8; si, n = si+8, n+6 {
			n1 := decoding[text[si]]
			n2 := decoding[text[si+1]]
			n3 := decoding[text[si+2]]
			n4 := decoding[text[si+3]]
			n5 := decoding[text[si+4]]
			n6 := decoding[text[si+5]]
			n7 := decoding[text[si+6]]
			n8 := decoding[text[si+7]]
			if n1|n2|n3|n4|n5|n6|n7|n8 == 0xff {
				return errors.New("sandid: invalid string")
			}

			binary.BigEndian.PutUint64(
				sID[n:],
				uint64(n1)<<58|
					uint64(n2)<<52|
					uint64(n3)<<46|
					uint64(n4)<<40|
					uint64(n5)<<34|
					uint64(n6)<<28|
					uint64(n7)<<22|
					uint64(n8)<<16,
			)
		}
	}

	for ; si <= 22-4 && n <= 16-4; si, n = si+4, n+3 {
		n1 := decoding[text[si]]
		n2 := decoding[text[si+1]]
		n3 := decoding[text[si+2]]
		n4 := decoding[text[si+3]]
		if n1|n2|n3|n4 == 0xff {
			return errors.New("sandid: invalid string")
		}

		binary.BigEndian.PutUint32(
			sID[n:],
			uint32(n1)<<26|
				uint32(n2)<<20|
				uint32(n3)<<14|
				uint32(n4)<<8,
		)
	}

	b := [4]byte{}
	for i := 0; i < 4 && si < 22; i, si = i+1, si+1 {
		if b[i] = decoding[text[si]]; b[i] == 0xff {
			return errors.New("sandid: invalid string")
		}
	}

	v := uint(b[0])<<18 | uint(b[1])<<12 | uint(b[2])<<6 | uint(b[3])
	sID[n] = byte(v >> 16)

	return nil
}

// MarshalBinary implements the [encoding.BinaryMarshaler].
func (sID SandID) MarshalBinary() ([]byte, error) {
	return sID[:], nil
}

// UnmarshalBinary implements the [encoding.BinaryUnmarshaler].
func (sID *SandID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return errors.New("sandid: invalid length bytes")
	}

	copy(sID[:], data)

	return nil
}

// MarshalJSON implements the [json.Marshaler].
func (sID SandID) MarshalJSON() ([]byte, error) {
	return json.Marshal(sID.String())
}

// UnmarshalJSON implements the [json.Unmarshaler].
func (sID *SandID) UnmarshalJSON(data []byte) error {
	s := ""
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return sID.UnmarshalText([]byte(s))
}

// Equal reports whether the a and b are equal.
func Equal(a, b SandID) bool {
	return Compare(a, b) == 0
}

// Compare returns an integer comparing the a and b lexicographically. The
// result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func Compare(a, b SandID) int {
	return bytes.Compare(a[:], b[:])
}

// NullSandID represents an instance of the [SandID] that may be null. It
// implements the [database/sql.Scanner] so it can be used as a scan
// destination.
type NullSandID struct {
	SandID SandID
	Valid  bool
}

// Scan implements the [database/sql.Scanner].
func (nsID *NullSandID) Scan(value interface{}) error {
	if value == nil {
		nsID.SandID, nsID.Valid = SandID{}, false
		return nil
	}

	nsID.Valid = true

	return nsID.SandID.Scan(value)
}

// Value implements the [driver.Valuer].
func (nsID NullSandID) Value() (driver.Value, error) {
	if !nsID.Valid {
		return nil, nil
	}

	return nsID.SandID.Value()
}
