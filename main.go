package main

import (
	"encoding/json"
	"fmt"
	"os"
	"flag"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"path"
)

func getModulesCallsRecursive(folder string) (map[string]*tfconfig.Module, error) {
	modules := map[string]*tfconfig.Module{}

	module, diags := tfconfig.LoadModule(folder)
	if diags != nil {
		return map[string]*tfconfig.Module{}, fmt.Errorf("Error processing: %s\n", diags)
	}

	modules[folder] = module

	for _, res := range module.ModuleCalls {
		expandedModulePath := path.Clean(folder + "/" + res.Source)

		deps, err := getModulesCallsRecursive(expandedModulePath)
		if err != nil {
			return map[string]*tfconfig.Module{}, err
		}
		for k, v := range deps {
			modules[k] = v
		}
	}

	return modules, nil
}

func realMain() int {
	folderDeps := map[string][]string{}

	flag.Parse()
	for _, folder := range flag.Args() {
		modules, err := getModulesCallsRecursive(folder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Err: %s\n", err)
			continue
		}

		for moduleFolder, _ := range modules {
			folderDeps[folder] = append(folderDeps[folder], moduleFolder)
		}
	}

	output, err := json.MarshalIndent(folderDeps, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling output\n")
		return 1
	}

	fmt.Printf("%s\n", output)

	return 0
}

func main() {
	os.Exit(realMain())
}
