package files

import (
	"fmt"
	"io/fs"
	"os"
)

func OpenFile(path string) (*os.File, error) {
	ok := fs.ValidPath(path)

	if !ok {
		return nil, fmt.Errorf("bad path provided %q", path)
	}
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	return f, nil
}
