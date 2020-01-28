package main

import (
	"testing"
)

func TestExpandPathFromGitRoot(t *testing.T) {
	testCases := []struct {
		paths       string
		wd, gitRoot string
		want        string
	}{
		{
			"env/production",
			"/home/user/repo/terraform",
			"/home/user/repo",
			"terraform/env/production",
		},
		{
			"terraform/env/production",
			"/home/user/repo",
			"/home/user/repo",
			"terraform/env/production",
		},
		{
			"/home/user/repo/terraform/env/production",
			"/home/user/repo",
			"/home/user/repo",
			"terraform/env/production",
		},
	}

	for i, testCase := range testCases {
		path := expandPathFromGitRoot(testCase.paths, testCase.wd, testCase.gitRoot)
		if path != testCase.want {
			t.Errorf("failed test %d: got %s, wanted %s\n", i, path, testCase.want)
		}
	}
}

func TestGetFoldersToPlan(t *testing.T) {

	testCases := []struct {
		deps           map[string][]string
		changedFolders []string
		want           []string
	}{
		{
			map[string][]string{
				"terraform/staging":    []string{"modules/1", "modules/2"},
				"terraform/production": []string{"modules/1", "modules/2", "modules/3"},
			},
			[]string{"modules/3"},
			[]string{"terraform/production"},
		},
		{
			map[string][]string{
				"terraform/staging":    []string{"modules/1", "modules/2"},
				"terraform/production": []string{"modules/1", "modules/2", "modules/3"},
			},
			[]string{"modules/3", "modules/1"},
			[]string{"terraform/production", "terraform/staging"},
		},
		{
			map[string][]string{
				"terraform/staging":    []string{"modules/1", "modules/2"},
				"terraform/production": []string{"modules/1", "modules/2", "modules/3"},
			},
			[]string{},
			[]string{},
		},
		{
			map[string][]string{
				"terraform/staging":    []string{"modules/1", "modules/2"},
				"terraform/production": []string{"modules/1", "modules/2", "modules/3"},
			},
			[]string{"modules/5"},
			[]string{},
		},
		{
			map[string][]string{
				"terraform/staging":    []string{"modules/1", "modules/2"},
				"terraform/production": []string{"modules/1", "modules/2", "modules/3"},
			},
			[]string{"modules/1", "modules/2", "modules/3"},
			[]string{"terraform/staging", "terraform/production"},
		},
	}

	for testI, testCase := range testCases {
		toPlan := getFoldersToPlan(testCase.deps, testCase.changedFolders)
		for _, folder := range toPlan {
			found := false
			for _, f := range testCase.want {
				if f == folder {
					found = true
				}
			}
			if found == false {
				t.Errorf("failed test %d, got %+v, wanted %+v", testI, toPlan, testCase.want)
			}
		}
	}
}
