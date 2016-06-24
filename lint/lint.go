/**
 * Copyright (C) 2015 Hewlett Packard Enterprise Development LP
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

//Package lints implements linting methods
package lint

import (
	//"fmt"

	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var lintersFlag = map[string]string{
	"aligncheck":  `aligncheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`,
	"deadcode":    `deadcode  `,
	"dupl":        `dupl -plumbing -threshold {duplthreshold} ./*.go:^(?P<path>[^\s][^:]+?\.go):(?P<line>\d+)-\d+:\s*(?P<message>.*)$`,
	"errcheck":    `errcheck -abspath`,
	"goconst":     `goconst -min-occurrences {min_occurrences} .:PATH:LINE:COL:MESSAGE`,
	"gocyclo":     `gocyclo -over {mincyclo} .:^(?P<cyclo>\d+)\s+\S+\s(?P<function>\S+)\s+(?P<path>[^:]+):(?P<line>\d+):(\d+)$`,
	"gofmt":       `gofmt -l -s ./*.go:^(?P<path>[^\n]+)$`,
	"goimports":   `goimports -w `,
	"golint":      "golint -min_confidence {min_confidence} .:PATH:LINE:COL:MESSAGE",
	"gotype":      "gotype -e {tests=-a} .:PATH:LINE:COL:MESSAGE",
	"ineffassign": `ineffassign -n .:PATH:LINE:COL:MESSAGE`,
	"interfacer":  `interfacer ./:PATH:LINE:COL:MESSAGE`,
	"lll":         `lll -g -l {maxlinelength} ./*.go:PATH:LINE:MESSAGE`,
	"structcheck": `structcheck {tests=-t} .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):\s*(?P<message>.+)$`,
	"test":        `go test:^--- FAIL: .*$\s+(?P<path>[^:]+):(?P<line>\d+): (?P<message>.*)$`,
	"testify":     `go test:Location:\s+(?P<path>[^:]+):(?P<line>\d+)$\s+Error:\s+(?P<message>[^\n]+)`,
	"varcheck":    `varcheck .:^(?:[^:]+: )?(?P<path>[^:]+):(?P<line>\d+):(?P<col>\d+):[\s\t]+(?P<message>.*)$`,
	"vet":         "go tool vet ./*.go:PATH:LINE:MESSAGE",
	"vetshadow":   "go tool vet --shadow ./*.go:PATH:LINE:MESSAGE",
	"unconvert":   "unconvert .:PATH:LINE:COL:MESSAGE",
	"gosimple":    "gosimple .:PATH:LINE:COL:MESSAGE",
	"staticcheck": "staticcheck .:PATH:LINE:COL:MESSAGE",
	"misspell":    "misspell ./*.go:PATH:LINE:COL:MESSAGE",
}

var packageLinters = []string{
	"errcheck",
}
var dirLinters = []string{
	"goimports",
}

func LintPackages(listOfPackages []string) (err error) {
	//Find out what the Root Path is
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}

	for _, aPackage := range listOfPackages {
		log.Debug("Checking package:", aPackage)
		for _, linter := range packageLinters {
			log.Debug("Running package checker:", linter)
			cmdString := lintersFlag[linter]
			splitCmd := strings.Split(cmdString, " ")
			msg := fmt.Sprintf("CMD: %s %s %s", splitCmd[0], splitCmd[1], aPackage)
			log.Debug(msg)
			cmd := exec.Command(splitCmd[0], splitCmd[1], aPackage)
			out, err := cmd.CombinedOutput()
			if err != nil {
				err = fmt.Errorf("%s found error in:%s", linter, string(out))
				log.Error(err)
				return err
			}
		}

	}

	return nil
}

func LintDirs(listOfDirs []string) (err error) {
	//Find out what the Root Path is
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	tmpRootPath, err := cmd.Output()
	if err != nil {
		return err
	}

	//Trim return character
	rootPath := strings.TrimSpace(string(tmpRootPath))
	err = os.Chdir(rootPath)
	if err != nil {
		return err
	}

	for _, aDir := range listOfDirs {
		log.Debug("Checking dir:", aDir)
		for _, linter := range dirLinters {
			log.Debug("Running dir checker:", linter)
			cmdString := lintersFlag[linter]
			splitCmd := strings.Split(cmdString, " ")
			msg := fmt.Sprintf("CMD: %s %s %s", splitCmd[0], splitCmd[1], aDir)
			log.Debug(msg)
			cmd := exec.Command(splitCmd[0], splitCmd[1], aDir)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Debug("args:", cmd.Args)
				err = fmt.Errorf("%s found error in:%s", linter, string(out))
				log.Error(err)
				return err
			}
		}

	}

	return nil
}