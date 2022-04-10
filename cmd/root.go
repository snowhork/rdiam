/*
Copyright © 2021 Kyoshiro Maruo

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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/snowhork/rdiam/pkg/redash"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rdiam",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		verbose, flagErr := rootCmd.Flags().GetBool("verbose")
		if flagErr != nil {
			fmt.Printf("failed to parse verbose flag %+v\n", flagErr)
		}

		if verbose {
			if e, ok := err.(interface{ StackTrace() errors.StackTrace }); ok {
				fmt.Printf("%s\n%+v\n", e, e.StackTrace())
			} else {
				fmt.Println("failed to get stack trace")
			}
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rdiam.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "print verbose log")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".rdiam" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".rdiam")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if err := setGlobalConfigFromViper(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
			return
		}

		fmt.Println("Failed to use config file:", viper.ConfigFileUsed())
	}

	fmt.Print("Enter you Redash endpoint (e.g. https://redash.yourdomain.com): ")
	redashEndpoint, err := enterValue()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Enter your Redash user API Key (available at %s/users/me): ", redashEndpoint)
	redashUserAPIKey, err := enterValue()
	if err != nil {
		panic(err)
	}

	viper.Set("RedashEndpoint", redashEndpoint)
	viper.Set("RedashUserAPIKey", redashUserAPIKey)
	if err := viper.WriteConfig(); err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			panic(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := setGlobalConfigFromViper(); err != nil {
		panic(err)
	}

	fmt.Printf("settings are writen to %s\n", viper.ConfigFileUsed())
}

func setGlobalConfigFromViper() error {
	conf := config{
		RedashEndPoint:   strings.Trim(viper.Get("RedashEndpoint").(string), "/"),
		RedashUserAPIKey: viper.Get("RedashUserAPIKEY").(string),
	}

	if conf.RedashEndPoint == "" {
		return errors.New("RedashEndPoint is empty")
	}
	if conf.RedashUserAPIKey == "" {
		return errors.New("RedashUserAPIKey is empty")
	}
	globalClient = redash.NewClient(conf.RedashEndPoint, conf.RedashUserAPIKey)
	return nil
}

func enterValue() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')

	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(res), nil
}

type config struct {
	RedashEndPoint   string
	RedashUserAPIKey string
}

var globalClient *redash.Client
