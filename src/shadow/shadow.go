package shadow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/mikemackintosh/wonka/src/libs/locker"
	r "github.com/mikemackintosh/wonka/src/libs/rand"
	"github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

const FILE_SHADOW = "/etc/shadow"

type Entries []*Entry

type Entry struct {
	Username           string
	Password           string
	PasswordUpdated    bool
	LastPasswordChange time.Time
	MinimumPasswordAge *time.Duration
	MaximumPasswordAge *time.Duration
	WarningPeriod      *time.Duration
	InactivityPeriod   *time.Duration
	ExpirationPeriod   *time.Duration
	Unused             interface{}
	Errors             []error
}

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

		var errs []error
		var lastChange time.Time
		var day = time.Duration(24) * time.Hour

		// Split the lines on the delim, ":".
		parts := strings.Split(line, ":")
		if len(parts) < 8 {
			errs = append(errs, errors.New("shadow entry has less than 8 segments"))
		}

		// Populate the new entry.
		entry := &Entry{}

		// Check if username is provided or not.
		entry.Username = parts[0]
		if len(parts[0]) < 1 {
			errs = append(errs, errors.New("invalid username parsed"))
		}

		// Check if password is provided or not.
		entry.Password = parts[1]

		// Check if lastPasswordChange is a valid int or not.
		if len(parts[2]) > 0 {
			lastChangeDays, err := strconv.Atoi(parts[2])
			if err != nil {
				errs = append(errs, errors.New("invalid lastPasswordChange"))
			}
			fd := (time.Duration(lastChangeDays) * time.Hour * 24)
			lastChange = time.Unix(0, 0).Add(fd)
			entry.LastPasswordChange = lastChange
		}

		// Check if the minAge is a valid int.
		if len(parts[3]) > 0 {
			minAgeField, err := strconv.Atoi(parts[3])
			if err != nil {
				errs = append(errs, errors.New("invalid minAge"))
			}
			minAge := time.Duration(minAgeField) * day
			entry.MinimumPasswordAge = &minAge
		}

		// Check if the maxAge is a valid int.
		if len(parts[4]) > 0 {
			maxAgeField, err := strconv.Atoi(parts[4])
			if err != nil {
				errs = append(errs, errors.New("invalid maxAge"))
			}
			maxAge := time.Duration(maxAgeField) * day
			entry.MaximumPasswordAge = &maxAge
		}

		// Check if the warning is a valid int.
		if len(parts[5]) > 0 {
			warningField, err := strconv.Atoi(parts[5])
			if err != nil {
				errs = append(errs, errors.New("invalid warning"))
			}
			warning := time.Duration(warningField) * day
			entry.WarningPeriod = &warning
		}

		// Check if the gid is a valid int.
		if len(parts[6]) > 0 {
			inactivityField, err := strconv.Atoi(parts[6])
			if err != nil {
				errs = append(errs, errors.New("invalid inactivity"))
			}
			inactivity := time.Duration(inactivityField) * day
			entry.InactivityPeriod = &inactivity
		}

		if len(parts[7]) > 0 {
			expirationField, err := strconv.Atoi(parts[7])
			if err != nil {
				errs = append(errs, errors.New("invalid expiration"))
			}
			expiration := time.Duration(expirationField) * day
			entry.ExpirationPeriod = &expiration
		}

		if len(errs) > 0 {
			entry.Errors = errs
		}

		//passwd = append(passwd, pwdentry)
		*outfile = append(*outfile, entry)
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
	var err error

	// Loop through the entries
	for _, entry := range in {
		var line []string
		// Check for username, uid and gid. Return if there is an error.
		// We check -2 for uid and gid since on macOS, -2 is an unprivileged user.
		if len(entry.Username) == 0 {
			return nil, errors.New("attempting to save invalid entry")
		}
		line = append(line, entry.Username)

		var password = entry.Password
		if entry.PasswordUpdated {
			cryptor := sha512_crypt.New()

			password, err = cryptor.Generate([]byte(entry.Password), []byte("$6$"+r.String(8)))
			if err != nil {
				return nil, fmt.Errorf("error generating password, %s", err)
			}

			entry.LastPasswordChange = time.Now()
		}
		line = append(line, password)

		line = append(line, fmt.Sprintf("%d", int(entry.LastPasswordChange.Sub(time.Unix(0, 0)).Hours()/24)))

		if entry.MinimumPasswordAge != nil {
			line = append(line, fmt.Sprintf("%d", int(entry.MinimumPasswordAge.Hours()/24)))
		} else {
			line = append(line, "")
		}

		if entry.MaximumPasswordAge != nil {
			line = append(line, fmt.Sprintf("%d", int(entry.MaximumPasswordAge.Hours()/24)))
		} else {
			line = append(line, "")
		}

		if entry.WarningPeriod != nil {
			line = append(line, fmt.Sprintf("%d", int(entry.WarningPeriod.Hours()/24)))
		} else {
			line = append(line, "")
		}

		if entry.InactivityPeriod != nil {
			line = append(line, fmt.Sprintf("%d", int(entry.InactivityPeriod.Hours()/24)))
		} else {
			line = append(line, "")
		}

		if entry.ExpirationPeriod != nil {
			line = append(line, fmt.Sprintf("%d", int(entry.ExpirationPeriod.Hours()/24)))
		} else {
			line = append(line, "")
		}

		// Unused segment
		line = append(line, "")

		out = append(out, strings.Join(line, ":"))
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
	if err = locker.WriteWithLock(FILE_SHADOW, b); err != nil {
		return err
	}

	return nil
}

// LoadFromDisk will read an /etc/passwd file and return parsed Entries or error.
func LoadFromDisk() (*Entries, error) {
	b, err := ioutil.ReadFile(FILE_SHADOW)
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
func (e *Entries) NewEntry(new *Entry) {
	*e = append(*e, new)
}

// RemoveEntry adds a new entry to Entries.
func (e *Entries) RemoveEntry(rm *Entry) error {
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

// UpdatePassword will update the password.
func (e *Entry) UpdatePassword(password string) {
	e.PasswordUpdated = true
	e.Password = password
}

//GetUserEntry will search the entries list for a user.
func (e *Entries) GetUserEntry(user string) *Entry {
	for _, entry := range *e {
		if entry.Username == user {
			return entry
		}
	}

	return nil
}
