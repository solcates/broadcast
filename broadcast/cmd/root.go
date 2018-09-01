// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"github.com/solcates/broadcast"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	port    int
	payload string
	self    bool
	quiet   bool
	timeout int
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "Discover devices on your network via UDP Broadcasting",
	Long: `Discover devices on your network via UDP Broadcasting.
	
This simple tool is used to print a list of discovered UDP Servers/Services that respond to a payload broadcasted on a given port.	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Create a broadcaster instance
		bc := broadcast.NewUDPBroadcaster(port, payload)

		// to include this node as well, set findself to true
		bc.SetFindself(self)

		// set timeout if needed.
		bc.SetTimeout(timeout)

		// discover the nodes on this LAN
		nodes, err := bc.Discover()
		if err != nil {
			log.Fatal(err)
		}
		// iterate over the nodes found.
		for _, node := range nodes {
			if quiet {
				fmt.Println(node)
			} else {
				log.Infof("Found Node @ %v", node)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.broadcast.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolVar(&self, "self", false, "Discover this node as well")
	RootCmd.Flags().StringVar(&payload, "payload", "", "Set the Payload to send out during UDP broadcasting")
	RootCmd.Flags().IntVar(&port, "port", 0, "Set the UDP port to broadcast too")
	RootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Set the UDP port to broadcast too")
	RootCmd.MarkFlagRequired("port")
	RootCmd.MarkFlagRequired("payload")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".broadcast" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".broadcast")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
