package fileutils

import (
	"os"
)

func Delete(files ...*os.File) (err error) {
	for _, file := range files {
		file.Close()
		if err = os.Remove(file.Name()); err != nil {
			return
		}
	}
	return
}
