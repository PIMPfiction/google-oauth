package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// REFRESH TOKEN
// https://developers.google.com/identity/protocols/oauth2/web-server#offline

func TakeAuthCode(code string) map[string]interface{} {
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

func RefreshTokenRequest(refreshToken string) string {
	// make post request to url
	client_id := Config["client_id"].(string)
	client_secret := Config["client_secret"].(string)
	data := url.Values{}
	data.Set("client_id", client_id)
	data.Set("client_secret", client_secret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")
	hc := &http.Client{}
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	fmt.Println("response Body:", result)
	authToken := result["access_token"].(string)
	return authToken
}

func CheckAccessToken(accessToken string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/tokeninfo?access_token="+accessToken, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	if result["error"] != nil {
		return false
	}
	if result["expires_in"].(float64) > 0 {
		return true
	} else {
		return false
	}
}
