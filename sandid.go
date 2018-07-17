package sandid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"
)

// SandID is an ID of sand.
type SandID [16]byte

var (
	storageOnce     sync.Once
	storageMutex    sync.Mutex
	clockSequence   uint16
	hardwareAddress [6]byte
	lastTime        uint64
)

// New returns a new instance of the `SandID`.
func New() SandID {
	storageOnce.Do(func() {
		b := make([]byte, 2)
		if _, err := rand.Read(b); err != nil {
			panic(err)
		}

		clockSequence = binary.BigEndian.Uint16(b)

		if is, err := net.Interfaces(); err == nil {
			for _, i := range is {
				if len(i.HardwareAddr) >= 6 {
					copy(hardwareAddress[:], i.HardwareAddr)
					return
				}
			}
		}

		if _, err := rand.Read(hardwareAddress[:]); err != nil {
			panic(err)
		}

		hardwareAddress[0] |= 0x01
	})

	storageMutex.Lock()
	defer storageMutex.Unlock()

	timeNow := 122192928000000000 + uint64(time.Now().UnixNano()/100)
	if timeNow <= lastTime {
		clockSequence++
	}
	lastTime = timeNow

	sID := SandID{}

	binary.BigEndian.PutUint16(sID[0:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(sID[2:], uint16(timeNow>>32))
	binary.BigEndian.PutUint32(sID[4:], uint32(timeNow))
	binary.BigEndian.PutUint16(sID[8:], clockSequence)

	copy(sID[10:], hardwareAddress[:])

	return sID
}

// Parse parses the s into a new instance of the `SandID`.
func Parse(s string) (SandID, error) {
	sID := SandID{}
	return sID, sID.UnmarshalText([]byte(s))
}

// MustParse is like the `Parse()`, but panics if the s cannot be parsed.
func MustParse(s string) SandID {
	sID, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return sID
}

// IsZero reports whether the sID is a zero instance of the `SandID`.
func (sID SandID) IsZero() bool {
	return Equal(sID, SandID{})
}

// String returns the serialization of the sID.
func (sID SandID) String() string {
	b := make([]byte, 32)
	hex.Encode(b, sID[:])
	return string(b)
}

// Scan implements the `sql.Scanner`.
//
// value must be a `[]byte`.
func (sID *SandID) Scan(value interface{}) error {
	switch value := value.(type) {
	case string:
		return sID.UnmarshalText([]byte(value))
	case []byte:
		return sID.UnmarshalBinary(value)
	}
	return errors.New("sandid: invalid type SandID value")
}

// Value implements the `driver.Valuer`.
func (sID SandID) Value() (driver.Value, error) {
	return sID[:], nil
}

// MarshalText implements the `encoding.TextMarshaler`.
func (sID SandID) MarshalText() ([]byte, error) {
	return []byte(sID.String()), nil
}

// UnmarshalText implements the `encoding.TextUnmarshaler`.
func (sID *SandID) UnmarshalText(text []byte) error {
	if len(text) != 32 {
		return errors.New("sandid: invalid length SandID string")
	}
	_, err := hex.Decode(sID[:], text)
	return err
}

// MarshalBinary implements the `encoding.BinaryMarshaler`.
func (sID SandID) MarshalBinary() ([]byte, error) {
	return sID[:], nil
}

// UnmarshalBinary implements the `encoding.BinaryUnmarshaler`.
func (sID *SandID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return errors.New("sandid: invalid length SandID bytes")
	}
	copy(sID[:], data)
	return nil
}

// MarshalJSON implements the `json.Marshaler`.
func (sID SandID) MarshalJSON() ([]byte, error) {
	return json.Marshal(sID.String())
}

// UnmarshalJSON implements the `json.Unmarshaler`.
func (sID *SandID) UnmarshalJSON(data []byte) error {
	s := ""
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return sID.UnmarshalText([]byte(s))
}

// Equal reports whether the a and the b are equal.
func Equal(a, b SandID) bool {
	return Compare(a, b) == 0
}

// Compare returns an integer comparing the a and the b lexicographically. The
// result will be 0 if a == b, -1 if a < b, and +1 if a > b.
func Compare(a, b SandID) int {
	return bytes.Compare(a[:], b[:])
}

// NullSandID represents an instance of the `SandID` that may be null.
// NullSandID implements the `sql.Scanner` so it can be used as a scan
// destination.
type NullSandID struct {
	SandID SandID
	Valid  bool
}

// Scan implements the `sql.Scanner`.
func (nsID *NullSandID) Scan(value interface{}) error {
	if value == nil {
		nsID.SandID, nsID.Valid = SandID{}, false
		return nil
	}
	nsID.Valid = true
	return nsID.SandID.Scan(value)
}

// Value implements the `driver.Valuer`.
func (nsID NullSandID) Value() (driver.Value, error) {
	if !nsID.Valid {
		return nil, nil
	}
	return nsID.SandID.Value()
}
