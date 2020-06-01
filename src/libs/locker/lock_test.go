package locker

import (
	"os"
	"testing"
)

func TestWriteWithLock(t *testing.T) {
	path := "/tmp/test"

	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
	}

	err = WriteWithLock(path, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatal(err)
	}
}
