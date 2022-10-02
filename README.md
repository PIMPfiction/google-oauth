
# Google Oauth2 Callback Server

This server can be used as callback url in Google APIs. You can authorize google users through your API with these server. Server will redirect user to google authentication/login page. User will be prompted to confirm your app. After confirmation; Server will generate Bearer token which can be used in google apis. 
Callback server will generate acccess token and refresh token.
- How to create Google oauth client - https://www.youtube.com/watch?v=rUTc4Yl7-_I
- Demo Video on Heroku - ( coming soon)


## API Reference

#### Main Route - Redirect Users to Google Authentication Page

```http
  GET /
```

#### Check and Refresh Token - Checks validity of access token, if it is expired it will be refreshed

```http
  GET /check_token/?token=${token}&refresh_token=${refresh_token}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `token`      | `string` | **Required**. Your Access Token (from downlaoded csv file) |
| `refresh_token`      | `string` | **Required**.Your Refresh Token (from downlaoded csv file) |


## Directly Run from Source

- Clone the project
    ```bash
    git clone https://github.com/PIMPfiction/google-oauth
    ```
- Go to your project folder
    ```bash
    cd google-oauth
    ```
- Edit config.json, replace client-id and client-secret. (which can be accessible from https://console.cloud.google.com/apis/credentials)
- Configure config.json before running the server. If you have golang simply execute ⬇
    ```bash
    go run main.go
    ```
- Note: You should add 127.0.0.1:9922/auth_handler to authorized redirect uri from google console. 
    - Look at these screenshots. 
        - https://prnt.sc/npKp3d-72MfT
        - https://prnt.sc/ZZ4SxJ-GecKY
        - https://prnt.sc/ZmZsLRjpd07q (your config should look like this)

## Deploy on Heroku (so it will be accessible by all users)
  
- Create account on heroku  (https://www.heroku.com/)

- Install heroku CLI - (https://devcenter.heroku.com/articles/heroku-cli#install-the-heroku-cli)

- Run below code to login to heroku ⬇
    ```bash
        heroku login
    ```
- Then execute below commands in order to deploy code to heroku 
    ```bash
        git init
        git add . 
        git commit -m "First commit of google-oauth-server"
    ```
    Below command will give you your apps link (it is important) ⚠ 
    Copy that link to your json file, and replace "redirect_uri"
    ```bash
        heroku create -a your-app-name
        heroku config:set HEROKU_ENV=vats --app your-app-name
    ```

- Edit config.json file. You will edit values in "heroku" key. Replace 'redirect_uri' with heroku app's link.
    ```text
        For example your heroku url is https://awesome-server-2299.herokuapp.com/
        So your redirect_uri must be "https://awesome-server-2299.herokuapp.com/auth_handler"
    ```

- Now we are executing final command, this will deploy our app on heroku.
    ```bash
        git push heroku main
    ```
