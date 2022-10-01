package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	echo "github.com/labstack/echo/v4"
)

func takeAuthCode(code string) map[string]interface{} {
	var client_id string
	var client_secret string
	hc := &http.Client{}

	client_id = Config["client_id"].(string)
	client_secret = Config["client_secret"].(string)

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("grant_type", "authorization_code")

	data.Set("redirect_uri", Config["redirect_uri"].(string))

	fmt.Println("DATA:", data)
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var result map[string]interface{}
	json.Unmarshal([]byte(bodyString), &result)
	fmt.Println("Access Token: ", result["access_token"])
	return result

}

func createCsv(response map[string]interface{}) []byte { // Generates csv data from customer object
	data := []byte("access_token,refresh_token\n" + fmt.Sprintf("%s,%s", response["access_token"], response["refresh_token"]))
	return data
}

var Wg sync.WaitGroup

type csvHolder struct {
	csvObject []byte
}

var CsvDict map[string]csvHolder
var Config map[string]interface{}

func main() {
	CsvDict = make(map[string]csvHolder)
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

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		url := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=%s&response_type=code&scope=https://www.googleapis.com/auth/youtube&redirect_uri=%s&state=state_parameter_passthrough_value&include_granted_scopes=true", Config["client_id"].(string), Config["redirect_uri"].(string))
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})
	e.GET("/auth_handler", func(c echo.Context) error {
		code := c.QueryParam("code")
		fmt.Println("CODE RECEIVED:", code)
		response := takeAuthCode(code)
		csvData := createCsv(response)
		csvId := string(rand.Intn(99999))
		CsvDict[csvId] = csvHolder{csvObject: csvData}
		return c.HTML(http.StatusOK, fmt.Sprintf(`<body onload="redirect()">You can use these informations while connecting to google's APIs<br></body>`+
			`<script>
				function wait(ms){
					var start = new Date().getTime();
					var end = start;
					while(end < start + ms) {
						end = new Date().getTime();
					};
				};
				function redirect(){
					wait(2000);
					window.location.href = "/auth_file?csvId=%s";
				};
			</script>`, csvId))
	})

	e.GET("/auth_file", func(c echo.Context) error {
		csvId := c.QueryParam("csvId")
		csv := CsvDict[csvId]
		delete(CsvDict, csvId)
		return c.Blob(http.StatusOK, "text/csv", csv.csvObject)

	})

	if os.Getenv("HEROKU_ENV") == "vats" {
		fmt.Println("Starting server on port:" + os.Getenv("PORT"))
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	} else {
		e.Logger.Fatal(e.Start(":9922"))
	}

}
