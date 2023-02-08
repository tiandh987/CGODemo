package file

import "testing"

func TestPathExists(t *testing.T) {
	exists, err := PathExists("./file.go")
	t.Logf("exists: %t, err: %+v", exists, err)

	exists, err = PathExists("./test.go")
	t.Logf("exists: %t, err: %+v", exists, err)
}
