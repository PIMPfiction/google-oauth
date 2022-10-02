
# Google Oauth2 Callback Server

This server can be used as callback url in Google APIs. You can authorize google users through your API with these server. Server will redirect user to google authentication/login page. User will be prompted to confirm your app. After confirmation; Server will generate Bearer token which can be used in google apis. 
Callback server will generate acccess token and refresh token.
- How to create Google oauth client - https://www.youtube.com/watch?v=rUTc4Yl7-_I
- Demo Video on Heroku - ( coming soon)

## Directly Run from Source

- Clone the project
    ```bash
    git clone https://github.com/PIMPfiction/google-oauth
    ```
- Go to your project folder
    ```bash
    cd google-oauth
    ```
- Configure config.json before running the server. If you have golang simply execute ⬇
    ```bash
    go run main.go
    ```

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
