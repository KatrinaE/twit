// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/golang/go/wiki/SQLInterface
// hi
package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"local/twit/internal/twit"
	"log"
	"net/http"
	"strconv"
)

func main() {
	log.Print("Starting up...")
	log.Print("Parsing flags")
	var configPath string
	var configFilename string
	var port int
	flag.StringVar(&configPath, "configpath", ".",
		"path to configuration file (absolute or relative)")
	flag.StringVar(&configFilename, "configfile", "config",
		"name of config file (no extension")
	flag.IntVar(&port, "port", 8080, "port for server to listen on")
	flag.Parse()

	log.Print("Setting config")
	viper.SetConfigType("yaml")
	viper.SetConfigName(configFilename)
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	log.Print("Registering routes")
	mux := twit.RegisterRoutes()
	http.Handle("/", mux)
	portStr := strconv.Itoa(port)
	log.Print("Starting HTTP listener on port " + portStr)
	err = http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Print("Listening on port " + portStr)
}
