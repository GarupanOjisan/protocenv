/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// globalCmd represents the global command
var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "set global protoc version",
	Long:  `set global protoc version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		// check if the specified version has been installed
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configRoot := fmt.Sprintf("%s/.protocenv", homeDir)
		versionDir := fmt.Sprintf("%s/versions/%s", configRoot, version)
		if _, err := os.Stat(versionDir); os.IsNotExist(err) {
			return fmt.Errorf("%s is not installed", version)
		}

		// create symbolic link
		oldName := fmt.Sprintf("%s/bin", versionDir)
		newName := fmt.Sprintf("%s/bin", configRoot)
		if _, err := os.Lstat(newName); err == nil {
			if err := os.Remove(newName); err != nil {
				return err
			}
		}
		if err := os.Symlink(oldName, newName); err != nil {
			return fmt.Errorf("symlink() failed: %v", err)
		}

		fmt.Printf(`
Now you can use protoc %s globally.
Do not forget set PATH : %s
`, version, newName)

		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires version")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(globalCmd)

	// Here you will define your flags and configuration settings.
	globalCmd.SetUsageTemplate(`Usage:
protocenv global <version>
`)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// globalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// globalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
