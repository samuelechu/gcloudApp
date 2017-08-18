package main

import (
    "fmt"
    "net/http"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/appengine"
    "os"
    "log"
    "google.golang.org/api/drive/v2"
    "golang.org/x/net/context"
)

var (
    googleOauthConfig = &oauth2.Config{
        RedirectURL:    "https://gotesting-175718.appspot.com/googleCallback",
        ClientID:     os.Getenv("CLIENT_ID"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
        ClientSecret: os.Getenv("CLIENT_SECRET"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
        Scopes:       []string{"https://www.googleapis.com/auth/drive", "https://www.googleapis.com/auth/drive.file", "https://www.googleapis.com/auth/gmail.readonly"},
        Endpoint:     google.Endpoint,
    }
// Some random string, random for each request
    oauthStateString = "random"
)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`

func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/GoogleLogin", handleGoogleLogin)
    http.HandleFunc("/googleCallback", handleGoogleCallback)

    log.Print("Listening on port 8080")
    http.ListenAndServe(":8080", nil)
    appengine.Main()
}

func handleMain(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
    if appengine.IsDevAppServer(){
        googleOauthConfig.RedirectURL = "https://8080-dot-2979131-dot-devshell.appspot.com/googleCallback"
    }

    url := googleOauthConfig.AuthCodeURL(oauthStateString)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")

    token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
    if err != nil {
        log.Print("oauthConf.Exchange() failed with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

    //response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
    driveService, err := drive.New(client)

    myR, err := driveService.Files.List().MaxResults(10).Do()
    if err != nil {
        fmt.Fprintf(w, "Couldn't retrieve files ", err)
    }
    if len(myR.Items) > 0 {
        for _, i := range myR.Items {
            fmt.Fprintf(w, i.Title, " ", i.Id)
        }
    } else {
        fmt.Fprintf(w, "No files found.")
    }

    //defer response.Body.Close()
    //contents, err := ioutil.ReadAll(response.Body)
    //fmt.Fprintf(w, "Content: %s\n", contents)
}