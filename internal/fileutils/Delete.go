package fileutils

import (
	"os"
)

func Delete(paths ...string) (err error) {
	for _, path := range paths {
		if err = os.Remove(path); err != nil {
			return
		}
	}
	return
}
