package wonka

import (
	"fmt"
	"io/ioutil"
)

const (
	defaultFilePasswd = "/etc/passwd"
	defaultFileGroups = "/etc/groups"
	defaultFileShadow = "/etc/shadow"
)

type Options struct {
	filePasswd string
	fileGroups string
	fileShadow string
	expiry     int
}

type Instance struct {
	Options Options
}

func New() Instance {
	return NewWithOptions(Options{
		filePasswd: defaultFilePasswd,
		fileGroups: defaultFileGroups,
		fileShadow: defaultFileShadow,
		expiry:     0,
	})
}

func NewWithOptions(options Options) Instance {
	return Instance{}
}

type ErrListFileFailed struct {
	err  error
	file string
}

func (e *ErrListFileFailed) Error() string {
	return fmt.Sprintf("failed to list %s, %s", e.file, e.err)
}

func List(file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return &ErrListFileFailed{err, file}
	}

	fmt.Printf("%s", f)

	return nil
}
