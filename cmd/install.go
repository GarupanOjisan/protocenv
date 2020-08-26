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
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

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

type InstallOptions struct {
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
				fmt.Printf("   %s: %s\n", v.Name, v.ZipballURL)
			}
			return nil
		}

		if len(args) >= 1 {
			version := args[0]
			err := installVersion(version)
			if err != nil {
				return err
			}
			return nil
		}

		return errors.New("requires a installing version or some flags")
	},
}

var (
	showVersionList bool
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.SetUsageTemplate(`Usage:
	protocenv install <version>
	protocenv install -l|--list

	-l|--list		list all available versions
`)
	// Here you will define your flags and configuration settings.
	installCmd.Flags().BoolVarP(&showVersionList, "list", "l", false, "")

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

func installVersion(version string) error {
	srcPath, err := downloadZip(version)
	if err != nil {
		return err
	}

	// create protocenv home directory if not exists
	if err := initializeConfigDirectory(); err != nil {
		return err
	}

	// unzip to under protocenv home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dstPath := fmt.Sprintf("%s/.protocenv/versions/%s", homeDir, version)
	_, err = unzip(srcPath, dstPath)
	if err != nil {
		return err
	}

	return nil
}

func unzip(src, dst string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	rootPath := ""
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		path := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			if rootPath == "" {
				rootPath = path
			}
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return "", err
			}
		} else {
			buf := make([]byte, f.UncompressedSize64)
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				return "", err
			}
			if err = ioutil.WriteFile(path, buf, f.Mode()); err != nil {
				return "", err
			}
		}
		rc.Close()
	}
	return rootPath, nil
}

func downloadZip(version string) (string, error) {
	url := fmt.Sprintf("https://github.com/protocolbuffers/protobuf/releases/download/%s/protoc-%s-osx-x86_64.zip", version, version[1:])
	fmt.Printf("downloading %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(fmt.Sprintf("/tmp/%s", version))
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return out.Name(), err
}

func initializeConfigDirectory() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	versionsRoot := fmt.Sprintf("%s/.protocenv/versions", homedir)
	if _, err := os.Stat(versionsRoot); os.IsNotExist(err) {
		if err := os.MkdirAll(versionsRoot, 0774); err != nil {
			return err
		}
	}
	return nil
}
