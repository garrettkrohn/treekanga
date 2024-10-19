/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package main

import (
	"fmt"
	"log"

	"github.com/garrettkrohn/treekanga/cmd"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("treekanga")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/treekanga")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	initConfig()
	fmt.Print(viper.GetString("repos.test.defaultBranch"))
	fmt.Print(viper.GetStringSlice("repos.offensive-security-platform.zoxideFolders"))

	cmd.Execute()
}
