package sandid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		sID := New()
		assert.NotZero(t, sID)
	}

	lastTime = 1<<64 - 1
	assert.NotZero(t, New())
}

func TestParse(t *testing.T) {
	sID, err := Parse("AAAAAAAAAAAAAAAAAAAAAA")
	assert.Zero(t, sID)
	assert.NoError(t, err)
	sID, err = Parse("AQIDBAUGBwgJCgsMDQ4PEA")
	assert.NotZero(t, sID)
	assert.NoError(t, err)
	sID, err = Parse("AAAAAAAAAAAAAAAAAAAAA=")
	assert.Zero(t, sID)
	assert.Error(t, err)
}

func TestMustParse(t *testing.T) {
	assert.NotPanics(t, func() {
		MustParse("AQIDBAUGBwgJCgsMDQ4PEA")
	})
	assert.Panics(t, func() {
		MustParse("AAAAAAAAAAAAAAAAAAAAA=")
	})
}

func TestSandIDIsZero(t *testing.T) {
	assert.True(t, SandID{}.IsZero())
	assert.False(t, New().IsZero())
}

func TestSandIDString(t *testing.T) {
	assert.Equal(t, "AAAAAAAAAAAAAAAAAAAAAA", SandID{}.String())
}

func TestSandIDScan(t *testing.T) {
	sID := SandID{}
	assert.NoError(t, sID.Scan("AQIDBAUGBwgJCgsMDQ4PEA"))
	assert.Equal(t, "AQIDBAUGBwgJCgsMDQ4PEA", sID.String())
	assert.Error(t, sID.Scan(""))
	assert.Error(t, sID.Scan([]byte{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
		17,
	}))
	assert.Error(t, sID.Scan(0))
}

func TestSandIDValue(t *testing.T) {
	sID := New()
	v, err := sID.Value()
	assert.NoError(t, err)
	b, ok := v.([]byte)
	assert.True(t, ok)
	assert.Equal(t, sID[:], b)
}

func TestSandIDMarshalText(t *testing.T) {
	sID := New()
	b, err := sID.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, sID.String(), string(b))
}

func TestSandIDUnmarshalText(t *testing.T) {
	sID := SandID{}
	assert.NoError(t, sID.UnmarshalText([]byte{
		65, 81, 73, 68,
		66, 65, 85, 71,
		66, 119, 103, 74,
		67, 103, 115, 77,
		68, 81, 52, 80,
		69, 65,
	}))
	assert.Equal(t, "AQIDBAUGBwgJCgsMDQ4PEA", sID.String())
}

func TestSandIDMarshalBinary(t *testing.T) {
	sID := New()
	b, err := sID.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, sID[:], b)
}

func TestSandIDUnmarshalBinary(t *testing.T) {
	sID := SandID{}
	assert.NoError(t, sID.UnmarshalBinary([]byte{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}))
	assert.Equal(t, "AQIDBAUGBwgJCgsMDQ4PEA", sID.String())
}

func TestSandIDMarshalJSON(t *testing.T) {
	sID := New()
	b, err := sID.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\""+sID.String()+"\"", string(b))
}

func TestSandIDUnmarshalJSON(t *testing.T) {
	sID := SandID{}
	assert.Error(t, sID.UnmarshalJSON(nil))
	assert.NoError(t, sID.UnmarshalJSON([]byte(`"AQIDBAUGBwgJCgsMDQ4PEA"`)))
	assert.Equal(t, "AQIDBAUGBwgJCgsMDQ4PEA", sID.String())
}

func TestEqual(t *testing.T) {
	assert.True(t, Equal(SandID{}, SandID{}))
	assert.False(t, Equal(New(), New()))
}

func TestCompare(t *testing.T) {
	assert.Equal(t, -1, Compare(
		MustParse("AAAAAAAAAAAAAAAAAAAAAQ"),
		MustParse("AAAAAAAAAAAAAAAAAAAAAg"),
	))
	assert.Equal(t, 0, Compare(SandID{}, SandID{}))
	assert.Equal(t, 1, Compare(
		MustParse("AAAAAAAAAAAAAAAAAAAAAg"),
		MustParse("AAAAAAAAAAAAAAAAAAAAAQ"),
	))
}

func TestNullSandIDScan(t *testing.T) {
	nsID := NullSandID{}
	assert.NoError(t, nsID.Scan(nil))
	assert.True(t, nsID.SandID.IsZero())
	assert.False(t, nsID.Valid)
	assert.NoError(t, nsID.Scan([]byte{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}))
	assert.False(t, nsID.SandID.IsZero())
	assert.True(t, nsID.Valid)
	assert.Equal(t, "AQIDBAUGBwgJCgsMDQ4PEA", nsID.SandID.String())
}

func TestNullSandIDValue(t *testing.T) {
	nsID := NullSandID{}
	v, err := nsID.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
	nsID.SandID = MustParse("AQIDBAUGBwgJCgsMDQ4PEA")
	nsID.Valid = true
	v, err = nsID.Value()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	b, ok := v.([]byte)
	assert.True(t, ok)
	assert.Equal(t, nsID.SandID[:], b)
}
