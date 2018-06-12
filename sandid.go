package sandid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net"
	"strings"
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
func Parse(s string) SandID {
	sID := SandID{}
	hex.Decode(sID[:], []byte(s))
	return sID
}

// SplitParse slices s into all substrings separated by the sep and parses the
// substrings between those separators into a new slice of the `SandID`.
func SplitParse(s, sep string) []SandID {
	ss := strings.Split(s, sep)
	sIDs := make([]SandID, 0, len(ss))
	for _, s := range ss {
		sIDs = append(sIDs, SandIDFromString(s))
	}
	return sIDs
}

// IsZero reports whether the sID is a zero instance of the `SandID`.
func (sID SandID) IsZero() bool {
	return sID.Equal(SandID{})
}

// Equal reports whether the sID and the a are equal.
func (sID SandID) Equal(a SandID) bool {
	return sID.Compare(a) == 0
}

// Compare returns an integer comparing the sID and the a lexicographically. The
// result will be 0 if sID == a, -1 if sID < a, and +1 if sID > a.
func (sID SandID) Compare(a SandID) int {
	return bytes.Compare(sID[:], a[:])
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
	copy(sID[:], value.([]byte))
	return nil
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
	_, err := hex.Decode(sID[:], text)
	return err
}

// MarshalBinary implements the `encoding.BinaryMarshaler`.
func (sID SandID) MarshalBinary() ([]byte, error) {
	return sID[:], nil
}

// UnmarshalBinary implements the `encoding.BinaryUnmarshaler`.
func (sID *SandID) UnmarshalBinary(data []byte) error {
	copy(sID[:], data)
	return nil
}

// MarshalJSON implements the `json.Marshaler`.
func (sID SandID) MarshalJSON() ([]byte, error) {
	return json.Marshal(sID.String())
}

// UnmarshalJSON implements the `json.Unmarshaler`.
func (sID *SandID) UnmarshalJSON(data []byte) error {
	return sID.UnmarshalText(data)
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
