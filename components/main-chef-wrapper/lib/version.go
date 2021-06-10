//
// Copyright (c) Chef Software, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package lib

import (
	"fmt"
	"os"
	"path"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/chef/chef-workstation/components/main-chef-wrapper/dist"
)



var gemManifestMap map[string]interface{}
var manifestMap map[string]interface{}

func init() {
	gemManifestMap = gemManifestHash()
	manifestMap = manifestHash()
}
func Version(){
	if omnibusInstall() == true {
		showVersionViaVersionManifest()
	} else {
		fmt.Fprintln(os.Stderr, "ERROR:", "Can not find omnibus installation directory for", dist.WorkstationProduct)
	}
}

func showVersionViaVersionManifest()  {
	fmt.Printf("%v version: %v", dist.WorkstationProduct, componentVersion("build_version")  )
	productMap := map[string]string{
		dist.ClientProduct: dist.CLIWrapperExec,
		dist.InspecProduct: dist.InspecCli,
		dist.CliProduct: dist.CliGem,
		dist.HabProduct: dist.HabSoftwareName,
		"Test Kitchen": "test-kitchen",
		"Cookstyle": "cookstyle",
	}
	for prodName, component := range productMap {
		fmt.Printf("\n%v version: %v", prodName, componentVersion(component)  )
	}
	fmt.Printf("\n")
}
func componentVersion(component string) string  {
	v, ok := gemManifestMap[component]
	if ok {
		stringifyVal := v.([]interface{})[0]
		return stringifyVal.(string)
	} else if v, ok := manifestMap[component]; ok {
		return v.(string)
	} else {
		success, _ := Dig(manifestMap, "software", component, "locked_version")
		if success == nil {
			return "unknown"
		}else {
			return success.(string)
		}
	}
}
func gemManifestHash() map[string]interface{} {
	filepath := path.Join(omnibusRoot(),"gem-version-manifest.json")
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var gemManifestHash map[string]interface{}
	json.Unmarshal([]byte(byteValue), &gemManifestHash)
	defer jsonFile.Close()
	return gemManifestHash
}
func manifestHash() map[string]interface{}  {
	filepath := path.Join(omnibusRoot(),"version-manifest.json")
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var manifestHash map[string]interface{}
	json.Unmarshal([]byte(byteValue), &manifestHash)
	defer jsonFile.Close()
	return manifestHash
}

func omnibusInstall() bool {
	//# We also check if the location we're running from (omnibus_root is relative to currently-running ruby)
	//# includes the version manifest that omnibus packages ship with. If it doesn't, then we're running locally
	//# or out of a gem - so not as an 'omnibus install'
	ExpectedOmnibusRoot := ExpectedOmnibusRoot()
	if _, err := os.Stat(ExpectedOmnibusRoot); err == nil {
		if _, err = os.Stat(path.Join(ExpectedOmnibusRoot,"version-manifest.json")); err == nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}


func omnibusRoot()string  {
	omnibusroot, err := filepath.Abs(path.Join(ExpectedOmnibusRoot()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", "Can not find omnibus installation directory for", dist.WorkstationProduct)
	}
	return omnibusroot
}

func ExpectedOmnibusRoot()string {
	groot := os.Getenv("GEM_ROOT")
	rootPath, err := filepath.Abs(path.Join(groot,"..","..", "..", "..", ".."))
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
	}
	return rootPath
}
