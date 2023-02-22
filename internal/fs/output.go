package fs

import (
	"fmt"
	"glabfilequery/internal/tui"
	"net/url"
	"os"
)

func SaveFiles(files map[string][]byte, dirpath string, pg tui.Program) error {
	pg.StageMsgSend("Saving files to output directory")
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
