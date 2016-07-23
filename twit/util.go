package twit

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
)

func writeJson (??, ??) {
	b, err := json.Marshal(tweet)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func GetDbConfig() (string, string) {
	env := viper.GetString("environment")
	dbDriverField := fmt.Sprintf("%s.driver", env)
	dbDriver := viper.GetString(dbDriverField)
	openField := fmt.Sprintf("%s.open", env)
	dbOpen := viper.GetString(openField)
	return dbDriver, dbOpen
}
