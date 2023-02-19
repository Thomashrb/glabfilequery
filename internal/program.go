package internal

import (
	"fmt"
	"glabfilequery/internal/fs"
	"glabfilequery/internal/gitlab"
	"glabfilequery/internal/tui"
	"os"

	bt "github.com/charmbracelet/bubbletea"
)

func Run(
	baseurl string,
	token string,
	fileSuffix string,
	outputDir string,
	dryRun bool,
) {
	tea := bt.NewProgram(tui.NewModel())

	go func() {
		err, projects := gitlab.ListProjects(baseurl, token, tea)
		if err != nil {
			fmt.Printf("Could not list projects from %s, error: %s", baseurl, err)
		}

		err, projectFiles := gitlab.ListProjectFiles(baseurl, token, fileSuffix, projects, tea)
		if err != nil {
			fmt.Printf("Could not list files from %s, error: %s", baseurl, err)
		}

		if !dryRun {
			err, files := gitlab.GetFiles(baseurl, token, projectFiles, tea)
			if err != nil {
				fmt.Printf("Could not pull files from %s, error: %s", baseurl, err)
			}
			err = fs.SaveFiles(files, outputDir, tea)
			if err != nil {
				fmt.Printf("Could not save files error: %s", err)
			}
		}

		os.Exit(0)
	}()

	if _, err := tea.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
