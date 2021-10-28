package sandid

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestNew(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		sID := New()
		if bytes.Equal(sID[:], zero[:]) {
			t.Error("want false")
		}
	}

	lastTime = 1<<64 - 1

	sID := New()
	if bytes.Equal(sID[:], zero[:]) {
		t.Error("want false")
	}
}

func TestParse(t *testing.T) {
	sID, err := Parse("AAAAAAAAAAAAAAAAAAAAAA")
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if !bytes.Equal(sID[:], zero[:]) {
		t.Error("want true")
	}

	sID, err = Parse("AAECAwQFBgcICQoLDA0ODw")
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if bytes.Equal(sID[:], zero[:]) {
		t.Error("want false")
	}

	sID, err = Parse("AAAAAAAAAAAAAAAAAAAAA=")
	if err == nil {
		t.Fatal("expected error")
	} else if !bytes.Equal(sID[:], zero[:]) {
		t.Error("want true")
	}
}

func TestMustParse(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unexpected panic %q", r)
			}
		}()

		MustParse("AAECAwQFBgcICQoLDA0ODw")
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic")
			}
		}()

		MustParse("AAAAAAAAAAAAAAAAAAAAA=")
	}()
}

func TestSandIDIsZero(t *testing.T) {
	if !(SandID{}).IsZero() {
		t.Error("want true")
	}

	if New().IsZero() {
		t.Error("want false")
	}
}

func TestSandIDString(t *testing.T) {
	sIDS := SandID{}.String()
	if want := "AAAAAAAAAAAAAAAAAAAAAA"; sIDS != want {
		t.Errorf("got %q, want %q", sIDS, want)
	}
}

func TestSandIDScan(t *testing.T) {
	sIDString := "AAECAwQFBgcICQoLDA0ODw"

	sID := SandID{}
	if err := sID.Scan(sIDString); err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if sIDS := sID.String(); sIDS != sIDString {
		t.Errorf("got %q, want %q", sIDS, sIDString)
	}

	if err := sID.Scan(""); err == nil {
		t.Fatal("expected error")
	}

	if err := sID.Scan([]byte{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
		16,
	}); err == nil {
		t.Fatal("expected error")
	}

	if err := sID.Scan(0); err == nil {
		t.Fatal("expected error")
	}
}

func TestSandIDValue(t *testing.T) {
	sID := New()

	v, err := sID.Value()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	b, ok := v.([]byte)
	if !ok {
		t.Error("want true")
	} else if equal := bytes.Equal(b, sID[:]); !equal {
		t.Error("want true")
	}
}

func TestSandIDMarshalText(t *testing.T) {
	sID := New()

	b, err := sID.MarshalText()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	s := string(b)
	if want := base64.URLEncoding.EncodeToString(sID[:])[:22]; s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}

func TestSandIDUnmarshalText(t *testing.T) {
	sID := SandID{}

	if err := sID.UnmarshalText([]byte{
		65, 65, 69, 67,
		65, 119, 81, 70,
		66, 103, 99, 73,
		67, 81, 111, 76,
		68, 65, 48, 79,
		68,
	}); err == nil {
		t.Fatal("expected error")
	}

	if err := sID.UnmarshalText([]byte{
		0xff, 65, 69, 67,
		65, 119, 81, 70,
		66, 103, 99, 73,
		67, 81, 111, 76,
		68, 65, 48, 79,
		68, 119,
	}); err == nil {
		t.Fatal("expected error")
	}

	if err := sID.UnmarshalText([]byte{
		65, 65, 69, 67,
		65, 119, 81, 70,
		66, 103, 99, 73,
		67, 81, 111, 76,
		68, 65, 48, 0xff,
		68, 119,
	}); err == nil {
		t.Fatal("expected error")
	}

	if err := sID.UnmarshalText([]byte{
		65, 65, 69, 67,
		65, 119, 81, 70,
		66, 103, 99, 73,
		67, 81, 111, 76,
		68, 65, 48, 79,
		68, 119,
	}); err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	b, err := base64.URLEncoding.DecodeString("AAECAwQFBgcICQoLDA0ODw==")
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if !bytes.Equal(b[:16], sID[:]) {
		t.Error("want true")
	}
}

