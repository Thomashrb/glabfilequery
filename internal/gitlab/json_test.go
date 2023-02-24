package gitlab

import (
	"testing"
)

const (
	projectsJson = `[
  {
    "id": 304,
    "description": null,
    "name": "projectname",
    "name_with_namespace": "Name Nameson / projectname",
    "path": "projectname",
    "path_with_namespace": "name.nameson/projectname",
    "created_at": "2022-12-10T15:17:06.847Z",
    "default_branch": "develop",
    "ssh_url_to_repo": "git@gitlab.companyname.no:name.nameson/projectname.git",
    "http_url_to_repo": "https://gitlab.companyname.no/name.nameson/projectname.git",
    "web_url": "https://gitlab.companyname.no/name.nameson/projectname",
    "readme_url": "https://gitlab.companyname.no/name.nameson/projectname/-/blob/develop/README.md",
    "avatar_url": null,
    "forks_count": 0,
    "star_count": 0,
    "last_activity_at": "2023-02-09T12:48:18.873Z",
    "packages_enabled": true,
    "empty_repo": false,
    "archived": false,
    "visibility": "private"
	}
]`
	filesJson = `[
  {
    "id": "690d5813be90cb27f7eef663bcfa6dec5625b538",
    "name": ".gitignore",
    "type": "blob",
    "path": ".gitignore",
    "mode": "100644"
  },
  {
    "id": "a71b913960ad019549989037448da91b1d964052",
    "name": "README.md",
    "type": "blob",
    "path": "README.md",
    "mode": "100644"
  }
]`
)

func TestToProjectsWithoutError(t *testing.T) {
	err, _ := ToProjects([]byte(projectsJson))
	if err != nil {
		t.Fatalf("ToProjects failed to parse test json")
	}
}

func TestToFilesWithoutError(t *testing.T) {
	err, _ := ToFiles([]byte(filesJson))
	if err != nil {
		t.Fatalf("ToFiles failed to parse test json")
	}
}

func TestToProjectsCorrectlyParsesTestInput(t *testing.T) {
	_, ps := ToProjects([]byte(projectsJson))
	p := ps[0]
	if p.Id != 304 {
		t.Fatalf("ToProjects produces the wrong Id: %d", p.Id)
	}
	if p.Archived != false {
		t.Fatalf("ToProjects produces the wrong Archived: %t", p.Archived)
	}
	if p.Visibility != "private" {
		t.Fatalf("ToProjects produces the wrong Visibility: %s", p.Visibility)
	}
	if p.DefaultBranch != "develop" {
		t.Fatalf("ToProjects produces the wrong DefaultBranch: %s", p.DefaultBranch)
	}
	if p.WebUrl != "https://gitlab.companyname.no/name.nameson/projectname" {
		t.Fatalf("ToProjects produces the wrong WebUrl: %s", p.WebUrl)
	}
	if p.LastActivityAt != "2023-02-09T12:48:18.873Z" {
		t.Fatalf("ToProjects produces the wrong LastActivityAt: %s", p.LastActivityAt)
	}
}

func TestToFilesCorrectlyParsesTestInput(t *testing.T) {
	_, fs := ToFiles([]byte(filesJson))
	f := fs[0]
	if f.Id != "690d5813be90cb27f7eef663bcfa6dec5625b538" {
		t.Fatalf("ToFiles produces the wrong Id: %s", f.Id)
	}
	if f.Type != "blob" {
		t.Fatalf("ToFiles produces the wrong Type: %s", f.Type)
	}
	if f.Path != ".gitignore" {
		t.Fatalf("ToFiles produces the wrong Path: %s", f.Path)
	}
}
