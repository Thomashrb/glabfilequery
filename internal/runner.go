package internal

import (
	"fmt"
	"glabfilequery/internal/fs"
	"glabfilequery/internal/gitlab"
	"glabfilequery/internal/tui"
	"os"
	"regexp"
)

func Run(
	baseurl string,
	token string,
	fileRegex *regexp.Regexp,
	outputDir string,
	dryRun bool,
) {
	p := tui.NewProgram()

	go func() {
		p.StageMsgSend("Listing projects")
		err, projects := gitlab.ListProjects(baseurl, token)
		if err != nil {
			fmt.Printf("Could not list projects from %s, error: %s", baseurl, err)
		}

		p.StageMsgSend("Listing project files")
		err, projectFiles := gitlab.ListProjectFiles(baseurl, token, fileRegex, projects)
		if err != nil {
			fmt.Printf("Could not list files from %s, error: %s", baseurl, err)
		}

		if !dryRun {
			p.StageMsgSend("Downloading files")
			err, files := gitlab.GetFiles(baseurl, token, projectFiles)
			if err != nil {
				fmt.Printf("Could not pull files from %s, error: %s", baseurl, err)
			}

			p.StageMsgSend("Saving files to output directory")
			err = fs.SaveFiles(files, outputDir)
			if err != nil {
				fmt.Printf("Could not save files error: %s", err)
			}
		}

		p.QuitMsgSend()
		p.Quit()
	}()

	if _, err := p.Program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
