package locker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
)

func WriteWithLock(filename string, data []byte) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, fileInfo.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	err = syscall.Flock(int(f.Fd()), 2) //2 is exclusive lock
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to lock %s for reading.", filename))
	}
	defer func() {
		err := syscall.Flock(int(f.Fd()), 8) //8 is unlock
		if err != nil {
			log.Fatal(fmt.Sprintf("Unable to unlock %s for reading.", filename))
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to lock %s for reading.", filename))
	}

	return nil
}
