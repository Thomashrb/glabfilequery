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
	fileRegex string,
	outputDir string,
	dryRun bool,
) {
	p := tui.NewProgram()

	go func() {
		err, projects := gitlab.ListProjects(baseurl, token, p)
		if err != nil {
			fmt.Printf("Could not list projects from %s, error: %s", baseurl, err)
		}

		re := regexp.MustCompile(fileRegex)
		err, projectFiles := gitlab.ListProjectFiles(baseurl, token, re, projects, p)
		if err != nil {
			fmt.Printf("Could not list files from %s, error: %s", baseurl, err)
		}

		if !dryRun {
			err, files := gitlab.GetFiles(baseurl, token, projectFiles, p)
			if err != nil {
				fmt.Printf("Could not pull files from %s, error: %s", baseurl, err)
			}
			err = fs.SaveFiles(files, outputDir, p)
			if err != nil {
				fmt.Printf("Could not save files error: %s", err)
			}
		}

		os.Exit(0)
	}()

	if _, err := p.Program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
