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
	"os"
	"path"
	"strings"

	"github.com/dbarzdys/docktest/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

var cfgFile *string
var lockFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "docktest",
	Short:        "Use docker to run integration tests againts other services.",
	SilenceUsage: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	cfgFile = rootCmd.PersistentFlags().StringP("config", "c", "./docktest.yaml", "config file")

}

func initConfig() (err error) {
	cfg, err = config.New(*cfgFile)
	if err != nil {
		return
	}
	lockFile = getLockPath()
	for k, e := range cfg.Export {
		e = os.Expand(e, func(s string) string {
			original := fmt.Sprintf("${%s}", s)
			path := strings.Split(s, ".")
			if len(path) == 0 {
				return original
			}
			switch path[0] {
			case "constants":
				if len(path) > 2 {
					return original
				}
				found, ok := cfg.Constants[path[1]]
				if !ok {
					return original
				}
				return found
			default:
				return original
			}
		})
		cfg.Export[k] = e
	}
	for name, svc := range cfg.Services {
		for k, v := range svc.Env {
			v = os.Expand(v, func(s string) string {
				original := fmt.Sprintf("${%s}", s)
				path := strings.Split(s, ".")
				if len(path) == 0 {
					return original
				}
				switch path[0] {
				case "constants":
					if len(path) > 2 {
						return original
					}
					found, ok := cfg.Constants[path[1]]
					if !ok {
						return original
					}
					return found
				default:
					return original
				}
			})
			svc.Env[k] = v
		}
		cfg.Services[name] = svc
	}
	return
}

func getLockPath() string {
	dir, file := path.Split(*cfgFile)
	return path.Join(dir, fmt.Sprintf(".%s.lock", file))
}
