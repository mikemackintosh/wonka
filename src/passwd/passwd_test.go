package passwd

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		Have []byte
		Want Entry
	}{
		{
			Have: []byte("root:x:0:"),
			Want: Entry{
				Username: "root",
				Password: "x",
				UID:      0,
				Error: []error{
					errors.New("passwd entry has less than 7 segments"),
					errors.New("invalid gid"),
					errors.New("invalid info field"),
					errors.New("invalid homedir"),
					errors.New("invalid shell"),
				},
			},
		},
		{
			Have: []byte("nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin"),
			Want: Entry{
				Username: "nobody",
				Password: "x",
				UID:      65534,
				GID:      65534,
				Info:     "nobody",
				HomeDir:  "/nonexistent",
				Shell:    "/usr/sbin/nologin",
			},
		},
	}

	for testNum, test := range tests {
		var passwd Entries
		err := Unmarshal(test.Have, &passwd)
		if err != nil {
			t.Error(err)
		}

		if len(passwd) == 0 {
			t.Fatalf("%d) expected passwd size to be > 0", testNum)
		}

		if len(passwd[0].Error) > 0 {
			if !reflect.DeepEqual(passwd[0].Error, test.Want.Error) {
				t.Errorf("%d) expected %#v, have %#v", testNum, passwd[0].Error, test.Want.Error)

				for _, err := range passwd[0].Error {
					t.Error(err)
				}
				t.Fatal()
			}
		}

		if !reflect.DeepEqual(passwd[0], test.Want) {
			t.Errorf("%d) expected %#v, have %#v", testNum, passwd[0], test.Want)
		}
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		Have Entries
		Want []byte
	}{
		{
			Have: Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
				},
			},
			Want: []byte("root:x:0:0:::"),
		},
		{
			Have: Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
				},
				Entry{
					Username: "nobody",
					Password: "x",
					UID:      0,
					Info:     "nobody",
					HomeDir:  "/nonexistent",
					Shell:    "/usr/sbin/nonexistent",
				},
			},
			Want: []byte("root:x:0:0:::\nnobody:x:0:0:nobody:/nonexistent:/usr/sbin/nonexistent"),
		},
	}

	for testNum, test := range tests {
		output, err := Marshal(test.Have)
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(output, test.Want) != 0 {
			t.Errorf("%d) expected %v, got %v", testNum, output, test.Want)
		}
	}
}

func TestRemoveEntry(t *testing.T) {
	tests := []struct {
		Have *Entries
		Want *Entries
	}{
		{
			Have: &Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
					GID:      0,
				},
				Entry{
					Username: "removeme",
					Password: "x",
					UID:      1,
					GID:      1,
				},
			},
			Want: &Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
					GID:      0,
				},
			},
		},
		{
			Have: &Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
					GID:      0,
				},
				Entry{
					Username: "boat",
					Password: "x",
					UID:      1,
					GID:      1,
				},
			},
			Want: &Entries{
				Entry{
					Username: "root",
					Password: "x",
					UID:      0,
					GID:      0,
				},
				Entry{
					Username: "boat",
					Password: "x",
					UID:      1,
					GID:      1,
				},
			},
		},
	}

	for i, test := range tests {
		test.Have.RemoveEntry(Entry{Username: "removeme"})
		if !reflect.DeepEqual(test.Have, test.Want) {
			t.Fatalf("%d) failed to remove entry as expected", i)
		}
	}
}
