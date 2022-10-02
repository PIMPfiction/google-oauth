package main

import (
	helpers "auth-server/helpers"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	echo "github.com/labstack/echo/v4"
)

func createCsv(response map[string]interface{}) []byte { // Generates csv data from customer object
	data := []byte("access_token,refresh_token\n" + fmt.Sprintf("%s,%s", response["access_token"], response["refresh_token"]))
	return data
}

var Wg sync.WaitGroup

type csvHolder struct {
	csvObject []byte
}

var CsvDict map[string]csvHolder

func main() {
	CsvDict = make(map[string]csvHolder)
	helpers.InitConfig()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		fmt.Println(helpers.Config["client_id"].(string), helpers.Config["redirect_uri"].(string))
		url := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?access_type=offline&client_id=%s&response_type=code&scope=https://www.googleapis.com/auth/youtube&redirect_uri=%s&state=state_parameter_passthrough_value&include_granted_scopes=true", helpers.Config["client_id"].(string), helpers.Config["redirect_uri"].(string))
		return c.Redirect(http.StatusTemporaryRedirect, url)
	})
	e.GET("/auth_handler", func(c echo.Context) error {
		code := c.QueryParam("code")
		fmt.Println("CODE RECEIVED:", code)
		response := helpers.TakeAuthCode(code)
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

	e.GET("/check_token", func(c echo.Context) error { // checks token validity, if not valid returns new token
		var status bool
		token := c.QueryParam("token")
		refresh_token := c.QueryParam("refresh_token")
		response := helpers.CheckAccessToken(token)
		if !response { // if token is valid
			status = false
			token = helpers.RefreshTokenRequest(refresh_token)
		} else {
			status = true
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"valid": status,
			"token": token,
		})
	})

	if os.Getenv("HEROKU_ENV") == "vats" {
		fmt.Println("Starting server on port:" + os.Getenv("PORT"))
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	} else {
		e.Logger.Fatal(e.Start(":9922"))
	}

}
