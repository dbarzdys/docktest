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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dbarzdys/docktest/docker"
	"github.com/dbarzdys/docktest/lock"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Spin-up test containers and export variables to next command",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(c *cobra.Command, args []string) (err error) {
		if err = initConfig(); err != nil {
			return
		}
		if lock.Exists(lockFile) {
			err = fmt.Errorf("lock file found: %s", lockFile)
			return err
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
				ExportFile: "./docktest.env",
				Containers: ids,
			},
		)
		if err != nil {
			return err
		}
		fmt.Printf("DockTest started with %d containers\n", len(containers))
		if len(args) == 1 {
			args = strings.Split(args[0], " ")
		}
		cmd := exec.Command(args[0], args[1:]...)
		stderr, _ := cmd.StderrPipe()
		for k, v := range cfg.Export {
			os.Setenv(k, v)
		}
		cmd.Start()
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		cmd.Wait()
		if !lock.Exists(lockFile) {
			fmt.Printf("Could not find lock file: %s\n", lockFile)
			fmt.Println("Nothing to remove")
			return
		}
		lockConfig, err := lock.Read(lockFile)
		if err != nil {
			return
		}
		remover := docker.NewRemover(client)
		err = remover.Remove(lockConfig.Containers...)
		if err != nil {
			return
		}
		err = lock.Remove(lockFile)
		if err != nil {
			return
		}
		fmt.Println("DockTest containers successfuly removed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
