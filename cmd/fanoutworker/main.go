package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"local/twit"
	"log"
)

func main() {
	log.Print("Starting fanout worker...")
	log.Print("Parsing flags")
	var configPath string
	var configFilename string
	flag.StringVar(&configPath, "configpath", ".",
		"path to configuration file (absolute or relative)")
	flag.StringVar(&configFilename, "configfile", "config",
		"name of config file (no extension")
	flag.Parse()

	log.Print("Setting config")
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFilename)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Running fanout loop")
	twit.FanoutLoop(db)
}
