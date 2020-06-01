package locker

import "testing"

func TestWriteWithLock(t *testing.T) {
	filename := "/tmp/test"
	err := WriteWithLock(filename, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
}
