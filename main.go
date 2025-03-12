/*
Copyright © 2024 Garrett Krohn <garrettkrohn@gmail.com>
*/
package main

import (
	"log"

	"github.com/garrettkrohn/treekanga/cmd"
	"github.com/spf13/viper"
)

var version = "dev"

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

	cmd.Execute(version)
}
