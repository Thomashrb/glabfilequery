package fs

import (
	"fmt"
	"glabfilequery/internal/tui"
	"net/url"
	"os"

	bt "github.com/charmbracelet/bubbletea"
)

func SaveFiles(files map[string][]byte, dirpath string, tea *bt.Program) error {
	tea.Send(tui.StageMsg("Saving files to output directory"))
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		tea.Send(tui.JobMsg("Creating directory"))
		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error creating directory for saving files: %s", err)
		}
	}

	tea.Send(tui.JobMsg(fmt.Sprintf("Saving %d files", len(files))))
	for fn, fbytes := range files {
		err := os.WriteFile(fmt.Sprintf("%s/%s", dirpath, url.PathEscape(fn)), fbytes, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
