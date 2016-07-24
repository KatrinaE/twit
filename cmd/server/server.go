// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/golang/go/wiki/SQLInterface
// hi
package main

import (
	"fmt"
	"github.com/spf13/viper"
	"local/twit/internal/twit"
	"log"
	"net/http"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../../")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	mux := twit.RegisterRoutes()
	http.Handle("/", mux)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
