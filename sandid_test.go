package sandid

import (
	"bytes"
	"encoding/base64"
	"reflect"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	for i := 0; i < 1_000_000; i++ {
		sID := New()
		if bytes.Equal(sID[:], zeroSandID[:]) {
			t.Error("want false")
		}
	}

	lastTime = 1<<64 - 1

	sID := New()
	if bytes.Equal(sID[:], zeroSandID[:]) {
		t.Error("want false")
	}
}

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		name    string
		s       string
		want    SandID
		wantErr bool
	}{
		{"Zero", "AAAAAAAAAAAAAAAAAAAAAA", zeroSandID, false},
		{"NonZero", "AAECAwQFBgcICQoLDA0ODw", SandID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, false},
		{"Invalid", "AAAAAAAAAAAAAAAAAAAAA=", zeroSandID, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sID, err := Parse(tt.s)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				if !bytes.Equal(sID[:], tt.want[:]) {
					t.Errorf("got %v, want %v", sID, tt.want)
				}
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	for _, tt := range []struct {
		name      string
		s         string
		wantPanic bool
	}{
		{"Valid", "AAECAwQFBgcICQoLDA0ODw", false},
		{"Invalid", "AAAAAAAAAAAAAAAAAAAAA=", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic {
					if r == nil {
						t.Fatal("expected panic")
					}
				} else {
					if r != nil {
						t.Fatalf("unexpected panic %v", r)
					}
				}
			}()
			MustParse(tt.s)
		})
	}
}

func TestSandIDIsZero(t *testing.T) {
	for _, tt := range []struct {
		name string
		sID  SandID
		want bool
	}{
		{"Zero", zeroSandID, true},
		{"Empty", SandID{}, true},
		{"NonZero", New(), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sID.IsZero()
			if got != tt.want {
				t.Errorf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func TestSandIDString(t *testing.T) {
	got := zeroSandID.String()
	want := "AAAAAAAAAAAAAAAAAAAAAA"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSandIDScan(t *testing.T) {
	for _, tt := range []struct {
		name    string
		value   interface{}
		want    string
		wantErr bool
	}{
		{"ValidString", "AAECAwQFBgcICQoLDA0ODw", "AAECAwQFBgcICQoLDA0ODw", false},
		{"EmptyString", "", "", true},
		{"InvalidByteSliceLength", make([]byte, 17), "", true},
		{"InvalidType", 0, "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var sID SandID
			err := sID.Scan(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				got := sID.String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestSandIDValue(t *testing.T) {
	sID := New()

	v, err := sID.Value()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	got, ok := v.([]byte)
	if !ok {
		t.Error("want true")
	}
	want := sID[:]
	if !bytes.Equal(got, want) {
		t.Error("want true")
	}
}

func TestSandIDMarshalText(t *testing.T) {
	sID := New()

	b, err := sID.MarshalText()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	got := string(b)
	want := base64.URLEncoding.EncodeToString(sID[:])[:22]
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSandIDUnmarshalText(t *testing.T) {
	for _, tt := range []struct {
		name    string
		text    []byte
		want    []byte
		wantErr bool
	}{
		{
			name: "Valid",
			text: []byte{
				65, 65, 69, 67,
				65, 119, 81, 70,
				66, 103, 99, 73,
				67, 81, 111, 76,
				68, 65, 48, 79,
				68, 119,
			},
			want: []byte{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			},
			wantErr: false,
		},
		{
			name: "InvalidLength",
			text: []byte{
				65, 65, 69, 67,
				65, 119, 81, 70,
				66, 103, 99, 73,
				67, 81, 111, 76,
				68, 65, 48, 79,
				68,
			},
			wantErr: true,
		},
		{
			name: "InvalidCharAtBeginning",
			text: []byte{
				0xff, 65, 69, 67,
				65, 119, 81, 70,
				66, 103, 99, 73,
				67, 81, 111, 76,
				68, 65, 48, 79,
				68, 119,
			},
			wantErr: true,
		},
		{
			name: "InvalidCharInMiddle",
			text: []byte{
				65, 65, 69, 67,
				65, 119, 81, 70,
				66, 103, 99, 73,
				67, 81, 111, 76,
				68, 65, 48, 0xff,
				68, 119,
			},
			wantErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var sID SandID
			err := sID.UnmarshalText(tt.text)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				got := sID[:]
				if !bytes.Equal(got, tt.want) {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSandIDMarshalBinary(t *testing.T) {
	sID := New()
	got, err := sID.MarshalBinary()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	want := sID[:]
	if !bytes.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSandIDUnmarshalBinary(t *testing.T) {
	var sID SandID
	if err := sID.UnmarshalBinary([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	got := sID.String()
	want := "AAECAwQFBgcICQoLDA0ODw"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSandIDMarshalJSON(t *testing.T) {
	sID := New()
	b, err := sID.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	got := string(b)
	want := strconv.Quote(sID.String())
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSandIDUnmarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		name    string
		json    []byte
		want    string
		wantErr bool
	}{
		{"Valid", []byte(strconv.Quote("AAECAwQFBgcICQoLDA0ODw")), "AAECAwQFBgcICQoLDA0ODw", false},
		{"Invalid", []byte("{"), "", true},
		{"Nil", nil, "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var sID SandID
			err := sID.UnmarshalJSON(tt.json)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				got := sID.String()
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestEqual(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b SandID
		want bool
	}{
		{"ZeroAndEmpty", zeroSandID, SandID{}, true},
		{"SameAndCopy", MustParse("AAECAwQFBgcICQoLDA0ODw"), MustParse("AAECAwQFBgcICQoLDA0ODw"), true},
		{"Different", New(), New(), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Equal(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b SandID
		want int
	}{
		{"LessThan", MustParse("AAAAAAAAAAAAAAAAAAAAAQ"), MustParse("AAAAAAAAAAAAAAAAAAAAAg"), -1},
		{"EqualTo", zeroSandID, SandID{}, 0},
		{"GreaterThan", MustParse("AAAAAAAAAAAAAAAAAAAAAg"), MustParse("AAAAAAAAAAAAAAAAAAAAAQ"), 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNullSandIDScan(t *testing.T) {
	for _, tt := range []struct {
		name      string
		input     interface{}
		wantID    SandID
		wantValid bool
		wantErr   bool
	}{
		{"Valid", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, SandID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, true, false},
		{"Invalid", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, zeroSandID, false, true},
		{"Nil", nil, zeroSandID, false, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var nsID NullSandID
			err := nsID.Scan(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}
				if got, want := nsID.SandID, tt.wantID; !bytes.Equal(got[:], want[:]) {
					t.Errorf("got %v, want %v", got, want)
				}
				if got, want := nsID.Valid, tt.wantValid; got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestNullSandIDValue(t *testing.T) {
	for _, tt := range []struct {
		name string
		nsID NullSandID
		want interface{}
	}{
		{"Zero", NullSandID{}, nil},
		{"NonZero", NullSandID{SandID: MustParse("AAECAwQFBgcICQoLDA0ODw"), Valid: true}, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nsID.Value()
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
