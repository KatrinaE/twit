package main

import (
	"fmt"
	"github.com/spf13/viper"
	"local/twit/internal/twit"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	fmt.Println("a fanout loop")
	twit.FanoutLoop()
	fmt.Println("fanout loop done")
}
