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

//Package installer handles installation of packages
package installer

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"fmt"
)

var installMap = map[string]string{
	"golint":      "github.com/golang/lint/golint",
	"gotype":      "golang.org/x/tools/cmd/gotype",
	"goimports":   "golang.org/x/tools/cmd/goimports",
	"errcheck":    "github.com/kisielk/errcheck",
	"varcheck":    "github.com/opennota/check/cmd/varcheck",
	"structcheck": "github.com/opennota/check/cmd/structcheck",
	"aligncheck":  "github.com/opennota/check/cmd/aligncheck",
	"deadcode":    "github.com/tsenart/deadcode",
	"gocyclo":     "github.com/alecthomas/gocyclo",
	"misspell":    "github.com/client9/misspell/cmd/misspell",
	"ineffassign": "github.com/gordonklaus/ineffassign",
	"dupl":        "github.com/mibk/dupl",
	"interfacer":  "github.com/mvdan/interfacer/cmd/interfacer",
	"lll":         "github.com/walle/lll/cmd/lll",
	"unconvert":   "github.com/mdempsky/unconvert",
	"goconst":     "github.com/jgautheron/goconst/cmd/goconst",
	"gosimple":    "honnef.co/go/simple/cmd/gosimple",
	"staticcheck": "honnef.co/go/staticcheck/cmd/staticcheck",
}

//CheckExternalDependencies checks if a required component is installed, if not, it go gets it.
func CheckExternalDependencies()(err error) {
	log.SetLevel(log.DebugLevel)
	log.Debug("Hello world")

	for key := range installMap {
		packageToGet := installMap[key]
		log.Debug("Getting package:",packageToGet)
		_, err := exec.LookPath(key)

		if err != nil {
			log.Debug("Need to install:",packageToGet)
			cmd := exec.Command("go", "get", packageToGet,)
			getOut, err := cmd.Output()
			if err != nil {
				log.Error(err)
				return err
			}
			if _, err := exec.LookPath(key); err != nil {
				nErr := fmt.Errorf("After installing %s, still can't find it:%s",key,err)
				log.Error(nErr.Error())
				return nErr
			}

			log.Debug(getOut)
		}
		log.Debug("Found:",key)
	}
	return err

}