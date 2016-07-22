package twit

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetDbConfig() (string, string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("dbconf")
	viper.AddConfigPath("./db/") // right now dbconf is only config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	env := viper.GetString("environment")
	dbDriverField := fmt.Sprintf("%s.driver", env)
	dbDriver := viper.GetString(dbDriverField)
	openField := fmt.Sprintf("%s.open", env)
	dbOpen := viper.GetString(openField)
	return dbDriver, dbOpen
}
