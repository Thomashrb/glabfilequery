package gitlab

import (
	"encoding/json"
	"fmt"
)

type (
	Project struct {
		Id             int
		Archived       bool
		Visibility     string
		DefaultBranch  string `json:"default_branch"`
		WebUrl         string `json:"web_url"`
		LastActivityAt string `json:"last_activity_at"`
	}

	File struct {
		Id   string
		Path string
		Type string
	}
)

func ToProjects(pjson []byte) (error, []Project) {
	var ps []Project
	err := json.Unmarshal(pjson, &ps)
	if err != nil {
		return fmt.Errorf("ToProjects failed with: %s", err), nil
	}
	return nil, ps
}

func ToFiles(fjson []byte) (error, []File) {
	var fs []File
	err := json.Unmarshal(fjson, &fs)
	if err != nil {
		return fmt.Errorf("ToFiles failed with: %s", err), nil
	}
	return nil, fs
}