func TestSandIDMarshalBinary(t *testing.T) {
	sID := New()
	b, err := sID.MarshalBinary()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if !bytes.Equal(b, sID[:]) {
		t.Error("want true")
	}
}

func TestSandIDUnmarshalBinary(t *testing.T) {
	sID := SandID{}
	if err := sID.UnmarshalBinary([]byte{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
	}); err != nil {
		t.Fatalf("unexpected error %q", err)
	}

	sIDS := sID.String()
	if want := "AAECAwQFBgcICQoLDA0ODw"; sIDS != want {
		t.Errorf("got %q, want %q", sIDS, want)
	}
}

func TestSandIDMarshalJSON(t *testing.T) {
	sID := New()
	b, err := sID.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if want := "\"" + sID.String() + "\""; string(b) != want {
		t.Errorf("got %q, want %q", b, want)
	}
}

func TestSandIDUnmarshalJSON(t *testing.T) {
	sID := SandID{}
	if err := sID.UnmarshalJSON(nil); err == nil {
		t.Fatal("expected error")
	}

	sIDString := "AAECAwQFBgcICQoLDA0ODw"
	sIDJSON := []byte("\"" + sIDString + "\"")
	if err := sID.UnmarshalJSON(sIDJSON); err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if sIDS := sID.String(); sIDS != sIDString {
		t.Errorf("got %q, want %q", sIDS, sIDString)
	}
}

func TestEqual(t *testing.T) {
	if !Equal(SandID{}, SandID{}) {
		t.Error("want true")
	}

	if Equal(New(), New()) {
		t.Error("want false")
	}
}

func TestCompare(t *testing.T) {
	if result := Compare(
		MustParse("AAAAAAAAAAAAAAAAAAAAAQ"),
		MustParse("AAAAAAAAAAAAAAAAAAAAAg"),
	); result != -1 {
		t.Errorf("got %v, want -1", result)
	}

	if result := Compare(SandID{}, SandID{}); result != 0 {
		t.Errorf("got %v, want 0", result)
	}

	if result := Compare(
		MustParse("AAAAAAAAAAAAAAAAAAAAAg"),
		MustParse("AAAAAAAAAAAAAAAAAAAAAQ"),
	); result != 1 {
		t.Errorf("got %v, want 1", result)
	}
}

func TestNullSandIDScan(t *testing.T) {
	nsID := NullSandID{}

	if err := nsID.Scan(nil); err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if !nsID.SandID.IsZero() {
		t.Error("want true")
	} else if nsID.Valid {
		t.Error("want false")
	}

	if err := nsID.Scan([]byte{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
	}); err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if nsID.SandID.IsZero() {
		t.Error("want false")
	} else if !nsID.Valid {
		t.Error("want true")
	}

	nsIDS := nsID.SandID.String()
	if want := "AAECAwQFBgcICQoLDA0ODw"; nsIDS != want {
		t.Errorf("got %q, want %q", nsIDS, want)
	}
}

func TestNullSandIDValue(t *testing.T) {
	nsID := NullSandID{}

	v, err := nsID.Value()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if v != nil {
		t.Errorf("got %v, want nil", v)
	}

	nsID.SandID = MustParse("AAECAwQFBgcICQoLDA0ODw")
	nsID.Valid = true

	v, err = nsID.Value()
	if err != nil {
		t.Fatalf("unexpected error %q", err)
	} else if v == nil {
		t.Fatal("unexpected nil")
	}

	b, ok := v.([]byte)
	if !ok {
		t.Error("want true")
	} else if !bytes.Equal(b, nsID.SandID[:]) {
		t.Error("want true")
	}
}
