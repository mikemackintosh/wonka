package groups

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mikemackintosh/wonka/src/libs/locker"
)

const FILE_GROUP = "/etc/group"

type Entries []*Group

type Group struct {
	Name     string
	Password string
	GID      int
	Users    []string
	Errors   []error
}

// Unmarshal will unmarshal a provided passwd formatted file.
func Unmarshal(data []byte, dest interface{}) error {
	switch dest.(type) {
	case *Entries:
		break
	default:
		return errors.New("must unmarshal to pointer of groups.Entries")
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
		var name, password string
		var gid int
		var users []string

		// Split the lines on the delim, ":".
		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			errs = append(errs, errors.New("shadow Group has less than 4 segments"))
		}

		// Check if name is provided or not.
		name = parts[0]
		if len(parts[0]) < 1 {
			errs = append(errs, errors.New("invalid name parsed"))
		}

		// Check if password is provided or not.
		password = parts[1]

		// Check if the gid is a valid int.
		gid, err := strconv.Atoi(parts[2])
		if err != nil {
			errs = append(errs, errors.New("invalid gid"))
		}

		// Split the lines on the delim, ":".
		users = strings.Split(parts[3], ",")

		// Populate the new Group.
		Group := &Group{
			Name:     name,
			Password: password,
			GID:      gid,
			Users:    users,
			Errors:   errs,
		}

		//passwd = append(passwd, pwdGroup)
		*outfile = append(*outfile, Group)
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
	for _, group := range in {
		// Check for username, uid and gid. Return if there is an error.
		// We check -2 for uid and gid since on macOS, -2 is an unprivileged user.
		if len(group.Name) == 0 {
			return nil, errors.New("attempting to save invalid Group")
		}

		// Append this Group to the outslice
		out = append(out, fmt.Sprintf(
			"%s:%s:%d:%s",
			group.Name,
			group.Password,
			group.GID,
			strings.Join(group.Users, ","),
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
	if err = locker.WriteWithLock(FILE_GROUP, b); err != nil {
		return err
	}

	return nil
}

// LoadFromDisk will read an /etc/passwd file and return parsed Entries or error.
func LoadFromDisk() (*Entries, error) {
	b, err := ioutil.ReadFile(FILE_GROUP)
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

// NewGroup adds a new Group to Entries.
func (e *Entries) NewGroup(new *Group) {
	*e = append(*e, new)
}

// RemoveGroup adds a new Group to Entries.
func (e *Entries) RemoveGroup(rm *Group) error {
	if len(rm.Name) == 0 {
		return errors.New("must provide name to be removed")
	}

	// Look for the username, then remove it.
	for i, group := range *e {
		if group.Name == rm.Name {
			s := *e
			s = append(s[:i], s[i+1:]...)
			*e = s
			return nil
		}
	}

	// return an error if it's not found
	return &ErrNotFound{"Group not found"}
}

// GetGroup will get a group by name.
func (e *Entries) GetGroup(name string) *Group {
	for _, Group := range *e {
		if Group.Name == name {
			return Group
		}
	}

	return nil
}

// GetGroup will get a group by id.
func (e *Entries) GetGroupByID(id int) *Group {
	for _, Group := range *e {
		if Group.GID == id {
			return Group
		}
	}

	return nil
}

// TODO: add check for existing user
func (g *Group) AddUser(name string) {
	g.Users = append(g.Users, name)
}

// TODO: add check for existing user.
func (g *Group) RemoveUser(name string) error {
	// Look for the username, then remove it.
	for i, user := range g.Users {
		if user == name {
			g.Users = append(g.Users[:i], g.Users[i+1:]...)
			return nil
		}
	}

	return &ErrNotFound{"user not found in group"}
}
