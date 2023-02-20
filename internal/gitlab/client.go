package gitlab

import (
	"fmt"
	"glabfilequery/internal/tui"
	"io/ioutil"
	"net/http"
	"strings"
)

func ListProjects(baseurl string, token string, pg tui.Program) (error, []Project) {
	pg.StageMsgSend("Listing projects")

	var ps []Project
	page := 0
	for {
		pg.JobMsgSend(fmt.Sprintf("Queriying project page %d", page))
		err, body := authenticatedGetReq(fmt.Sprintf("%s/api/v4/projects?page=%d", baseurl, page), token)
		if err != nil {
			return err, nil
		}
		err, psPage := ToProjects(body)
		if err != nil {
			return err, nil
		}
		if len(psPage) <= 0 {
			break
		}
		ps = append(ps, psPage...)
		page++
	}
	return nil, ps
}

func ListProjectFiles(baseurl string, token string, fileSuffix string, projects []Project, pg tui.Program) (error, map[Project]File) {
	pg.StageMsgSend("Listing project files")
	projectFiles := make(map[Project]File)

	for _, p := range projects {
		pg.JobMsgSend(fmt.Sprintf("Quering file tree for %s", p.WebUrl))
		if p.Archived {
			continue
		}
		err, files := listFiles(baseurl, token, p.Id)
		if err != nil {
			return err, nil
		}
		for _, f := range files {
			if f.Type != "blob" {
				continue
			}
			if !strings.HasSuffix(f.Name, fileSuffix) {
				continue
			}
			projectFiles[p] = f
		}
	}
	return nil, projectFiles
}

func GetFiles(baseurl string, token string, projectFiles map[Project]File, pg tui.Program) (error, map[string][]byte) {
	pg.StageMsgSend("Downloading files")

	nameFiles := make(map[string][]byte)

	for p, f := range projectFiles {
		blobPath := fmt.Sprintf("%s/-/blob/%s/%s", p.WebUrl, p.DefaultBranch, f.Name)
		pg.JobMsgSend(blobPath)
		err, bytes := getRaw(baseurl, token, p.Id, f.Id)
		if err != nil {
			return err, nil
		}
		nameFiles[blobPath] = bytes
	}
	return nil, nameFiles
}

func authenticatedGetReq(url string, token string) (error, []byte) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err, nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("\nGitlab response failed with %d", resp.StatusCode)
		return nil, []byte("[]")
	}
	if err != nil {
		return err, nil
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	return nil, resBody
}

func listFiles(baseurl string, token string, projecId int) (error, []File) {
	var fs []File
	page := 0
	for {
		err, body := authenticatedGetReq(fmt.Sprintf("%s/api/v4/projects/%d/repository/tree?page=%d", baseurl, projecId, page), token)
		if err != nil {
			return fmt.Errorf("Listing files failed %s", err), nil
		}
		err, fsPage := ToFiles(body)
		if err != nil {
			return err, nil
		}
		if len(fsPage) <= 0 {
			break
		}
		fs = append(fs, fsPage...)
		page++
	}
	return nil, fs
}

func getRaw(baseurl string, token string, projectId int, fileId string) (error, []byte) {
	return authenticatedGetReq(fmt.Sprintf("%s/api/v4/projects/%d/repository/blobs/%s/raw", baseurl, projectId, fileId), token)
}
