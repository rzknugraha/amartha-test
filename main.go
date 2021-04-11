package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: " + viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
	err := viper.Unmarshal(&Config{})
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Database struct {
		Driver   string `json:"driver"`
		Dsn      string `json:"dsn"`
		MaxConns int    `json:"max_cons"`
		MaxIdle  int    `json:"max_idle"`
	} `json:"database"`
}

func main() {

}
