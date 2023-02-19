package main

import (
	"flag"
	"glabfilequery/internal"
)

var (
	baseurl    string
	token      string
	fileSuffix string
	outputDir  string
	dryRun     bool
)

func main() {
	flag.StringVar(&baseurl, "baseuri", "", "base uri for upstream gitlab instance ex: gitlab.{company}.com")
	flag.StringVar(&token, "token", "", "authentication token for upstream gitlab instance")
	flag.StringVar(&fileSuffix, "filesuffix", ".md", "suffix of files to look for")
	flag.StringVar(&outputDir, "outputdir", "output", "directory to store output in")
	flag.BoolVar(&dryRun, "dryrun", false, "do not store any output")
	flag.Parse()

	internal.Run(baseurl, token, fileSuffix, outputDir, dryRun)
}
