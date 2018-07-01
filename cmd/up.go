// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
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

package cmd

import (
	"fmt"
	"path"

	"github.com/dbarzdys/docktest/docker"
	"github.com/dbarzdys/docktest/export"
	"github.com/dbarzdys/docktest/lock"
	"github.com/spf13/cobra"
)

var exportFile *string

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Spin up test containers and export variables to .env file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = initConfig(); err != nil {
			return
		}
		if lock.Exists(lockFile) {
			err = fmt.Errorf("lock file found: %s", lockFile)
			return err
		}
		if export.Exists(*exportFile) {
			err = export.Remove(*exportFile)
			if err != nil {
				return
			}
		}
		client, err := docker.NewClient("")
		if err != nil {
			return err
		}
		runner := docker.NewRunner(client)
		containers, err := runner.Run(cfg.Services)
		if err != nil {
			return err
		}
		ids := make([]string, len(containers))
		for i, c := range containers {
			ids[i] = c.ID
		}
		err = lock.Create(
			lockFile,
			lock.LockConfig{
				ExportFile: path.Join(path.Dir(lockFile), *exportFile),
				Containers: ids,
			},
		)
		if err != nil {
			return err
		}
		err = export.Create(
			*exportFile,
			cfg.Export,
			containers,
		)
		if err != nil {
			return err
		}
		fmt.Printf("DockTest started with %d containers\n", len(containers))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	exportFile = upCmd.Flags().StringP("export_file", "e", "./docktest.env", "Variable export file location")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
