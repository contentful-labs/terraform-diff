package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"os"
	"os/exec"
	"path"
	"strings"
)

type cfg struct {
	gitCommitRange, outputFormat string
	folders                      []string
}

// Returns all modules used by the Terraform project in folder
func getTerraformDependencies(folder string) (map[string]*tfconfig.Module, error) {
	modules := map[string]*tfconfig.Module{}

	module, diags := tfconfig.LoadModule(folder)
	if diags != nil {
		return map[string]*tfconfig.Module{}, fmt.Errorf("Error processing: %s\n", diags)
	}

	modules[folder] = module

	for _, res := range module.ModuleCalls {
		expandedModulePath := path.Clean(folder + "/" + res.Source)

		deps, err := getTerraformDependencies(expandedModulePath)
		if err != nil {
			return map[string]*tfconfig.Module{}, err
		}
		for k, v := range deps {
			modules[k] = v
		}
	}

	return modules, nil
}

// If paths are relative to a subfolder of a GIT repo, expand them to start from the Git root folder
func expandPathFromGitRoot(folder string, wd, gitRoot string) string {
	pathFromGitRoot := strings.TrimPrefix(wd, gitRoot)
	pathFromGitRoot = strings.TrimPrefix(pathFromGitRoot, string(os.PathSeparator))

	folder = strings.TrimPrefix(folder, gitRoot)
	folder = strings.TrimPrefix(folder, string(os.PathSeparator))
	return path.Join(pathFromGitRoot, folder)
}

// Given a list of changed files, return a list of changed folders
func foldersContainingFiles(changedFiles []string) ([]string, error) {
	// easier & more efficient to use a map here to enforce uniqueness
	changedFoldersMap := map[string]bool{}
	for _, file := range changedFiles {
		folder := path.Dir(file)
		changedFoldersMap[folder] = true
	}

	changedFolders := []string{}
	for folder, _ := range changedFoldersMap {
		changedFolders = append(changedFolders, folder)
	}

	return changedFolders, nil
}

// Returns the files that have been changed in the commitRange
func filesInGitDiff(commitRange string) ([]string, error) {
	params := []string{"diff", "--name-only"}
	if len(commitRange) > 0 {
		params = append(params, commitRange)
	}

	out, err := exec.Command("git", params...).Output()
	if err != nil {
		return nil, fmt.Errorf("Error running git: %s", err)
	}

	files := []string{}
	for _, file := range strings.Split(string(out), "\n") {
		if len(file) > 0 {
			files = append(files, file)
		}
	}

	return files, nil
}

func getOutput(foldersToPlan []string, format string) (string, error) {
	output := ""

	if format == "json" {
		m, err := json.MarshalIndent(foldersToPlan, "", " ")
		if err != nil {
			return "", fmt.Errorf("failed to marshall result: %s", err)
		}
		output = string(m)
	} else {
		for _, folderToPlan := range foldersToPlan {
			output = output + folderToPlan + "\n"
		}
	}

	return output, nil
}

func getFoldersToPlan(tfDeps map[string][]string, changedFolders []string) []string {
	foldersToPlan := []string{}
	for folder, deps := range tfDeps {
		for _, dep := range deps {
			for _, cFolder := range changedFolders {
				if dep == cFolder {
					foldersToPlan = append(foldersToPlan, folder)
				}
			}
		}
	}

	return foldersToPlan
}

func getConfig() (*cfg, error) {
	gitCommitRange := flag.String("range", "", "git commit range")
	outputFormat := flag.String("output", "text", "output format (text or json)")
	flag.Parse()

	format := strings.ToLower(*outputFormat)

	if format != "text" && format != "json" {
		return nil, fmt.Errorf("-output parameter can only be text or json")
	}

	// remove all trailing slashes
	folders := []string{}
	for _, f := range flag.Args() {
		folders = append(folders, strings.TrimSuffix(f, string(os.PathSeparator)))
	}

	return &cfg{
		gitCommitRange: *gitCommitRange,
		outputFormat:   format,
		folders:        folders,
	}, nil
}

func realMain() int {
	cfg, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing command line: %s\n", err)
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting working directory: %s\n", err)
		return 1
	}
	cwd = strings.TrimSpace(cwd)
	cwd = strings.TrimSuffix(cwd, string(os.PathSeparator))

	cmdOutput, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting git top level directory (trying to run outside of a git repository?): %s\n", err)
		return 1
	}
	gitRoot := strings.TrimSpace(string(cmdOutput))
	gitRoot = strings.TrimSuffix(gitRoot, string(os.PathSeparator))

	changedFiles, err := filesInGitDiff(cfg.gitCommitRange)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed getting changed files from git diff: %s\n", err)
		return 1
	}

	changedFolders, err := foldersContainingFiles(changedFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed getting changed folders from git diff: %s\n", err)
		return 1
	}

	tfDeps := map[string][]string{} // List of TF modules each folder depends on
	for _, folder := range cfg.folders {
		modules, err := getTerraformDependencies(folder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed getting dependencies for folder %s: %s\n", folder, err)
			return 1
		}

		for moduleFolder, _ := range modules {
			tfDeps[folder] = append(tfDeps[folder], expandPathFromGitRoot(moduleFolder, cwd, gitRoot))
		}
	}

	if output, err := getOutput(getFoldersToPlan(tfDeps, changedFolders), cfg.outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	} else {
		fmt.Print(output)
	}

	return 0
}

func main() {
	os.Exit(realMain())
}
