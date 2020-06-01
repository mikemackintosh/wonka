package passwd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mikemackintosh/wonka/src/libs/locker"
)

const FILE_PASSWD = "/etc/passwd"

type Entries []Entry

// Unmarshal will unmarshal a provided passwd formatted file.
func Unmarshal(data []byte, dest interface{}) error {
	switch dest.(type) {
	case *Entries:
		break
	default:
		return errors.New("must unmarshal to pointer of passwd.Entries")
	}

	outfile := dest.(*Entries)
	file := strings.TrimSpace(string(data))
	lines := strings.Split(file, "\n")

	for _, line := range lines {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		var entryErrors []error
		var info, homedir, shell string

		// Split the lines on the delim, ":".
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			entryErrors = append(entryErrors, errors.New("passwd entry has less than 7 segments"))
		}

		// Check if username is provided or not.
		username := parts[0]
		if len(parts[0]) < 1 {
			entryErrors = append(entryErrors, errors.New("invalid username parsed"))
		}

		// Check if password is provided or not.
		password := parts[1]
		if parts[1] != "x" {
			entryErrors = append(entryErrors, errors.New("password not stored in /etc/shadow"))
		}

		// Check if uid is a valid int or not.
		uid, err := strconv.Atoi(parts[2])
		if err != nil {
			entryErrors = append(entryErrors, errors.New("invalid uid"))
		}

		// Check if the gid is a valid int.
		gid, err := strconv.Atoi(parts[3])
		if err != nil {
			entryErrors = append(entryErrors, errors.New("invalid gid"))
		}

		// Check if the info field is provided or not.
		if len(parts) >= 5 {
			info = parts[4]
		} else {
			entryErrors = append(entryErrors, errors.New("invalid info field"))
		}

		// Check if the homedir field is provided or not.
		if len(parts) >= 6 {
			homedir = parts[5]
		} else {
			entryErrors = append(entryErrors, errors.New("invalid homedir"))
		}

		if len(parts) == 7 {
			shell = parts[6]
		} else {
			entryErrors = append(entryErrors, errors.New("invalid shell"))
		}

		// Populate the new entry.
		pwdentry := Entry{
			Username: username,
			Password: password,
			UID:      uid,
			GID:      gid,
			Info:     info,
			HomeDir:  homedir,
			Shell:    shell,
			Error:    entryErrors,
		}

		//passwd = append(passwd, pwdentry)
		*outfile = append(*outfile, pwdentry)
	}

	dest = outfile
	return nil
}

// Marshal is a helper for passwd.Marshal().
func (e Entries) Marshal() ([]byte, error) {
	return Marshal(e)
}

// Marshal will parse the provided entries into a byte array for writing.
func Marshal(in Entries) ([]byte, error) {
	var out []string

	// Loop through the entries
	for _, entry := range in {
		// Check for username, uid and gid. Return if there is an error.
		// We check -2 for uid and gid since on macOS, -2 is an unprivileged user.
		if len(entry.Username) == 0 || entry.UID < -2 || entry.GID < -2 {
			return nil, errors.New("attempting to save invalid entry")
		}

		// Append this entry to the outslice
		out = append(out, fmt.Sprintf(
			"%s:%s:%d:%d:%s:%s:%s",
			entry.Username,
			entry.Password,
			entry.UID,
			entry.GID,
			entry.Info,
			entry.HomeDir,
			entry.Shell,
		))
	}

	// Join the slices by new lines.
	return []byte(strings.Join(out, "\n")), nil
}

// Save will take in entries.
func (e Entries) Save() error {
	b, err := e.Marshal()
	if err != nil {
		return err
	}

	// Will write the entries list with.
	if err = locker.WriteWithLock(FILE_PASSWD, b); err != nil {
		return err
	}

	return nil
}

// LoadFromDisk will read an /etc/passwd file and return parsed Entries or error.
func LoadFromDisk() (*Entries, error) {
	b, err := ioutil.ReadFile(FILE_PASSWD)
	if err != nil {
		return nil, err
	}

	var e Entries
	err = Unmarshal(b, &e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

// NewEntry adds a new entry to Entries.
func (e *Entries) NewEntry(new Entry) {
	*e = append(*e, new)
}

// RemoveEntry adds a new entry to Entries.
func (e *Entries) RemoveEntry(rm Entry) error {
	if len(rm.Username) == 0 {
		return errors.New("must provide username to be removed")
	}

	// Look for the username, then remove it.
	for i, entry := range *e {
		if entry.Username == rm.Username {
			s := *e
			s = append(s[:i], s[i+1:]...)
			*e = s
			return nil
		}
	}

	// return an error if it's not found
	return &ErrNotFound{"entry not found"}
}

// GetUser will get a user by name.
func (e *Entries) GetUser(name string) *Entry {
	for _, user := range *e {
		if user.Username == name {
			return &user
		}
	}

	return nil
}

// GetUserByID will get a user by id.
func (e *Entries) GetUserByID(id int) *Entry {
	for _, user := range *e {
		if user.UID == id {
			return &user
		}
	}

	return nil
}
