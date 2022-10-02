package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Config map[string]interface{}

func InitConfig() {
	Config = map[string]interface{}{}
	// configContent := map[string]interface{}{}
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(byteValue, &Config)

	if os.Getenv("HEROKU_ENV") == "vats" {
		Config = Config["heroku"].(map[string]interface{})
	} else {
		Config = Config["local"].(map[string]interface{})
	}

}
