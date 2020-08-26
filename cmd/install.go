/*
Copyright © 2020 Motohiro Nakamura <nakamura.motohiro.private@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type versions []version

type version struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		Sha string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

// installCmd represents the versions command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install specified version",
	Long: `install specified version:

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersionList {
			versions, err := getAllVersions()
			if err != nil {
				return fmt.Errorf("failed to get list of versions: %v", err)
			}
			for _, v := range versions {
				fmt.Printf("   %s\n", v.Name)
			}
		}

		return nil
	},
}

var (
	showVersionList bool
)

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.
	installCmd.Flags().BoolVarP(&showVersionList, "list", "l", false, "list all available versions")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getAllVersions() ([]*version, error) {
	var versions []*version
	for i := 1; ; i++ {
		vs, err := getVersions(i)
		if err != nil {
			return nil, err
		}
		if len(vs) == 0 {
			break
		}
		versions = append(versions, vs...)
	}
	return versions, nil
}

func getVersions(page int) ([]*version, error) {
	var versions []*version
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/protocolbuffers/protobuf/tags?page=%d", page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}
