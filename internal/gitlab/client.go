package gitlab

import (
	"fmt"
	"glabfilequery/internal/tui"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

func ListProjects(baseurl string, token string, pg tui.Program) (error, []Project) {
	pg.StageMsgSend("Listing projects")

	var ps []Project
	page := 0
	for {
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

func ListProjectFiles(baseurl string, token string, fileRegex *regexp.Regexp, projects []Project, pg tui.Program) (error, map[Project]File) {
	wg := sync.WaitGroup{}
	projectFiles := make(map[Project]File)
	var filteredPs []Project

	for _, p := range projects {
		if p.Archived {
			continue
		}
		wg.Add(1)
		filteredPs = append(filteredPs, p)
	}

	pg.StageMsgSend("Listing project files")
	for _, p := range filteredPs {
		go func(p Project) {
			defer wg.Done()
			err, files := listFiles(baseurl, token, p.Id)
			if err == nil {
				for _, f := range files {
					if f.Type != "blob" {
						continue
					}
					if !fileRegex.Match([]byte(f.Name)) {
						continue
					}
					projectFiles[p] = f
				}
			}
		}(p)
	}

	wg.Wait()
	return nil, projectFiles
}

func GetFiles(baseurl string, token string, projectFiles map[Project]File, pg tui.Program) (error, map[string][]byte) {
	wg := sync.WaitGroup{}
	wg.Add(len(projectFiles))
	nameFiles := make(map[string][]byte)

	pg.StageMsgSend("Downloading files")
	for p, f := range projectFiles {
		blobPath := fmt.Sprintf("%s/-/blob/%s/%s", p.WebUrl, p.DefaultBranch, f.Name)
		go func(p Project, f File) {
			defer wg.Done()
			err, bytes := getRaw(baseurl, token, p.Id, f.Id)
			if err == nil {
				nameFiles[blobPath] = bytes
			}
		}(p, f)
	}

	wg.Wait()
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
