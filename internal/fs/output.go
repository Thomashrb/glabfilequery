package fs

import (
	"fmt"
	"net/url"
	"os"
)

func SaveFiles(files map[string][]byte, dirpath string) error {
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating directory for saving files: %s", err)
		}
	}

	for fn, fbytes := range files {
		err := os.WriteFile(fmt.Sprintf("%s/%s", dirpath, url.PathEscape(fn)), fbytes, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
