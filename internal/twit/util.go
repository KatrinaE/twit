package twit

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
)

func writeErrorResponse(w http.ResponseWriter, err error) {
	Debug(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func writeJsonResponse(w http.ResponseWriter, val interface{}) {
	b, err := json.Marshal(val)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func getDbConfig() (string, string) {
	env := viper.GetString("environment")
	dbDriverField := fmt.Sprintf("%s.driver", env)
	dbDriver := viper.GetString(dbDriverField)
	openField := fmt.Sprintf("%s.open", env)
	dbOpen := viper.GetString(openField)
	return dbDriver, dbOpen
}

func getRedisConfig() (string, string, int) {
	env := viper.GetString("environment")
	addressField := fmt.Sprintf("%s.redis.address", env)
	address := viper.GetString(addressField)
	passwordField := fmt.Sprintf("%s.redis.password", env)
	password := viper.GetString(passwordField)
	dbField := fmt.Sprintf("%s.redis.db", env)
	dbStr := viper.GetString(dbField)
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		log.Fatalf("Non-integer value for db: %s", dbStr)
	}
	return address, password, db
}
